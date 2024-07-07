package main

import (
	"context"
	"fmt"
	"log"

	csvprocessor "go-scripting/pkg/csv_processor"
	"go-scripting/pkg/logger"
	"go-scripting/scripts/bulk_update_homeclub/service"

	"github.com/joho/godotenv"
)

const (
	filePath    = "homeclub.csv"
	envFilePath = "../../.env"
)

type CSVRow struct {
	Phone       string
	NewHomeClub string
}

var changeHomeClubService service.ChangeHomeClubInterface

func init() {
	loadEnv()
	initializeService()
}

func main() {
	csvReader := csvprocessor.NewCSVReader()

	data, err := csvReader.ReadCSV(context.Background(), filePath)
	if err != nil {
		logger.LogError(fmt.Sprintf("Error reading CSV: %v", err))
	}

	processCSVData(data)
}

func loadEnv() {
	if err := godotenv.Load(envFilePath); err != nil {
		logger.LogError(fmt.Sprintf("Error loading .env file: %v", err))
	}
}

func initializeService() {
	changeHomeClubService = service.NewChangeHomeClubService()
}

func processCSVData(data [][]string) {
	rows := mapToCSVRows(data)
	for rowNumber, row := range rows {
		if err := processRow(row, rowNumber+1); err != nil {
			logger.LogError(err.Error())
			continue
		}
	}
}

func processRow(row CSVRow, rowCount int) error {
	log.Printf("Processing row %d - Phone: %s, NewHomeClub: %s", rowCount, row.Phone, row.NewHomeClub)

	err := changeHomeClubService.BulkChangeHomeClub(context.Background(), row.Phone, row.NewHomeClub)
	if err != nil {
		errorMsg := fmt.Sprintf("Member (%s) failed on Row (%d): %v", row.Phone, rowCount, err)
		logger.LogError(errorMsg)
		return fmt.Errorf("operation failed on Row %d: %v", rowCount, err)
	}

	successMsg := fmt.Sprintf("ChangeHomeClub (%v) on Member (%s) updated successfully on Row (%v)", row.NewHomeClub, row.Phone, rowCount)
	logger.LogInfo(successMsg)
	return nil
}

func mapToCSVRows(data [][]string) []CSVRow {
	var result []CSVRow
	for _, row := range data {
		result = append(result, CSVRow{
			Phone:       row[0],
			NewHomeClub: row[1],
		})
	}
	return result
}
