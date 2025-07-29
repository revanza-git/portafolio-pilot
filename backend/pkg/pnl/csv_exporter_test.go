package pnl

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/defi-dashboard/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCSVExporter_ExportToWriter(t *testing.T) {
	exporter := NewCSVExporter("/tmp")
	
	data := []models.PnLExportData{
		{
			WalletAddress:     "0x1234567890123456789012345678901234567890",
			TokenSymbol:       "ETH",
			TokenAddress:      "0x0000000000000000000000000000000000000000",
			TransactionHash:   "0xabcdef1234567890",
			Type:              "buy",
			Quantity:          "1.5",
			PriceUSD:          "2000.00",
			RemainingQuantity: "1.5",
			RealizedPnLUSD:    "0",
			Timestamp:         time.Date(2023, 1, 15, 12, 0, 0, 0, time.UTC),
			BlockNumber:       12345,
		},
		{
			WalletAddress:     "0x1234567890123456789012345678901234567890",
			TokenSymbol:       "ETH",
			TokenAddress:      "0x0000000000000000000000000000000000000000",
			TransactionHash:   "0x1234567890abcdef",
			Type:              "sell",
			Quantity:          "0.5",
			PriceUSD:          "2200.00",
			RemainingQuantity: "0.5",
			RealizedPnLUSD:    "100.00",
			Timestamp:         time.Date(2023, 1, 20, 14, 30, 0, 0, time.UTC),
			BlockNumber:       12567,
		},
	}

	var buffer bytes.Buffer
	err := exporter.ExportToWriter(&buffer, data)
	require.NoError(t, err)

	output := buffer.String()
	
	// Check header
	assert.Contains(t, output, "Wallet Address,Token Symbol,Token Address,Transaction Hash,Type,Quantity,Price USD,Remaining Quantity,Realized PnL USD,Timestamp,Block Number")
	
	// Check first data row
	assert.Contains(t, output, "0x1234567890123456789012345678901234567890,ETH,0x0000000000000000000000000000000000000000,0xabcdef1234567890,buy,1.5,2000.00,1.5,0,2023-01-15 12:00:00,12345")
	
	// Check second data row
	assert.Contains(t, output, "0x1234567890123456789012345678901234567890,ETH,0x0000000000000000000000000000000000000000,0x1234567890abcdef,sell,0.5,2200.00,0.5,100.00,2023-01-20 14:30:00,12567")
	
	// Verify CSV structure
	lines := strings.Split(strings.TrimSpace(output), "\n")
	assert.Len(t, lines, 3) // Header + 2 data rows
}

func TestCSVExporter_ExportToCSV(t *testing.T) {
	// Create temporary directory for test
	tempDir := t.TempDir()
	exporter := NewCSVExporter(tempDir)
	
	data := []models.PnLExportData{
		{
			WalletAddress:     "0x1234567890123456789012345678901234567890",
			TokenSymbol:       "BTC",
			TokenAddress:      "0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599",
			TransactionHash:   "0xabcdef1234567890",
			Type:              "buy",
			Quantity:          "0.1",
			PriceUSD:          "45000.00",
			RemainingQuantity: "0.1",
			RealizedPnLUSD:    "0",
			Timestamp:         time.Date(2023, 1, 15, 12, 0, 0, 0, time.UTC),
			BlockNumber:       12345,
		},
	}

	// Export to CSV
	filePath, err := exporter.ExportToCSV(data, "0x12345678")
	require.NoError(t, err)
	
	// Verify file was created
	assert.True(t, strings.HasPrefix(filepath.Base(filePath), "pnl_export_0x123456"))
	assert.True(t, strings.HasSuffix(filePath, ".csv"))
	
	// Check file exists
	_, err = os.Stat(filePath)
	require.NoError(t, err)
	
	// Read file content
	content, err := os.ReadFile(filePath)
	require.NoError(t, err)
	
	contentStr := string(content)
	assert.Contains(t, contentStr, "Wallet Address,Token Symbol")
	assert.Contains(t, contentStr, "0x1234567890123456789012345678901234567890,BTC")
	
	// Clean up
	os.Remove(filePath)
}

func TestCSVExporter_CleanupFile(t *testing.T) {
	// Create temporary directory and file
	tempDir := t.TempDir()
	exporter := NewCSVExporter(tempDir)
	
	// Create a test file
	testFile := filepath.Join(tempDir, "test_file.csv")
	err := os.WriteFile(testFile, []byte("test content"), 0644)
	require.NoError(t, err)
	
	// Verify file exists
	_, err = os.Stat(testFile)
	require.NoError(t, err)
	
	// Clean up file
	err = exporter.CleanupFile(testFile)
	require.NoError(t, err)
	
	// Verify file is gone
	_, err = os.Stat(testFile)
	assert.True(t, os.IsNotExist(err))
}

func TestCSVExporter_ScheduleCleanup(t *testing.T) {
	// Create temporary directory and file
	tempDir := t.TempDir()
	exporter := NewCSVExporter(tempDir)
	
	// Create a test file
	testFile := filepath.Join(tempDir, "test_cleanup.csv")
	err := os.WriteFile(testFile, []byte("test content"), 0644)
	require.NoError(t, err)
	
	// Verify file exists
	_, err = os.Stat(testFile)
	require.NoError(t, err)
	
	// Schedule cleanup with very short duration
	exporter.ScheduleCleanup(testFile, 10*time.Millisecond)
	
	// Wait for cleanup to happen
	time.Sleep(50 * time.Millisecond)
	
	// Verify file is gone
	_, err = os.Stat(testFile)
	assert.True(t, os.IsNotExist(err))
}

func TestCSVExporter_EmptyData(t *testing.T) {
	exporter := NewCSVExporter("/tmp")
	
	var buffer bytes.Buffer
	err := exporter.ExportToWriter(&buffer, []models.PnLExportData{})
	require.NoError(t, err)

	output := buffer.String()
	
	// Should only contain header
	lines := strings.Split(strings.TrimSpace(output), "\n")
	assert.Len(t, lines, 1) // Only header
	assert.Contains(t, output, "Wallet Address,Token Symbol,Token Address")
}

func TestCSVExporter_SpecialCharacters(t *testing.T) {
	exporter := NewCSVExporter("/tmp")
	
	data := []models.PnLExportData{
		{
			WalletAddress:     "0x1234567890123456789012345678901234567890",
			TokenSymbol:       "TOKEN,WITH,COMMAS",
			TokenAddress:      "0x0000000000000000000000000000000000000000",
			TransactionHash:   "0xabcdef1234567890",
			Type:              "buy",
			Quantity:          "1.5",
			PriceUSD:          "2000.00",
			RemainingQuantity: "1.5",
			RealizedPnLUSD:    "0",
			Timestamp:         time.Date(2023, 1, 15, 12, 0, 0, 0, time.UTC),
			BlockNumber:       12345,
		},
	}

	var buffer bytes.Buffer
	err := exporter.ExportToWriter(&buffer, data)
	require.NoError(t, err)

	output := buffer.String()
	
	// CSV should properly escape commas in token symbol
	assert.Contains(t, output, "\"TOKEN,WITH,COMMAS\"")
}