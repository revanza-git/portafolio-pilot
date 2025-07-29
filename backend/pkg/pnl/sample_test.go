package pnl

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/defi-dashboard/backend/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSampleDataset tests the PnL calculation with a realistic dataset
func TestSampleDataset(t *testing.T) {
	calculator := NewCalculator(FIFO)

	// Create sample data that simulates real trading activity
	lots := createSampleTradingData()

	currentPrice := "1500.00" // Current ETH price

	result, err := calculator.CalculatePnL(lots, currentPrice)
	require.NoError(t, err)

	// Verify the calculation makes sense
	assert.Equal(t, "fifo", result.Method)
	assert.NotEmpty(t, result.RealizedPnLUSD)
	assert.NotEmpty(t, result.UnrealizedPnLUSD)
	assert.NotEmpty(t, result.TotalPnLUSD)
	assert.NotEmpty(t, result.CurrentQuantity)

	t.Logf("Sample Dataset Results:")
	t.Logf("Realized PnL: %s USD", result.RealizedPnLUSD)
	t.Logf("Unrealized PnL: %s USD", result.UnrealizedPnLUSD)
	t.Logf("Total PnL: %s USD", result.TotalPnLUSD)
	t.Logf("Current Quantity: %s", result.CurrentQuantity)
	t.Logf("Current Value: %s USD", result.CurrentValueUSD)
	t.Logf("Cost Basis: %s USD", result.TotalCostBasisUSD)
}

// TestCSVExportWithSampleData tests CSV export with realistic data
func TestCSVExportWithSampleData(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()
	exporter := NewCSVExporter(tempDir)

	// Create sample export data
	exportData := createSampleExportData()

	// Test file export
	filePath, err := exporter.ExportToCSV(exportData, "0x1234abcd")
	require.NoError(t, err)

	t.Logf("CSV file created at: %s", filePath)

	// Verify file can be opened (this would work in Excel)
	content, err := ReadCSVFile(filePath)
	require.NoError(t, err)
	
	assert.True(t, len(content) > 1) // Header + data rows
	assert.Contains(t, content[0], "Wallet Address") // Header check

	t.Logf("CSV contains %d rows (including header)", len(content))
	
	// Print first few rows for manual verification
	for i, row := range content {
		if i < 3 { // Print header + first 2 data rows
			t.Logf("Row %d: %s", i, row)
		}
	}

	// Clean up
	exporter.CleanupFile(filePath)
}

// createSampleTradingData creates realistic trading data for testing
func createSampleTradingData() []models.PnLLot {
	baseTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	
	return []models.PnLLot{
		// Initial buy
		{
			ID:                uuid.New(),
			Type:              "buy",
			Quantity:          "2.0",
			PriceUSD:          "1200.00",
			RemainingQuantity: "2.0",
			Timestamp:         baseTime,
			BlockNumber:       12000,
			TransactionHash:   "0x1111111111111111111111111111111111111111111111111111111111111111",
		},
		// Second buy at higher price
		{
			ID:                uuid.New(),
			Type:              "buy",
			Quantity:          "1.5",
			PriceUSD:          "1400.00",
			RemainingQuantity: "1.5",
			Timestamp:         baseTime.Add(7 * 24 * time.Hour),
			BlockNumber:       12100,
			TransactionHash:   "0x2222222222222222222222222222222222222222222222222222222222222222",
		},
		// Partial sell
		{
			ID:                uuid.New(),
			Type:              "sell",
			Quantity:          "1.0",
			PriceUSD:          "1600.00",
			RemainingQuantity: "1.0",
			Timestamp:         baseTime.Add(14 * 24 * time.Hour),
			BlockNumber:       12200,
			TransactionHash:   "0x3333333333333333333333333333333333333333333333333333333333333333",
		},
		// Another buy during dip
		{
			ID:                uuid.New(),
			Type:              "buy",
			Quantity:          "0.8",
			PriceUSD:          "1100.00",
			RemainingQuantity: "0.8",
			Timestamp:         baseTime.Add(21 * 24 * time.Hour),
			BlockNumber:       12300,
			TransactionHash:   "0x4444444444444444444444444444444444444444444444444444444444444444",
		},
		// Final sell
		{
			ID:                uuid.New(),
			Type:              "sell",
			Quantity:          "1.5",
			PriceUSD:          "1800.00",
			RemainingQuantity: "1.5",
			Timestamp:         baseTime.Add(30 * 24 * time.Hour),
			BlockNumber:       12400,
			TransactionHash:   "0x5555555555555555555555555555555555555555555555555555555555555555",
		},
	}
}

