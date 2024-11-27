package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	csvprocessor "go-scripting/pkg/csv_processor"
	"go-scripting/pkg/logger"
	"go-scripting/scripts/inject_free_days/service"

	"github.com/joho/godotenv"
)

const (
	filePath              = "input.csv"
	missingAttributeError = "error: %s is empty in CSV row %d"
	unsupportedTypeError  = "error: Unsupported Type '%s' in CSV row %d"
	operationFailedError  = "operation failed on Row %d: %v"
	numWorkers            = 100 // Number of workers
)

type CSVRow struct {
	ExecutionID string
	MemberId    string
	DaysToAdd   string
	Logs        string
}

var freeDaysService service.FreeDaysInterface

func init() {
	env := "../../.env"
	err := godotenv.Load(env)
	if err != nil {
		log.Fatal("Error loading .env file:", err.Error())
	}

	freeDaysService = service.NewFreeDaysService()
	logger.Init()
}

func main() {
	defer logger.CloseLogFile()

	csvReader := csvprocessor.NewCSVReader()

	data, err := csvReader.ReadCSV(context.Background(), filePath)
	if err != nil {
		log.Fatal("Error reading CSV:", err)
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
	log.Printf("Processing --++--> MemberId : %s, Days T0 Add : %s\n", row.MemberId, row.DaysToAdd)

	isSuccessExpiredDate, isSuccessTypeChange, err := freeDaysService.AddFreeDays(context.Background(), row.ExecutionID, row.MemberId, row.DaysToAdd, row.Logs)
	if err != nil {
		log.Println("Error:", err)
		logger.LogError(fmt.Sprintf("Member: (%s) failed", row.MemberId))
		return err
	}

	isSuccessUpdate := fmt.Sprintf("%v", isSuccessExpiredDate)
	isSuccessTypeChangeStr := fmt.Sprintf("%v", isSuccessTypeChange)

	logger.LogCustom(row.MemberId, row.DaysToAdd, isSuccessUpdate, isSuccessTypeChangeStr, err)
	return nil
}

func mapToCSVRow(data [][]string) []CSVRow {
	var result []CSVRow

	for _, row := range data {
		csvRow := CSVRow{
			ExecutionID: row[0],
			MemberId:    row[1],
			DaysToAdd:   row[2],
			Logs:        row[3],
		}

		result = append(result, csvRow)
	}

	return result
}
