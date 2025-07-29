package pnl

import (
	"errors"
	"math/big"
	"sort"
	"time"

	"github.com/defi-dashboard/backend/internal/models"
)

type CalculationMethod string

const (
	FIFO CalculationMethod = "fifo"
	LIFO CalculationMethod = "lifo"
)

type Calculator struct {
	method CalculationMethod
}

func NewCalculator(method CalculationMethod) *Calculator {
	return &Calculator{method: method}
}

// CalculatePnL calculates realized and unrealized PnL for a given set of lots
func (c *Calculator) CalculatePnL(lots []models.PnLLot, currentPriceUSD string) (*models.PnLCalculation, error) {
	if len(lots) == 0 {
		return nil, errors.New("no lots provided for calculation")
	}

	// Separate buy and sell lots
	var buys, sells []models.PnLLot
	for _, lot := range lots {
		if lot.Type == "buy" {
			buys = append(buys, lot)
		} else if lot.Type == "sell" {
			sells = append(sells, lot)
		}
	}

	// Sort based on method
	c.sortLots(buys)
	c.sortLots(sells)

	// Calculate realized PnL from matched lots
	realizedPnL, processedBuys, err := c.calculateRealizedPnL(buys, sells)
	if err != nil {
		return nil, err
	}

	// Calculate unrealized PnL from remaining buy lots
	unrealizedPnL, totalCostBasis, currentQuantity, err := c.calculateUnrealizedPnL(processedBuys, currentPriceUSD)
	if err != nil {
		return nil, err
	}

	// Calculate current value
	currentValue, err := c.multiplyDecimals(currentQuantity, currentPriceUSD)
	if err != nil {
		return nil, err
	}

	// Calculate total PnL
	totalPnL, err := c.addDecimals(realizedPnL, unrealizedPnL)
	if err != nil {
		return nil, err
	}

	// Get token info from first lot
	var walletAddress, tokenAddress, tokenSymbol string
	if len(lots) > 0 {
		// These would need to be populated from the calling context
		// since lots don't contain wallet/token addresses directly
	}

	return &models.PnLCalculation{
		WalletAddress:     walletAddress,
		TokenAddress:      tokenAddress,
		TokenSymbol:       tokenSymbol,
		Method:            string(c.method),
		RealizedPnLUSD:    realizedPnL,
		UnrealizedPnLUSD:  unrealizedPnL,
		TotalPnLUSD:       totalPnL,
		TotalCostBasisUSD: totalCostBasis,
		CurrentValueUSD:   currentValue,
		CurrentQuantity:   currentQuantity,
		Lots:              processedBuys,
		CalculatedAt:      time.Now(),
	}, nil
}

// sortLots sorts lots based on the calculation method
func (c *Calculator) sortLots(lots []models.PnLLot) {
	if c.method == FIFO {
		// Sort by timestamp ascending (earliest first)
		sort.Slice(lots, func(i, j int) bool {
			return lots[i].Timestamp.Before(lots[j].Timestamp)
		})
	} else { // LIFO
		// Sort by timestamp descending (latest first)
		sort.Slice(lots, func(i, j int) bool {
			return lots[i].Timestamp.After(lots[j].Timestamp)
		})
	}
}

