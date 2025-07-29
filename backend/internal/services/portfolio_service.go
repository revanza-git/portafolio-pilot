package services

import (
	"context"
	"fmt"
	"time"

	"github.com/defi-dashboard/backend/internal/models"
	"github.com/defi-dashboard/backend/internal/repos"
	"github.com/defi-dashboard/backend/pkg/blockchain"
	"github.com/defi-dashboard/backend/pkg/logger"
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

// GetBalances returns real token balances for an address from blockchain
func (s *PortfolioService) GetBalances(ctx context.Context, address string, chainID *int, hideSmall bool, alchemyAPIKey, coinGeckoAPIKey string) (*PortfolioBalances, error) {
	logger.Info("Fetching portfolio balances", "address", address, "chainID", chainID)

	// Default to Ethereum mainnet if no chain specified
	chain := 1
	if chainID != nil {
		chain = *chainID
	}

	// Validate chain is supported
	supportedChains := blockchain.GetSupportedChains()
	isSupported := false
	for _, supportedChain := range supportedChains {
		if chain == supportedChain {
			isSupported = true
			break
		}
	}

	if !isSupported {
		return nil, fmt.Errorf("unsupported chain ID: %d. Supported chains: %v", chain, supportedChains)
	}

	// Create blockchain service with dynamic API keys
	blockchainService := blockchain.NewBlockchainServiceWithDynamicKeys(alchemyAPIKey, coinGeckoAPIKey)
	
	// Get real balances from blockchain
	balances, totalValue, err := blockchainService.GetWalletBalances(ctx, address, chain)
	if err != nil {
		logger.Error("Failed to fetch wallet balances", "error", err, "address", address, "chainID", chain)
		return nil, fmt.Errorf("failed to fetch wallet balances: %w", err)
	}

	// Filter small balances if requested
	if hideSmall {
		filteredBalances := make([]*models.Balance, 0)
		filteredTotalValue := 0.0

		for _, balance := range balances {
			if balance.BalanceUSD != nil && *balance.BalanceUSD >= 1.0 {
				filteredBalances = append(filteredBalances, balance)
				filteredTotalValue += *balance.BalanceUSD
			}
		}

		balances = filteredBalances
		totalValue = filteredTotalValue
	}

	// Store/update balances in database (optional - for caching)
	if err := s.storeBalances(ctx, address, chain, balances); err != nil {
		logger.Error("Failed to store balances in database", "error", err)
		// Continue without storing - this is not critical
	}

	logger.Info("Successfully fetched portfolio balances", 
		"address", address, 
		"chainID", chain,
		"tokenCount", len(balances), 
		"totalValue", totalValue)

	return &PortfolioBalances{
		TotalValue: totalValue,
		Balances:   balances,
	}, nil
}

// GetHistory returns portfolio value history (currently mock - real implementation would need historical data)
func (s *PortfolioService) GetHistory(ctx context.Context, address string, chainID *int, period string, interval string, alchemyAPIKey, coinGeckoAPIKey string) ([]*PortfolioHistoryPoint, error) {
	logger.Info("Fetching portfolio history", "address", address, "period", period)

	// TODO: Implement real historical data fetching
	// For now, we generate mock history based on current portfolio value
	
	// Get current portfolio value
	currentBalances, err := s.GetBalances(ctx, address, chainID, false, alchemyAPIKey, coinGeckoAPIKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get current balances for history: %w", err)
	}

	currentValue := currentBalances.TotalValue
	if currentValue == 0 {
		currentValue = 1000 // Default fallback for empty wallets
	}

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

	// Generate realistic history based on current value
	history := make([]*PortfolioHistoryPoint, points+1)
	
	for i := 0; i <= points; i++ {
		timestamp := startTime.Add(time.Duration(i) * (endTime.Sub(startTime) / time.Duration(points)))
		
		// Create a realistic trend that ends at current value
		progress := float64(i) / float64(points)
		
		// Base value starts slightly lower and grows to current
		baseValue := currentValue * (0.85 + 0.15*progress)
		
		// Add some realistic volatility (smaller for shorter periods)
		volatilityFactor := 0.05 // 5% volatility
		if period == "1d" {
			volatilityFactor = 0.02 // 2% for daily
		} else if period == "1y" {
			volatilityFactor = 0.15 // 15% for yearly
		}
		
		// Use deterministic "randomness" based on timestamp for consistency
		seed := timestamp.Unix() % 1000
		variation := (float64(seed%200)/100.0 - 1.0) * volatilityFactor
		
		value := baseValue * (1 + variation)
		
		// Ensure last point is exactly current value
		if i == points {
			value = currentValue
		}
		
		history[i] = &PortfolioHistoryPoint{
			Timestamp:  timestamp,
			TotalValue: value,
		}
	}

	logger.Info("Generated portfolio history", 
		"address", address, 
		"points", len(history),
		"currentValue", currentValue)

	return history, nil
}

// storeBalances stores/updates balance data in database for caching
func (s *PortfolioService) storeBalances(ctx context.Context, address string, chainID int, balances []*models.Balance) error {
	// This is optional - for caching and historical tracking
	// Implementation would involve:
	// 1. Find or create wallet record
	// 2. Update/insert token records
	// 3. Update/insert balance records
	// 4. Clean up old balance records
	
	// For now, just log that we would store this data
	logger.Debug("Would store balances in database", 
		"address", address, 
		"chainID", chainID, 
		"balanceCount", len(balances))
	
	return nil
}

// GetMultiChainBalances gets balances across multiple chains
func (s *PortfolioService) GetMultiChainBalances(ctx context.Context, address string, hideSmall bool, alchemyAPIKey, coinGeckoAPIKey string) (*MultiChainPortfolio, error) {
	logger.Info("Fetching multi-chain portfolio", "address", address)

	supportedChains := blockchain.GetSupportedChains()
	chainBalances := make(map[int]*PortfolioBalances)
	totalValue := 0.0

	for _, chainID := range supportedChains {
		balances, err := s.GetBalances(ctx, address, &chainID, hideSmall, alchemyAPIKey, coinGeckoAPIKey)
		if err != nil {
			logger.Error("Failed to get balances for chain", "chainID", chainID, "error", err)
			// Continue with other chains
			continue
		}

		if balances.TotalValue > 0 {
			chainBalances[chainID] = balances
			totalValue += balances.TotalValue
		}
	}

	return &MultiChainPortfolio{
		TotalValue:    totalValue,
		ChainBalances: chainBalances,
	}, nil
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

type MultiChainPortfolio struct {
	TotalValue    float64                     `json:"total_value"`
	ChainBalances map[int]*PortfolioBalances  `json:"chain_balances"`
}

type Allocation struct {
	Token      *models.Token `json:"token"`
	Percentage float64       `json:"percentage"`
}