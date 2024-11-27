package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"go-scripting/pkg/logger"
)

const (
	filePath       = "dummy_data.csv"
	apiURL         = "http://localhost:8080/v1/referrals/apply"
	maxConcurrency = 2
	maxRetries     = 3
	httpTimeout    = 30 * time.Second
)

type CSVRow struct {
	ReferralType        string
	MemberID            string
	PhoneNumberReferrer string
	MemberIDReferee     string
	PhoneNumberReferee  string
	TrxIDNJM            string
	ClubNJM             string
	FreeDays            int
	ExecutionID         string
	Logs                string
}

type APIPayload struct {
	ReferrerPhone    string    `json:"referrer_phone"`
	ReferrerMemberID string    `json:"referrer_member_id"`
	RefereeDetails   []Referee `json:"referee_detail"`
	FreeDays         int       `json:"free_days"`
}

type Referee struct {
	RefereePhone  string `json:"referee_phone"`
	Location      string `json:"location"`
	TransactionID string `json:"transaction_id"`
}

var (
	headerOnce sync.Once
	httpClient = &http.Client{Timeout: httpTimeout} // Custom HTTP client dengan timeout
)

func init() {
	logger.Init()
}

func initializeCSVHeader() {
	headerOnce.Do(func() {
		header := []string{"timestamp", "MemberID", "FreeDays", "isSuccessUpdate", "isTypeChange", "Error", "APIStatus"}
		if err := logger.WriteHeader(header); err != nil {
			log.Fatalf("Error writing CSV header: %v", err)
		}
	})
}

func main() {
	defer logger.CloseLogFile()

	initializeCSVHeader()

	data, err := readCSV(filePath)
	if err != nil {
		log.Fatalf("Failed to read CSV: %v", err)
	}

	groupedData := groupByMemberID(data)
	semaphore := make(chan struct{}, maxConcurrency)
	var wg sync.WaitGroup

	for memberID, rows := range groupedData {
		semaphore <- struct{}{}
		wg.Add(1)

		go func(memberID string, rows []CSVRow) {
			defer wg.Done()
			defer func() { <-semaphore }()
			if err := processGroupedRows(memberID, rows); err != nil {
				log.Printf("Error processing MemberID %s: %v", memberID, err)
			}
		}(memberID, rows)
	}

	wg.Wait()
	log.Println("Processing completed.")
}

func readCSV(filePath string) ([]CSVRow, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rawData, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var rows []CSVRow
	for _, row := range rawData[1:] {
		freeDays, _ := strconv.Atoi(row[7])
		rows = append(rows, CSVRow{
			ReferralType:        row[0],
			MemberID:            row[1],
			PhoneNumberReferrer: row[2],
			MemberIDReferee:     row[3],
			PhoneNumberReferee:  row[4],
			TrxIDNJM:            row[5],
			ClubNJM:             row[6],
			FreeDays:            freeDays,
			ExecutionID:         row[8],
			Logs:                row[9],
		})
	}
	return rows, nil
}

func groupByMemberID(rows []CSVRow) map[string][]CSVRow {
	grouped := make(map[string][]CSVRow)
	for _, row := range rows {
		grouped[row.MemberID] = append(grouped[row.MemberID], row)
	}
	return grouped
}

func processGroupedRows(memberID string, rows []CSVRow) error {
	expectedFreeDays := rows[0].FreeDays
	referrerPhone := rows[0].PhoneNumberReferrer
	referrerMemberID := rows[0].MemberID

	var refereeDetails []Referee
	for _, row := range rows {
		if row.FreeDays != expectedFreeDays {
			err := fmt.Errorf("inconsistent FreeDays for MemberID %s", memberID)
			logger.LogCustomRefferal(memberID, strconv.Itoa(expectedFreeDays), "false", "false", err, "Inconsistent FreeDays")
			return err
		}

		if row.ReferralType == "NJM Referral Program" {
			refereeDetails = append(refereeDetails, Referee{
				RefereePhone:  row.PhoneNumberReferee,
				Location:      row.ClubNJM,
				TransactionID: row.TrxIDNJM,
			})
		}
	}

	if len(refereeDetails) == 0 {
		logger.LogCustomRefferal(memberID, strconv.Itoa(expectedFreeDays), "false", "false", nil, "No referees for API call")
		return nil
	}

	payload := APIPayload{
		ReferrerPhone:    referrerPhone,
		ReferrerMemberID: referrerMemberID,
		FreeDays:         expectedFreeDays,
		RefereeDetails:   refereeDetails,
	}

	return callAPIWithRetry(payload, maxRetries, memberID)
}

func callAPIWithRetry(payload APIPayload, retries int, memberID string) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	for i := 0; i < retries; i++ {
		resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(payloadBytes))
		fmt.Println("STATUS CODE", resp.StatusCode)
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			// Log success when API responds with 200 status
			logger.LogCustomRefferal(memberID, strconv.Itoa(payload.FreeDays), "true", "true", nil, "API Call Successful with Status 200")
			fmt.Printf("API call for MemberID %s was successful with status 200.\n", memberID)
			return nil
		}

		if resp != nil {
			resp.Body.Close()
			// Log failure with status code
			logger.LogCustomRefferal(memberID, strconv.Itoa(payload.FreeDays), "false", "false", nil, fmt.Sprintf("API Request Failed with Status: %d", resp.StatusCode))
		} else {
			// Log failure without response
			logger.LogCustomRefferal(memberID, strconv.Itoa(payload.FreeDays), "false", "false", err, "API Request Error")
		}

		time.Sleep(time.Second * time.Duration(2<<i)) // Eksponensial Backoff
	}

	return errors.New("API call failed after retries")
}
