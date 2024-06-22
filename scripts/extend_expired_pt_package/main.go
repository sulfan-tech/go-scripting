package main

import (
	"context"
	csvprocessor "go-scripting/pkg/csv_processor"
	"go-scripting/script/extend_expired_pt_package/service"
	"log"
	"sync"

	"github.com/joho/godotenv"
)

const (
	filePath   = "input.csv"
	numWorkers = 100 // Number of concurrent workers
)

type CSVRow struct {
	TransactionID     string
	MemberPhoneNumber string
	NewExpiredDate    string
	Type              string
	ExtendStatus      string
}

var extendMemberPTPackageService service.ExtendExpiredPTInterface

func init() {
	env := "../../.env"
	err := godotenv.Load(env)
	if err != nil {
		log.Fatal("Error loading .env file:", err.Error())
	}

	extendMemberPTPackageService = service.NewExtendExpiredPTService()
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
			if row.Type != "Extend" {
				log.Println("SKIP Unsupported Type")
				continue
			}
			row.MemberPhoneNumber = sanitizePhoneNumber(row.MemberPhoneNumber)
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

	req := service.RequestExtendMemberPTPackage{
		TransactionID:     row.TransactionID,
		MemberPhoneNumber: row.MemberPhoneNumber,
		NewExpiredDate:    row.NewExpiredDate,
		Type:              row.Type,
	}

	log.Printf("Processing - Member Phone: %s\n", row.MemberPhoneNumber)

	err := extendMemberPTPackageService.ExtendExpiredPT(ctx, req)
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
			TransactionID:     row[0],
			MemberPhoneNumber: row[1],
			NewExpiredDate:    row[2],
			Type:              row[3],
		}

		result = append(result, csvRow)
	}

	return result
}

func sanitizePhoneNumber(phoneNumber string) string {
	// Assuming that phoneNumber is always in the format "6281213390614"
	// You can modify this function based on the actual structure of your phone numbers
	return "+6" + phoneNumber[1:]
}
