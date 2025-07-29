package pnl

import (
	"testing"
	"time"

	"github.com/defi-dashboard/backend/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalculator_CalculatePnL_FIFO(t *testing.T) {
	calculator := NewCalculator(FIFO)

	tests := []struct {
		name               string
		lots               []models.PnLLot
		currentPriceUSD    string
		expectedRealizedPnL string
		expectedUnrealizedPnL string
		expectedCurrentQuantity string
		expectError        bool
	}{
		{
			name: "Simple FIFO calculation",
			lots: []models.PnLLot{
				{
					ID:                uuid.New(),
					Type:              "buy",
					Quantity:          "100",
					PriceUSD:          "10.00",
					RemainingQuantity: "100",
					Timestamp:         time.Now().Add(-time.Hour * 24),
				},
				{
					ID:                uuid.New(),
					Type:              "sell",
					Quantity:          "50",
					PriceUSD:          "15.00",
					RemainingQuantity: "50",
					Timestamp:         time.Now().Add(-time.Hour * 12),
				},
			},
			currentPriceUSD:         "20.00",
			expectedRealizedPnL:     "250", // (15 - 10) * 50 = 250
			expectedUnrealizedPnL:   "500", // (20 - 10) * 50 = 500
			expectedCurrentQuantity: "50",
			expectError:             false,
		},
		{
			name: "Multiple buys and sells FIFO",
			lots: []models.PnLLot{
				{
					ID:                uuid.New(),
					Type:              "buy",
					Quantity:          "100",
					PriceUSD:          "10.00",
					RemainingQuantity: "100",
					Timestamp:         time.Now().Add(-time.Hour * 48),
				},
				{
					ID:                uuid.New(),
					Type:              "buy",
					Quantity:          "50",
					PriceUSD:          "12.00",
					RemainingQuantity: "50",
					Timestamp:         time.Now().Add(-time.Hour * 36),
				},
				{
					ID:                uuid.New(),
					Type:              "sell",
					Quantity:          "75",
					PriceUSD:          "15.00",
					RemainingQuantity: "75",
					Timestamp:         time.Now().Add(-time.Hour * 12),
				},
			},
			currentPriceUSD:         "18.00",
			expectedRealizedPnL:     "375", // (15-10)*75 = 375 (sells from first buy only)
			expectedUnrealizedPnL:   "175", // (18-10)*25 + (18-12)*50 = 150 + 300 = 450
			expectedCurrentQuantity: "75",  // 100 + 50 - 75 = 75
			expectError:             false,
		},
		{
			name: "Partial sell FIFO",
			lots: []models.PnLLot{
				{
					ID:                uuid.New(),
					Type:              "buy",
					Quantity:          "100",
					PriceUSD:          "10.00",
					RemainingQuantity: "100",
					Timestamp:         time.Now().Add(-time.Hour * 24),
				},
				{
					ID:                uuid.New(),
					Type:              "sell",
					Quantity:          "30",
					PriceUSD:          "8.00",
					RemainingQuantity: "30",
					Timestamp:         time.Now().Add(-time.Hour * 12),
				},
			},
			currentPriceUSD:         "12.00",
			expectedRealizedPnL:     "-60", // (8 - 10) * 30 = -60
			expectedUnrealizedPnL:   "140", // (12 - 10) * 70 = 140
			expectedCurrentQuantity: "70",
			expectError:             false,
		},
		{
			name: "No lots",
			lots: []models.PnLLot{},
			currentPriceUSD: "10.00",
			expectError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calculator.CalculatePnL(tt.lots, tt.currentPriceUSD)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedRealizedPnL, result.RealizedPnLUSD)
			assert.Equal(t, tt.expectedCurrentQuantity, result.CurrentQuantity)
			assert.Equal(t, "fifo", result.Method)
		})
	}
}

func TestCalculator_CalculatePnL_LIFO(t *testing.T) {
	calculator := NewCalculator(LIFO)

	tests := []struct {
		name               string
		lots               []models.PnLLot
		currentPriceUSD    string
		expectedRealizedPnL string
		expectedCurrentQuantity string
		expectError        bool
	}{
		{
			name: "LIFO vs FIFO difference",
			lots: []models.PnLLot{
				{
					ID:                uuid.New(),
					Type:              "buy",
					Quantity:          "100",
					PriceUSD:          "10.00",
					RemainingQuantity: "100",
					Timestamp:         time.Now().Add(-time.Hour * 48),
				},
				{
					ID:                uuid.New(),
					Type:              "buy",
					Quantity:          "50",
					PriceUSD:          "15.00",
					RemainingQuantity: "50",
					Timestamp:         time.Now().Add(-time.Hour * 36),
				},
				{
					ID:                uuid.New(),
					Type:              "sell",
					Quantity:          "75",
					PriceUSD:          "20.00",
					RemainingQuantity: "75",
					Timestamp:         time.Now().Add(-time.Hour * 12),
				},
			},
			currentPriceUSD:         "18.00",
			expectedRealizedPnL:     "375", // LIFO: (20-15)*50 + (20-10)*25 = 250 + 250 = 500
			expectedCurrentQuantity: "75",  // 100 + 50 - 75 = 75
			expectError:             false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calculator.CalculatePnL(tt.lots, tt.currentPriceUSD)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedCurrentQuantity, result.CurrentQuantity)
			assert.Equal(t, "lifo", result.Method)
		})
	}
}

