package main

import (
	"context"
	"fmt"
	csvprocessor "go-scripting/pkg/csv_processor"
	"go-scripting/pkg/logger"
	"go-scripting/scripts/adjust_expired_membership/service"
	"log"
	"sync"

	"github.com/joho/godotenv"
)

const (
	filePath              = "input.csv"
	missingAttributeError = "error: %s is empty in CSV row %d"
	unsupportedTypeError  = "error: Unsupported Type '%s' in CSV row %d"
	operationFailedError  = "operation failed on Row %d: %v"
	numWorkers            = 50 // Number of workers
)

type CSVRow struct {
	MemberId  string
	DaysToAdd string
	Logs      string
}

var adjustExpiredMembershipService service.AdjustmentExpiredMembershipImpl

func init() {
	env := "../../.env"
	err := godotenv.Load(env)
	if err != nil {
		log.Fatal("Error loading .env file:", err.Error())
	}

	adjustExpiredMembershipService = service.NewInstanceAdjustmentExpiredMembership()
	logger.Init()
}

func main() {
	defer logger.CloseLogFile()

	csvReader := csvprocessor.NewCSVReader()

	data, err := csvReader.ReadCSV(context.Background(), filePath)
	if err != nil {
		log.Fatal("Error reading CSV:", err)
		return
	}
	mappedData := mapToCSVRow(data)

	rowChan := make(chan CSVRow, len(mappedData))
	resultChan := make(chan error, len(mappedData))
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(rowChan, resultChan, &wg)
	}

	go func() {
		for _, row := range mappedData {
			rowChan <- row
		}
		close(rowChan)
	}()

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	rowNumber := 1
	for err := range resultChan {
		if err != nil {
			log.Printf("Error processing row %d: %v", rowNumber, err)
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
	log.Printf("Processing - MemberId: %s, Days: %s\n", row.MemberId, row.DaysToAdd)

	err := adjustExpiredMembershipService.AdjustExpiredMembership(context.Background(), row.MemberId)
	if err != nil {
		log.Println("Error", err)
		logger.LogError(fmt.Sprintf("Member: (%s) failed", row.MemberId))
		return err
	}

	logger.LogInfo(fmt.Sprintf("Add (%v) Day on Member: (%s) updated successfully", row.DaysToAdd, row.MemberId))
	return nil
}

func mapToCSVRow(data [][]string) []CSVRow {
	var result []CSVRow

	for _, row := range data {
		csvRow := CSVRow{
			MemberId:  row[0],
			DaysToAdd: row[1],
			Logs:      row[2],
		}

		result = append(result, csvRow)
	}

	return result
}
