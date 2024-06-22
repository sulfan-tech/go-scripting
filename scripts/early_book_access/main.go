package main

import (
	"context"
	"fmt"
	"go-scripting/entities"
	csvprocessor "go-scripting/pkg/csv_processor"
	"go-scripting/pkg/logger"
	"go-scripting/script/early_book_access/service"
	"log"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

const (
	filePath   = "input.csv"
	numWorkers = 100 // Number of concurrent workers
)

type CSVRow struct {
	MemberID string
	Uid      string
}

var ebaService service.EBAInterface

func init() {
	env := "../../.env"
	err := godotenv.Load(env)
	if err != nil {
		log.Fatal("Error loading .env file:", err.Error())
	}

	ebaService = service.NewUserEBAService()
}

func main() {
	csvReader := csvprocessor.NewCSVReader()

	data, err := csvReader.ReadCSV(context.Background(), filePath)
	if err != nil {
		log.Fatal("Error reading CSV:", err)
		return
	}

	mappedData := mapToCSVRow(data)

	// Channel for CSV rows and for results
	rowChan := make(chan CSVRow, len(mappedData))
	resultChan := make(chan error, len(mappedData))
	var wg sync.WaitGroup

	// Launch worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(rowChan, resultChan, &wg)
	}

	// Send rows to rowChan
	go func() {
		for _, row := range mappedData {
			rowChan <- row
		}
		close(rowChan)
	}()

	// Close resultChan after all workers are done
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Process results
	rowNumber := 1
	for err := range resultChan {
		if err != nil {
			log.Printf("Error processing row %d: %v", rowNumber, err)
		} else {
			log.Printf("Successfully processed row %d", rowNumber)
		}
		rowNumber++
	}

	log.Println("Processing completed.")
}

func worker(rowChan <-chan CSVRow, resultChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	for row := range rowChan {
		err := processRow(row)
		resultChan <- err
	}
}

func processRow(row CSVRow) error {
	ctx := context.Background()

	eba := entities.EBA{
		UserID:        row.Uid,
		Slot:          1,
		ExpiredDate:   time.Now().AddDate(0, 0, 14),
		AvailableFrom: time.Now().AddDate(0, 0, -1),
	}
	bypass := entities.BYPASS{
		Uid:    row.Uid,
		Bypass: true,
	}

	req := service.RequestUserEBA{
		InsertEBA: eba,
		BypassEBA: bypass,
	}

	log.Printf("Processing row - MemberID: %s\n", row.MemberID)

	err := ebaService.InsertUserEBA(ctx, req)
	if err != nil {
		log.Printf("Error processing row - MemberID: %s: %s", row.MemberID, err.Error())
		return err
	}

	logger.LogInfo(fmt.Sprintf("MemberID (%v) successfully Add to EBA", row.MemberID))

	return nil
}

func mapToCSVRow(data [][]string) []CSVRow {
	var result []CSVRow

	for _, row := range data {
		csvRow := CSVRow{
			MemberID: row[0],
			Uid:      row[1],
		}

		result = append(result, csvRow)
	}

	return result
}
