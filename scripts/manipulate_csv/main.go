package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

func main() {
	// Create the CSV file
	file, err := os.Create("dummy_logs.csv")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header
	header := []string{"memberId", "daysToAdd", "logsMembership"}
	if err := writer.Write(header); err != nil {
		fmt.Println("Error writing header to CSV:", err)
		return
	}

	// Write the data
	for i := 0; i < 1000; i++ {
		memberId := "TDUMMY" + strconv.Itoa(i)
		daysToAdd := strconv.Itoa(i + 1)
		logsMembership := "log_" + memberId + "_" + daysToAdd
		row := []string{memberId, daysToAdd, logsMembership}

		if err := writer.Write(row); err != nil {
			fmt.Println("Error writing row to CSV:", err)
			return
		}
	}

	fmt.Println("Dummy logs generated and saved to dummy_logs.csv")
}
