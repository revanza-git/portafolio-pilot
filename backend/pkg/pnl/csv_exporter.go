package pnl

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/defi-dashboard/backend/internal/models"
)

type CSVExporter struct {
	tempDir string
}

func NewCSVExporter(tempDir string) *CSVExporter {
	return &CSVExporter{tempDir: tempDir}
}

// ExportToCSV exports PnL data to a CSV file and returns the file path
func (e *CSVExporter) ExportToCSV(data []models.PnLExportData, walletAddress string) (string, error) {
	// Create temporary file
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("pnl_export_%s_%s.csv", walletAddress[:8], timestamp)
	filepath := filepath.Join(e.tempDir, filename)

	file, err := os.Create(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	// Create CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"Wallet Address",
		"Token Symbol",
		"Token Address",
		"Transaction Hash",
		"Type",
		"Quantity",
		"Price USD",
		"Remaining Quantity",
		"Realized PnL USD",
		"Timestamp",
		"Block Number",
	}

	if err := writer.Write(header); err != nil {
		return "", fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data rows
	for _, row := range data {
		record := []string{
			row.WalletAddress,
			row.TokenSymbol,
			row.TokenAddress,
			row.TransactionHash,
			row.Type,
			row.Quantity,
			row.PriceUSD,
			row.RemainingQuantity,
			row.RealizedPnLUSD,
			row.Timestamp.Format("2006-01-02 15:04:05"),
			strconv.FormatInt(row.BlockNumber, 10),
		}

		if err := writer.Write(record); err != nil {
			return "", fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	return filepath, nil
}

// ExportToWriter exports PnL data directly to a writer (for streaming)
func (e *CSVExporter) ExportToWriter(writer io.Writer, data []models.PnLExportData) error {
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	// Write header
	header := []string{
		"Wallet Address",
		"Token Symbol",
		"Token Address",
		"Transaction Hash",
		"Type",
		"Quantity",
		"Price USD",
		"Remaining Quantity",
		"Realized PnL USD",
		"Timestamp",
		"Block Number",
	}

	if err := csvWriter.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data rows
	for _, row := range data {
		record := []string{
			row.WalletAddress,
			row.TokenSymbol,
			row.TokenAddress,
			row.TransactionHash,
			row.Type,
			row.Quantity,
			row.PriceUSD,
			row.RemainingQuantity,
			row.RealizedPnLUSD,
			row.Timestamp.Format("2006-01-02 15:04:05"),
			strconv.FormatInt(row.BlockNumber, 10),
		}

		if err := csvWriter.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	return nil
}

// CleanupFile removes a temporary file
func (e *CSVExporter) CleanupFile(filepath string) error {
	return os.Remove(filepath)
}

// ScheduleCleanup schedules a file for cleanup after the specified duration
func (e *CSVExporter) ScheduleCleanup(filepath string, duration time.Duration) {
	go func() {
		time.Sleep(duration)
		e.CleanupFile(filepath)
	}()
}