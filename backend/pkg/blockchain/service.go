package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/defi-dashboard/backend/internal/models"
	"github.com/defi-dashboard/backend/pkg/external"
	"github.com/defi-dashboard/backend/pkg/logger"
	"github.com/google/uuid"
)

type BlockchainService struct {
	alchemyClient   *AlchemyClient
	coinGeckoClient *external.CoinGeckoClient
}

func NewBlockchainService(alchemyAPIKey, coinGeckoAPIKey string) *BlockchainService {
	return &BlockchainService{
		alchemyClient:   NewAlchemyClient(alchemyAPIKey),
		coinGeckoClient: external.NewCoinGeckoClient(coinGeckoAPIKey),
	}
}

// NewBlockchainServiceWithDynamicKeys creates a blockchain service with runtime API keys
func NewBlockchainServiceWithDynamicKeys(alchemyAPIKey, coinGeckoAPIKey string) *BlockchainService {
	// Use fallback keys from environment if headers are empty
	if alchemyAPIKey == "" {
		alchemyAPIKey = os.Getenv("ALCHEMY_API_KEY")
	}
	if coinGeckoAPIKey == "" {
		coinGeckoAPIKey = os.Getenv("COINGECKO_API_KEY")
	}
	
	logger.Debug("Creating blockchain service with dynamic keys", 
		"hasAlchemy", alchemyAPIKey != "",
		"hasCoinGecko", coinGeckoAPIKey != "")
	
	return NewBlockchainService(alchemyAPIKey, coinGeckoAPIKey)
}

// GetWalletBalances fetches complete wallet balances with USD values
func (s *BlockchainService) GetWalletBalances(ctx context.Context, address string, chainID int) ([]*models.Balance, float64, error) {
	logger.Info("Fetching wallet balances", "address", address, "chainID", chainID)

	// Get token balances
	balances, err := s.alchemyClient.GetTokenBalances(ctx, address, chainID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get token balances: %w", err)
	}

	// Get ETH balance
	ethBalance, err := s.alchemyClient.GetETHBalance(ctx, address, chainID)
	if err != nil {
		logger.Error("Failed to get ETH balance", "error", err)
	} else if ethBalance.Cmp(big.NewInt(0)) > 0 {
		// Add ETH balance to results
		ethToken := s.createETHToken(chainID)
		balance := &models.Balance{
			ID:       uuid.New(),
			WalletID: uuid.New(),
			TokenID:  ethToken.ID,
			Token:    ethToken,
			Balance:  ethBalance.String(),
		}
		balances = append(balances, balance)
	}

	// Get USD prices for all tokens
	totalValue, err := s.enrichBalancesWithPrices(ctx, balances)
	if err != nil {
		logger.Error("Failed to enrich balances with prices", "error", err)
		// Continue without prices
	}

	logger.Info("Successfully fetched wallet balances", 
		"address", address, 
		"tokenCount", len(balances), 
		"totalValue", totalValue)

	return balances, totalValue, nil
}

// enrichBalancesWithPrices adds USD price data to balances
func (s *BlockchainService) enrichBalancesWithPrices(ctx context.Context, balances []*models.Balance) (float64, error) {
	if len(balances) == 0 {
		return 0, nil
	}

	// Map token symbols to CoinGecko IDs
	tokenIDs := make([]string, 0, len(balances))
	symbolToToken := make(map[string]*models.Token)
	
	for _, balance := range balances {
		if balance.Token != nil {
			symbol := strings.ToLower(balance.Token.Symbol)
			if coingeckoID, exists := external.TokenIDMappings[symbol]; exists {
				tokenIDs = append(tokenIDs, coingeckoID)
				symbolToToken[coingeckoID] = balance.Token
			}
		}
	}

	if len(tokenIDs) == 0 {
		return 0, nil
	}

	// Get prices from CoinGecko
	prices, err := s.coinGeckoClient.GetTokenPrices(ctx, tokenIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to get token prices: %w", err)
	}

	// Calculate USD values
	totalValue := 0.0
	for _, balance := range balances {
		if balance.Token == nil {
			continue
		}

		symbol := strings.ToLower(balance.Token.Symbol)
		coingeckoID, exists := external.TokenIDMappings[symbol]
		if !exists {
			continue
		}

		priceData, exists := prices[coingeckoID]
		if !exists {
			continue
		}

		// Update token with price data
		balance.Token.PriceUSD = &priceData.USD
		change24h := priceData.USD24hChange
		balance.Token.PriceChange24h = &change24h

		// Calculate balance USD value
		if balance.Token.Decimals > 0 {
			// Convert from wei to decimal
			balanceInt, ok := new(big.Int).SetString(balance.Balance, 10)
			if !ok {
				logger.Error("Failed to parse balance", "balance", balance.Balance)
				continue
			}

			divisor := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(balance.Token.Decimals)), nil))
			balanceFloat := new(big.Float).SetInt(balanceInt)
			balanceFloat.Quo(balanceFloat, divisor)

			decimalBalance, _ := balanceFloat.Float64()
			usdValue := decimalBalance * priceData.USD
			balance.BalanceUSD = &usdValue
			totalValue += usdValue
		}
	}

	return totalValue, nil
}

