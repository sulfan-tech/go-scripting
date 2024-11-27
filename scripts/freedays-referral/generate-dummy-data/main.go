// generate_csv.go
package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

func main() {
	file, err := os.Create("dummy_data.csv")
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"ReferralType", "MemberID", "PhoneNumberReferrer", "MemberIDReferee", "PhoneNumberReferee", "TrxIDNJM", "ClubNJM", "FreeDays", "ExecutionID", "Logs"}
	if err := writer.Write(header); err != nil {
		log.Fatalf("Failed to write header: %v", err)
	}

	// Generate dummy rows
	for i := 1; i <= 1000; i++ {
		row := []string{
			"NJM Referral Program",
			"Member" + strconv.Itoa(i),
			"08123456789",
			"RefereeMember" + strconv.Itoa(i),
			"08234567890",
			"TrxID" + strconv.Itoa(i),
			"Club" + strconv.Itoa(i),
			strconv.Itoa(5),
			"Exec" + strconv.Itoa(i),
			"Log" + strconv.Itoa(i),
		}
		if err := writer.Write(row); err != nil {
			log.Fatalf("Failed to write row %d: %v", i, err)
		}
	}

	log.Println("Dummy CSV data generated successfully!")
}
