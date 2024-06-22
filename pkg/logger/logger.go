package logger

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	writer  *csv.Writer
	logFile *os.File
)

func init() {
	// Create the "logs" directory if it doesn't exist
	if err := os.MkdirAll("logs", os.ModePerm); err != nil {
		log.Fatal(err)
	}

	logFilePath := generateFileNameWithTimestamp("logs", "app")
	var err error
	logFile, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}

	writer = csv.NewWriter(logFile)
}

func LogInfo(message string) {
	writeCSV("INFO", message)
}

func LogError(message string) {
	writeCSV("ERROR", message)
}

func writeCSV(level string, message string) {
	record := []string{level, message}
	if err := writer.Write(record); err != nil {
		log.Println("Error writing to CSV:", err)
	}
	writer.Flush()
}

// CloseLogFile closes the log file. Should be called when the application is done logging.
func CloseLogFile() {
	if logFile != nil {
		writer.Flush()
		if err := logFile.Close(); err != nil {
			log.Println("Error closing log file:", err)
		}
	}
}

func generateFileNameWithTimestamp(basePath, prefix string) string {
	currentTime := time.Now().Format("20060102_150405") // Format: YYYYMMDD_HHMMSS
	fileName := fmt.Sprintf("%s_%s.csv", prefix, currentTime)
	return filepath.Join(basePath, fileName)
}
