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
		writeHeader() // Write header once when initialized
	})
}

func writeHeader() {
	headers := []string{"timestamp", "memberId", "daysToAdd", "isSuccessUpdateExpiredMembership", "isSuccessTypeChange"}
	if err := writer.Write(headers); err != nil {
		log.Println("Error writing header to CSV:", err)
	}
	writer.Flush()
}

// LogInfo logs an informational message.
func LogInfo(message string, details ...interface{}) {
	writeCSV("INFO", message, details...)
}

func LogCustom(memberId, daysToAdd, isSuccessUpdatedExpired, isSuccessTypeChange string) {
	writeCSV(memberId, daysToAdd, isSuccessUpdatedExpired, isSuccessTypeChange)
}

// LogError logs an error message.
func LogError(message string, details ...interface{}) {
	writeCSV("ERROR", message, details...)
}

func writeCSV(level, message string, details ...interface{}) {
	mutex.Lock()
	defer mutex.Unlock()

	Init() // Ensure logger is initialized before writing

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
	var result []string
	for _, detail := range details {
		result = append(result, fmt.Sprint(detail))
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
	fileName := fmt.Sprintf("%s_%s.csv", prefix, currentTime)
	return filepath.Join(basePath, fileName)
}
