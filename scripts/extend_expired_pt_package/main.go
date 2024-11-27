package main

import (
	"context"
	csvprocessor "go-scripting/pkg/csv_processor"
	"go-scripting/scripts/extend_expired_pt_package/service"

	"log"

	"github.com/joho/godotenv"
)

const (
	filePath = "pt-stack - Sheet1.csv"
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

	for i, row := range mappedData {
		if row.Type != "Extend" {
			log.Println("SKIP Unsupported Type")
			continue
		}
		row.MemberPhoneNumber = sanitizePhoneNumber(row.MemberPhoneNumber)

		err := processRow(row)
		if err != nil {
			log.Printf("Error processing row %d: %v", i+1, err)
		}
	}

	log.Println("Processing completed.")
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
	// Format phone numbers starting with "0" to start with "+62"
	if len(phoneNumber) > 0 && phoneNumber[0] == '0' {
		return "+62" + phoneNumber[1:]
	}
	return phoneNumber
}