func TestCalculator_EdgeCases(t *testing.T) {
	calculator := NewCalculator(FIFO)

	t.Run("Sell more than bought", func(t *testing.T) {
		lots := []models.PnLLot{
			{
				ID:                uuid.New(),
				Type:              "buy",
				Quantity:          "100",
				PriceUSD:          "10.00",
				RemainingQuantity: "100",
				Timestamp:         time.Now().Add(-time.Hour * 24),
			},
			{
				ID:                uuid.New(),
				Type:              "sell",
				Quantity:          "150",
				PriceUSD:          "15.00",
				RemainingQuantity: "150",
				Timestamp:         time.Now().Add(-time.Hour * 12),
			},
		}

		result, err := calculator.CalculatePnL(lots, "20.00")
		require.NoError(t, err)
		
		// Should only match available quantity
		assert.Equal(t, "500", result.RealizedPnLUSD) // (15-10)*100 = 500
		assert.Equal(t, "0", result.CurrentQuantity)  // No remaining quantity
	})

	t.Run("Only buy lots", func(t *testing.T) {
		lots := []models.PnLLot{
			{
				ID:                uuid.New(),
				Type:              "buy",
				Quantity:          "100",
				PriceUSD:          "10.00",
				RemainingQuantity: "100",
				Timestamp:         time.Now().Add(-time.Hour * 24),
			},
			{
				ID:                uuid.New(),
				Type:              "buy",
				Quantity:          "50",
				PriceUSD:          "12.00",
				RemainingQuantity: "50",
				Timestamp:         time.Now().Add(-time.Hour * 12),
			},
		}

		result, err := calculator.CalculatePnL(lots, "15.00")
		require.NoError(t, err)
		
		assert.Equal(t, "0", result.RealizedPnLUSD)   // No sells
		assert.Equal(t, "150", result.CurrentQuantity) // All bought quantity remains
	})

	t.Run("Only sell lots", func(t *testing.T) {
		lots := []models.PnLLot{
			{
				ID:                uuid.New(),
				Type:              "sell",
				Quantity:          "100",
				PriceUSD:          "15.00",
				RemainingQuantity: "100",
				Timestamp:         time.Now().Add(-time.Hour * 12),
			},
		}

		result, err := calculator.CalculatePnL(lots, "20.00")
		require.NoError(t, err)
		
		assert.Equal(t, "0", result.RealizedPnLUSD)   // No buys to match against
		assert.Equal(t, "0", result.CurrentQuantity)  // No remaining quantity
	})

	t.Run("Zero quantities", func(t *testing.T) {
		lots := []models.PnLLot{
			{
				ID:                uuid.New(),
				Type:              "buy",
				Quantity:          "0",
				PriceUSD:          "10.00",
				RemainingQuantity: "0",
				Timestamp:         time.Now().Add(-time.Hour * 24),
			},
		}

		result, err := calculator.CalculatePnL(lots, "20.00")
		require.NoError(t, err)
		
		assert.Equal(t, "0", result.RealizedPnLUSD)
		assert.Equal(t, "0", result.CurrentQuantity)
	})

	t.Run("Invalid price format", func(t *testing.T) {
		lots := []models.PnLLot{
			{
				ID:                uuid.New(),
				Type:              "buy",
				Quantity:          "100",
				PriceUSD:          "invalid",
				RemainingQuantity: "100",
				Timestamp:         time.Now(),
			},
		}

		_, err := calculator.CalculatePnL(lots, "20.00")
		assert.Error(t, err)
	})

	t.Run("Invalid current price", func(t *testing.T) {
		lots := []models.PnLLot{
			{
				ID:                uuid.New(),
				Type:              "buy",
				Quantity:          "100",
				PriceUSD:          "10.00",
				RemainingQuantity: "100",
				Timestamp:         time.Now(),
			},
		}

		_, err := calculator.CalculatePnL(lots, "invalid")
		assert.Error(t, err)
	})
}

func TestCalculator_SortLots(t *testing.T) {
	now := time.Now()
	
	lots := []models.PnLLot{
		{
			ID:        uuid.New(),
			Timestamp: now.Add(-time.Hour * 12), // Middle
		},
		{
			ID:        uuid.New(),
			Timestamp: now.Add(-time.Hour * 24), // Oldest
		},
		{
			ID:        uuid.New(),
			Timestamp: now, // Newest
		},
	}

	t.Run("FIFO sorting", func(t *testing.T) {
		calculator := NewCalculator(FIFO)
		lotsCopy := make([]models.PnLLot, len(lots))
		copy(lotsCopy, lots)
		
		calculator.sortLots(lotsCopy)
		
		// Should be sorted oldest first
		assert.True(t, lotsCopy[0].Timestamp.Before(lotsCopy[1].Timestamp))
		assert.True(t, lotsCopy[1].Timestamp.Before(lotsCopy[2].Timestamp))
	})

	t.Run("LIFO sorting", func(t *testing.T) {
		calculator := NewCalculator(LIFO)
		lotsCopy := make([]models.PnLLot, len(lots))
		copy(lotsCopy, lots)
		
		calculator.sortLots(lotsCopy)
		
		// Should be sorted newest first
		assert.True(t, lotsCopy[0].Timestamp.After(lotsCopy[1].Timestamp))
		assert.True(t, lotsCopy[1].Timestamp.After(lotsCopy[2].Timestamp))
	})
}