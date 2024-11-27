package main

import (
	"context"
	csvprocessor "go-scripting/pkg/csv_processor"
	"go-scripting/scripts/change_salesBy/service"

	"log"

	"github.com/joho/godotenv"
)

const (
	filePath = "Untitled spreadsheet - Sheet1 (5).csv"
)

type CSVRow struct {
	TransactionID string
	NewSalesBy    string
}

var changeSalesByService service.ChangeSalesByInterface

func init() {
	env := "../../.env"
	err := godotenv.Load(env)
	if err != nil {
		log.Fatal("Error loading .env file:", err.Error())
	}

	changeSalesByService = service.NewChangeSalesByService()
}

func main() {
	csvReader := csvprocessor.NewCSVReader()

	data, err := csvReader.ReadCSV(context.Background(), filePath)
	if err != nil {
		log.Fatal("Error reading CSV:", err)
		return
	}

	mappedData := mapToCSVRow(data)

	for i, row := range mappedData {
		err := processRow(row)
		if err != nil {
			log.Printf("Error processing row %d: %v", i+1, err)
		}
	}

	log.Println("Processing completed.")
}

func processRow(row CSVRow) error {
	ctx := context.Background()

	log.Printf("Processing - TransactionId: %s\n", row.TransactionID)
	err := changeSalesByService.ChangeSalesBy(ctx, row.TransactionID, row.NewSalesBy)
	if err != nil {
		log.Printf("Error processing row - TransactionID: %s, Error: %s", row.TransactionID, err.Error())
		return err
	}

	return nil
}

func mapToCSVRow(data [][]string) []CSVRow {
	var result []CSVRow

	for _, row := range data {
		csvRow := CSVRow{
			TransactionID: row[0],
			NewSalesBy:    row[1],
		}

		result = append(result, csvRow)
	}

	return result
}
