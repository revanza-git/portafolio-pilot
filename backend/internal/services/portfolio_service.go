package services

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/defi-dashboard/backend/internal/models"
	"github.com/defi-dashboard/backend/internal/repos"
	"github.com/defi-dashboard/backend/pkg/utils"
	"github.com/google/uuid"
)

type PortfolioService struct {
	walletRepo repos.WalletRepository
	tokenRepo  repos.TokenRepository
}

func NewPortfolioService(walletRepo repos.WalletRepository, tokenRepo repos.TokenRepository) *PortfolioService {
	return &PortfolioService{
		walletRepo: walletRepo,
		tokenRepo:  tokenRepo,
	}
}

// GetBalances returns token balances for an address
func (s *PortfolioService) GetBalances(ctx context.Context, address string, chainID *int, hideSmall bool) (*PortfolioBalances, error) {
	// TODO: Fetch actual balances from blockchain
	// For now, return mock data

	// Get mock tokens
	tokens := s.getMockTokens(chainID)
	
	balances := make([]*models.Balance, 0)
	totalValue := 0.0

	for _, token := range tokens {
		// Generate random balance
		balance := s.generateMockBalance(token)
		
		// Skip small balances if requested
		if hideSmall && balance.BalanceUSD != nil && *balance.BalanceUSD < 1.0 {
			continue
		}

		balances = append(balances, balance)
		if balance.BalanceUSD != nil {
			totalValue += *balance.BalanceUSD
		}
	}

	return &PortfolioBalances{
		TotalValue: totalValue,
		Balances:   balances,
	}, nil
}

// GetHistory returns portfolio value history
func (s *PortfolioService) GetHistory(ctx context.Context, address string, chainID *int, period string, interval string) ([]*PortfolioHistoryPoint, error) {
	// TODO: Fetch actual historical data
	// For now, generate mock history

	// Determine time range based on period
	endTime := time.Now()
	var startTime time.Time
	var points int

	switch period {
	case "1d":
		startTime = endTime.Add(-24 * time.Hour)
		points = 24
	case "1w":
		startTime = endTime.Add(-7 * 24 * time.Hour)
		points = 7
	case "1m":
		startTime = endTime.Add(-30 * 24 * time.Hour)
		points = 30
	case "3m":
		startTime = endTime.Add(-90 * 24 * time.Hour)
		points = 90
	case "1y":
		startTime = endTime.Add(-365 * 24 * time.Hour)
		points = 365
	default:
		startTime = endTime.Add(-7 * 24 * time.Hour)
		points = 7
	}

	// Generate mock history points
	history := make([]*PortfolioHistoryPoint, points)
	baseValue := 10000.0
	
	for i := 0; i < points; i++ {
		timestamp := startTime.Add(time.Duration(i) * (endTime.Sub(startTime) / time.Duration(points)))
		
		// Add some random variation
		variation := (rand.Float64() - 0.5) * 0.1 // ±5% variation
		value := baseValue * (1 + variation + float64(i)/float64(points)*0.2) // Overall 20% growth
		
		history[i] = &PortfolioHistoryPoint{
			Timestamp:  timestamp,
			TotalValue: value,
		}
	}

	return history, nil
}

// Helper functions

func (s *PortfolioService) getMockTokens(chainID *int) []*models.Token {
	chain := 1
	if chainID != nil {
		chain = *chainID
	}

	// Mock tokens with prices
	return []*models.Token{
		{
			ID:       uuid.New(),
			Address:  "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
			ChainID:  chain,
			Symbol:   "WETH",
			Name:     "Wrapped Ether",
			Decimals: 18,
			LogoURI:  utils.StrPtr("https://assets.coingecko.com/coins/images/2518/small/weth.png"),
			PriceUSD: utils.Float64Ptr(2345.67),
		},
		{
			ID:       uuid.New(),
			Address:  "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
			ChainID:  chain,
			Symbol:   "USDC",
			Name:     "USD Coin",
			Decimals: 6,
			LogoURI:  utils.StrPtr("https://assets.coingecko.com/coins/images/6319/small/USD_Coin_icon.png"),
			PriceUSD: utils.Float64Ptr(1.0),
		},
		{
			ID:       uuid.New(),
			Address:  "0x6B175474E89094C44Da98b954EedeAC495271d0F",
			ChainID:  chain,
			Symbol:   "DAI",
			Name:     "Dai Stablecoin",
			Decimals: 18,
			LogoURI:  utils.StrPtr("https://assets.coingecko.com/coins/images/9956/small/4943.png"),
			PriceUSD: utils.Float64Ptr(1.0),
		},
		{
			ID:       uuid.New(),
			Address:  "0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599",
			ChainID:  chain,
			Symbol:   "WBTC",
			Name:     "Wrapped Bitcoin",
			Decimals: 8,
			LogoURI:  utils.StrPtr("https://assets.coingecko.com/coins/images/7598/small/wrapped_bitcoin_wbtc.png"),
			PriceUSD: utils.Float64Ptr(45678.90),
		},
	}
}

func (s *PortfolioService) generateMockBalance(token *models.Token) *models.Balance {
	// Generate random balance amount
	amounts := map[string]float64{
		"WETH": 2.5 + rand.Float64()*2,      // 2.5-4.5 ETH
		"USDC": 5000 + rand.Float64()*10000, // 5k-15k USDC
		"DAI":  2000 + rand.Float64()*5000,  // 2k-7k DAI
		"WBTC": 0.1 + rand.Float64()*0.2,    // 0.1-0.3 BTC
	}

	amount := amounts[token.Symbol]
	if amount == 0 {
		amount = 100 + rand.Float64()*1000
	}

	// Calculate USD value
	var balanceUSD float64
	if token.PriceUSD != nil {
		balanceUSD = amount * (*token.PriceUSD)
	}

	return &models.Balance{
		ID:         uuid.New(),
		WalletID:   uuid.New(),
		TokenID:    token.ID,
		Token:      token,
		Balance:    formatBalance(amount, token.Decimals),
		BalanceUSD: &balanceUSD,
	}
}

func formatBalance(amount float64, decimals int) string {
	// Convert to smallest unit (like wei)
	// For simplicity, just return as string
	return fmt.Sprintf("%.0f", amount * math.Pow(10, float64(decimals)))
}

// Response types

type PortfolioBalances struct {
	TotalValue float64           `json:"total_value"`
	Balances   []*models.Balance `json:"balances"`
}

type PortfolioHistoryPoint struct {
	Timestamp  time.Time `json:"timestamp"`
	TotalValue float64   `json:"total_value"`
}

type Allocation struct {
	Token      *models.Token `json:"token"`
	Percentage float64       `json:"percentage"`
}

