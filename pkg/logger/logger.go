package logger

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	writer     *csv.Writer
	logFile    *os.File
	once       sync.Once
	mutex      sync.Mutex
	closeMutex sync.Mutex
)

// Init initializes the logger by setting up the CSV writer.
func Init() {
	once.Do(func() {
		if err := os.MkdirAll("logs", os.ModePerm); err != nil {
			log.Fatal("Error creating logs directory:", err)
		}

		logFilePath := generateFileNameWithTimestamp("logs", "app")
		var err error
		logFile, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal("Error opening log file:", err)
		}

		writer = csv.NewWriter(logFile)
	})
}

func WriteHeader(header []string) error {
	mutex.Lock()
	defer mutex.Unlock()

	Init() // logger is initialized before writing

	// Write the header to the CSV file
	if err := writer.Write(header); err != nil {
		log.Println("Error writing header to CSV:", err)
		return err
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		log.Println("Error flushing writer after writing header:", err)
		return err
	}

	return nil
}

func LogInfo(message string, details ...interface{}) {
	writeCSV("INFO", message, details...)
}

func LogError(message string, details ...interface{}) {
	writeCSV("ERROR", message, details...)
}

func LogSalesBy(transactionId, salesBy, status, errorDetail string) {
	writeCSV(transactionId, salesBy, status, errorDetail)
}

func LogCustom(memberId, daysToAdd, isSuccessUpdatedExpired, isSuccessTypeChange string, err error) {
	writeCSV(memberId, daysToAdd, isSuccessUpdatedExpired, isSuccessTypeChange, err)
}

func LogGo(trxId, executionId, executionType, startDate, endDate, membershipLogs, isSuccessUpdateExpiredMembership, isSuccessTypeChange, isSuccessUpdateV1, isSuccessUpdateV6, errorStatus string) {
	writeCSV(trxId, executionId, startDate, endDate, membershipLogs, isSuccessUpdateExpiredMembership, isSuccessTypeChange, isSuccessUpdateV1, isSuccessUpdateV6, errorStatus)
}

func LogCustomRefferal(memberId, freeDays, isSuccessUpdate, isSuccessTypeChangeStr string, err error, apiStatus string) {
	writeCSV(memberId, freeDays, isSuccessUpdate, isSuccessTypeChangeStr, err, apiStatus)
}

func writeCSV(level, message string, details ...interface{}) {
	mutex.Lock()
	defer mutex.Unlock()

	Init() // logger is initialized before writing

	timestamp := time.Now().Format(time.RFC3339)
	record := append([]string{timestamp, level, message}, toStringSlice(details)...)
	if err := writer.Write(record); err != nil {
		log.Println("Error writing to CSV:", err)
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		log.Println("Error flushing writer:", err)
	}
}

// toStringSlice converts a slice of interface{} to a slice of string
func toStringSlice(details []interface{}) []string {
	result := make([]string, len(details))
	for i, detail := range details {
		result[i] = fmt.Sprint(detail)
	}
	return result
}

// CloseLogFile closes the log file. Should be called when the application is done logging.
func CloseLogFile() {
	closeMutex.Lock()
	defer closeMutex.Unlock()

	if logFile != nil {
		writer.Flush()
		if err := writer.Error(); err != nil {
			log.Println("Error flushing writer:", err)
		}
		if err := logFile.Close(); err != nil {
			log.Println("Error closing log file:", err)
		}
		logFile = nil // Prevents multiple closes
	}
}

func generateFileNameWithTimestamp(basePath, prefix string) string {
	currentTime := time.Now().Format("20060102_150405") // Format: YYYYMMDD_HHMMSS
	return filepath.Join(basePath, fmt.Sprintf("%s_%s.csv", prefix, currentTime))
}
