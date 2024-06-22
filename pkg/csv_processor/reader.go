package csvreader

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
)

type CSVReader struct{}

type CSVReaderInterface interface {
	ReadCSV(ctx context.Context, filePath string) ([][]string, error)
}

func NewCSVReader() *CSVReader {
	return &CSVReader{}
}

// ReadCSV reads the CSV file from the given file path, ignoring the header row.
func (r *CSVReader) ReadCSV(ctx context.Context, filePath string) ([][]string, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read and discard the header row
	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("error reading header: %w", err)
	}

	// Read all remaining rows
	data, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV: %w", err)
	}

	return data, nil
}