// calculateRealizedPnL matches sell lots against buy lots to calculate realized PnL
func (c *Calculator) calculateRealizedPnL(buys, sells []models.PnLLot) (string, []models.PnLLot, error) {
	// Create copies to avoid modifying originals
	buysCopy := make([]models.PnLLot, len(buys))
	copy(buysCopy, buys)

	totalRealizedPnL := "0"

	for _, sell := range sells {
		sellQuantity, err := c.stringToBigFloat(sell.Quantity)
		if err != nil {
			return "", nil, err
		}

		remainingSellQuantity := new(big.Float).Set(sellQuantity)
		sellPrice, err := c.stringToBigFloat(sell.PriceUSD)
		if err != nil {
			return "", nil, err
		}

		// Match against buy lots
		for i := range buysCopy {
			if remainingSellQuantity.Cmp(big.NewFloat(0)) <= 0 {
				break
			}

			buyRemaining, err := c.stringToBigFloat(buysCopy[i].RemainingQuantity)
			if err != nil {
				return "", nil, err
			}

			if buyRemaining.Cmp(big.NewFloat(0)) <= 0 {
				continue
			}

			buyPrice, err := c.stringToBigFloat(buysCopy[i].PriceUSD)
			if err != nil {
				return "", nil, err
			}

			// Calculate quantity to match
			matchedQuantity := new(big.Float)
			if remainingSellQuantity.Cmp(buyRemaining) <= 0 {
				matchedQuantity.Set(remainingSellQuantity)
				remainingSellQuantity.SetFloat64(0)
			} else {
				matchedQuantity.Set(buyRemaining)
				remainingSellQuantity.Sub(remainingSellQuantity, buyRemaining)
			}

			// Update remaining quantity in buy lot
			newRemaining := new(big.Float).Sub(buyRemaining, matchedQuantity)
			buysCopy[i].RemainingQuantity = newRemaining.String()

			// Calculate realized PnL for this match
			costBasis := new(big.Float).Mul(matchedQuantity, buyPrice)
			proceeds := new(big.Float).Mul(matchedQuantity, sellPrice)
			pnl := new(big.Float).Sub(proceeds, costBasis)

			// Add to total realized PnL
			currentTotal, err := c.stringToBigFloat(totalRealizedPnL)
			if err != nil {
				return "", nil, err
			}
			currentTotal.Add(currentTotal, pnl)
			totalRealizedPnL = currentTotal.String()
		}
	}

	return totalRealizedPnL, buysCopy, nil
}

// calculateUnrealizedPnL calculates unrealized PnL from remaining buy lots
func (c *Calculator) calculateUnrealizedPnL(buys []models.PnLLot, currentPriceUSD string) (string, string, string, error) {
	currentPrice, err := c.stringToBigFloat(currentPriceUSD)
	if err != nil {
		return "", "", "", err
	}

	totalUnrealizedPnL := big.NewFloat(0)
	totalCostBasis := big.NewFloat(0)
	totalCurrentQuantity := big.NewFloat(0)

	for _, buy := range buys {
		remaining, err := c.stringToBigFloat(buy.RemainingQuantity)
		if err != nil {
			return "", "", "", err
		}

		if remaining.Cmp(big.NewFloat(0)) <= 0 {
			continue
		}

		buyPrice, err := c.stringToBigFloat(buy.PriceUSD)
		if err != nil {
			return "", "", "", err
		}

		// Calculate cost basis for remaining quantity
		costBasis := new(big.Float).Mul(remaining, buyPrice)
		totalCostBasis.Add(totalCostBasis, costBasis)

		// Calculate current value for remaining quantity
		currentValue := new(big.Float).Mul(remaining, currentPrice)

		// Calculate unrealized PnL
		pnl := new(big.Float).Sub(currentValue, costBasis)
		totalUnrealizedPnL.Add(totalUnrealizedPnL, pnl)

		// Add to total quantity
		totalCurrentQuantity.Add(totalCurrentQuantity, remaining)
	}

	return totalUnrealizedPnL.String(), totalCostBasis.String(), totalCurrentQuantity.String(), nil
}

// Helper functions for decimal arithmetic
func (c *Calculator) stringToBigFloat(s string) (*big.Float, error) {
	f, _, err := big.ParseFloat(s, 10, 256, big.ToNearestEven)
	return f, err
}

func (c *Calculator) addDecimals(a, b string) (string, error) {
	aFloat, err := c.stringToBigFloat(a)
	if err != nil {
		return "", err
	}
	bFloat, err := c.stringToBigFloat(b)
	if err != nil {
		return "", err
	}
	result := new(big.Float).Add(aFloat, bFloat)
	return result.String(), nil
}

func (c *Calculator) multiplyDecimals(a, b string) (string, error) {
	aFloat, err := c.stringToBigFloat(a)
	if err != nil {
		return "", err
	}
	bFloat, err := c.stringToBigFloat(b)
	if err != nil {
		return "", err
	}
	result := new(big.Float).Mul(aFloat, bFloat)
	return result.String(), nil
}