// createSampleExportData creates sample data for CSV export testing
func createSampleExportData() []models.PnLExportData {
	baseTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	
	return []models.PnLExportData{
		{
			WalletAddress:     "0x1234567890123456789012345678901234567890",
			TokenSymbol:       "ETH",
			TokenAddress:      "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
			TransactionHash:   "0x1111111111111111111111111111111111111111111111111111111111111111",
			Type:              "buy",
			Quantity:          "2.0",
			PriceUSD:          "1200.00",
			RemainingQuantity: "1.0", // After partial sell
			RealizedPnLUSD:    "400.00",
			Timestamp:         baseTime,
			BlockNumber:       12000,
		},
		{
			WalletAddress:     "0x1234567890123456789012345678901234567890",
			TokenSymbol:       "ETH",
			TokenAddress:      "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
			TransactionHash:   "0x2222222222222222222222222222222222222222222222222222222222222222",
			Type:              "buy",
			Quantity:          "1.5",
			PriceUSD:          "1400.00",
			RemainingQuantity: "1.5",
			RealizedPnLUSD:    "0.00",
			Timestamp:         baseTime.Add(7 * 24 * time.Hour),
			BlockNumber:       12100,
		},
		{
			WalletAddress:     "0x1234567890123456789012345678901234567890",
			TokenSymbol:       "ETH",
			TokenAddress:      "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
			TransactionHash:   "0x3333333333333333333333333333333333333333333333333333333333333333",
			Type:              "sell",
			Quantity:          "1.0",
			PriceUSD:          "1600.00",
			RemainingQuantity: "0.0",
			RealizedPnLUSD:    "400.00",
			Timestamp:         baseTime.Add(14 * 24 * time.Hour),
			BlockNumber:       12200,
		},
	}
}

// ReadCSVFile is a helper function to read and parse CSV for testing
func ReadCSVFile(filepath string) ([]string, error) {
	// This is a simplified version - in real tests you'd use csv.Reader
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	return lines, nil
}

// TestFIFOvsLIFOComparison demonstrates the difference between FIFO and LIFO
func TestFIFOvsLIFOComparison(t *testing.T) {
	lots := createVolatileTradingData()
	currentPrice := "1500.00"

	// Test FIFO
	fifoCalc := NewCalculator(FIFO)
	fifoResult, err := fifoCalc.CalculatePnL(lots, currentPrice)
	require.NoError(t, err)

	// Test LIFO
	lifoCalc := NewCalculator(LIFO)
	lifoResult, err := lifoCalc.CalculatePnL(lots, currentPrice)
	require.NoError(t, err)

	t.Logf("FIFO vs LIFO Comparison:")
	t.Logf("FIFO Realized PnL: %s USD", fifoResult.RealizedPnLUSD)
	t.Logf("LIFO Realized PnL: %s USD", lifoResult.RealizedPnLUSD)
	t.Logf("FIFO Total PnL: %s USD", fifoResult.TotalPnLUSD)
	t.Logf("LIFO Total PnL: %s USD", lifoResult.TotalPnLUSD)

	// The results should be different for FIFO vs LIFO
	assert.NotEqual(t, fifoResult.RealizedPnLUSD, lifoResult.RealizedPnLUSD)
}

// createVolatileTradingData creates data where FIFO and LIFO will show different results
func createVolatileTradingData() []models.PnLLot {
	baseTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	
	return []models.PnLLot{
		// Buy low
		{
			ID:                uuid.New(),
			Type:              "buy",
			Quantity:          "1.0",
			PriceUSD:          "1000.00",
			RemainingQuantity: "1.0",
			Timestamp:         baseTime,
		},
		// Buy high
		{
			ID:                uuid.New(),
			Type:              "buy",
			Quantity:          "1.0",
			PriceUSD:          "2000.00",
			RemainingQuantity: "1.0",
			Timestamp:         baseTime.Add(time.Hour),
		},
		// Sell one unit (FIFO will use $1000 cost, LIFO will use $2000 cost)
		{
			ID:                uuid.New(),
			Type:              "sell",
			Quantity:          "1.0",
			PriceUSD:          "1500.00",
			RemainingQuantity: "1.0",
			Timestamp:         baseTime.Add(2 * time.Hour),
		},
	}
}