// createETHToken creates an ETH token model for the given chain
func (s *BlockchainService) createETHToken(chainID int) *models.Token {
	var symbol, name string
	
	switch chainID {
	case 1: // Ethereum
		symbol = "ETH"
		name = "Ether"
	case 137: // Polygon
		symbol = "MATIC"
		name = "Polygon"
	case 42161: // Arbitrum
		symbol = "ETH"
		name = "Ether"
	case 10: // Optimism
		symbol = "ETH" 
		name = "Ether"
	default:
		symbol = "ETH"
		name = "Ether"
	}

	return &models.Token{
		ID:       uuid.New(),
		Address:  "0x0000000000000000000000000000000000000000", // Native token
		ChainID:  chainID,
		Symbol:   symbol,
		Name:     name,
		Decimals: 18,
	}
}

// GetTransactionHistory fetches transaction history for an address
func (s *BlockchainService) GetTransactionHistory(ctx context.Context, address string, chainID int, limit int) ([]*models.Transaction, error) {
	logger.Info("Fetching transaction history", "address", address, "chainID", chainID, "limit", limit)

	transactions, err := s.alchemyClient.GetTransactions(ctx, address, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	// Limit results
	if limit > 0 && len(transactions) > limit {
		transactions = transactions[:limit]
	}

	logger.Info("Successfully fetched transaction history", 
		"address", address, 
		"transactionCount", len(transactions))

	return transactions, nil
}

// ParseTokenAmount converts a token amount string to a float64 based on decimals
func ParseTokenAmount(amount string, decimals int) (float64, error) {
	if amount == "" || amount == "0" {
		return 0, nil
	}

	// Parse as big integer
	amountInt, ok := new(big.Int).SetString(amount, 10)
	if !ok {
		return 0, fmt.Errorf("invalid amount: %s", amount)
	}

	// Convert to decimal
	divisor := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil))
	amountFloat := new(big.Float).SetInt(amountInt)
	amountFloat.Quo(amountFloat, divisor)

	result, _ := amountFloat.Float64()
	return result, nil
}

// FormatTokenAmount formats a token amount for display
func FormatTokenAmount(amount float64, decimals int) string {
	if amount == 0 {
		return "0"
	}

	// Determine appropriate precision
	var precision int
	if amount >= 1000 {
		precision = 2
	} else if amount >= 1 {
		precision = 4
	} else {
		precision = 6
	}

	return strconv.FormatFloat(amount, 'f', precision, 64)
}

// Chain ID constants
const (
	ChainIDEthereum = 1
	ChainIDPolygon  = 137
	ChainIDArbitrum = 42161
	ChainIDOptimism = 10
)

// GetChainName returns the chain name for a given chain ID
func GetChainName(chainID int) string {
	switch chainID {
	case ChainIDEthereum:
		return "Ethereum"
	case ChainIDPolygon:
		return "Polygon"
	case ChainIDArbitrum:
		return "Arbitrum"
	case ChainIDOptimism:
		return "Optimism"
	default:
		return fmt.Sprintf("Chain %d", chainID)
	}
}

// GetSupportedChains returns list of supported chain IDs
func GetSupportedChains() []int {
	return []int{ChainIDEthereum, ChainIDPolygon, ChainIDArbitrum, ChainIDOptimism}
}