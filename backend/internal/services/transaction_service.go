package services

import (
	"context"
	"fmt"
	"time"

	"github.com/defi-dashboard/backend/internal/models"
	"github.com/defi-dashboard/backend/internal/repos"
	"github.com/defi-dashboard/backend/pkg/blockchain"
	"github.com/defi-dashboard/backend/pkg/errors"
	"github.com/defi-dashboard/backend/pkg/logger"
	"github.com/google/uuid"
)

type TransactionService struct {
	transactionRepo repos.TransactionRepository
}

func NewTransactionService(transactionRepo repos.TransactionRepository) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
	}
}

// GetTransactions returns real transactions for an address from blockchain
func (s *TransactionService) GetTransactions(ctx context.Context, address string, chainID *int, txType *string, page, limit int, alchemyAPIKey, coinGeckoAPIKey string) (*TransactionResponse, error) {
	logger.Info("Fetching transactions", "address", address, "chainID", chainID, "type", txType)

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
		return nil, errors.BadRequest(fmt.Sprintf("unsupported chain ID: %d", chain))
	}

	// Create blockchain service with dynamic API keys
	blockchainService := blockchain.NewBlockchainServiceWithDynamicKeys(alchemyAPIKey, coinGeckoAPIKey)
	
	// Get transactions from blockchain
	transactions, err := blockchainService.GetTransactionHistory(ctx, address, chain, limit*2) // Get more to allow filtering
	if err != nil {
		logger.Error("Failed to fetch transactions from blockchain", "error", err)
		return nil, errors.Internal("Failed to fetch transactions")
	}

	// Filter by transaction type if specified
	if txType != nil {
		filteredTxs := make([]*models.Transaction, 0)
		for _, tx := range transactions {
			if tx.Type == *txType {
				filteredTxs = append(filteredTxs, tx)
			}
		}
		transactions = filteredTxs
	}

	// Apply pagination
	offset := (page - 1) * limit
	total := len(transactions)
	end := offset + limit
	
	if offset >= total {
		transactions = []*models.Transaction{}
	} else {
		if end > total {
			end = total
		}
		transactions = transactions[offset:end]
	}

	// Store transactions in database for caching (optional)
	if err := s.storeTransactions(ctx, address, chain, transactions); err != nil {
		logger.Error("Failed to store transactions in database", "error", err)
		// Continue without storing - this is not critical
	}

	logger.Info("Successfully fetched transactions", 
		"address", address, 
		"chainID", chain, 
		"count", len(transactions),
		"total", total)

	return &TransactionResponse{
		Data: transactions,
		Meta: &PaginationMeta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: (total + limit - 1) / limit,
		},
	}, nil
}

// GetApprovals returns token approvals for an address (placeholder - requires specialized API)
func (s *TransactionService) GetApprovals(ctx context.Context, address string, chainID *int, activeOnly bool) ([]*TokenApproval, error) {
	logger.Info("Fetching token approvals", "address", address, "chainID", chainID)

	// TODO: Implement real approval fetching
	// This requires querying approval events from blockchain or using specialized APIs like Etherscan
	// For now, we'll return empty results with a note that this needs implementation
	
	logger.Warn("Token approval fetching not yet implemented - returning empty results")
	
	return []*TokenApproval{}, nil
}

// RevokeApproval revokes a token approval (placeholder - requires transaction signing)
func (s *TransactionService) RevokeApproval(ctx context.Context, address, token, spender string) (string, error) {
	logger.Info("Revoking token approval", "address", address, "token", token, "spender", spender)

	// TODO: Implement actual blockchain transaction to revoke approval
	// This would require:
	// 1. Create a transaction to call approve(spender, 0) on the token contract
	// 2. Sign the transaction (requires private key or wallet integration)
	// 3. Broadcast the transaction to the network
	// 4. Return the transaction hash

	return "", errors.BadRequest("Approval revocation not yet implemented - requires wallet integration")
}

// storeTransactions stores transaction data in database for caching
func (s *TransactionService) storeTransactions(ctx context.Context, address string, chainID int, transactions []*models.Transaction) error {
	// This is optional - for caching and historical tracking
	// Implementation would involve storing transactions in the database
	
	logger.Debug("Would store transactions in database", 
		"address", address, 
		"chainID", chainID, 
		"transactionCount", len(transactions))
	
	return nil
}

// GetTransactionsByType returns transactions filtered by type
func (s *TransactionService) GetTransactionsByType(ctx context.Context, address string, chainID int, txType string, limit int, alchemyAPIKey, coinGeckoAPIKey string) ([]*models.Transaction, error) {
	logger.Info("Fetching transactions by type", "address", address, "chainID", chainID, "type", txType)

	// Create blockchain service with dynamic API keys
	blockchainService := blockchain.NewBlockchainServiceWithDynamicKeys(alchemyAPIKey, coinGeckoAPIKey)

	transactions, err := blockchainService.GetTransactionHistory(ctx, address, chainID, limit*3) // Get more to filter
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %w", err)
	}

	// Filter by type
	var filtered []*models.Transaction
	for _, tx := range transactions {
		if tx.Type == txType {
			filtered = append(filtered, tx)
			if len(filtered) >= limit {
				break
			}
		}
	}

	return filtered, nil
}

// GetTransactionsSummary returns a summary of transaction activity
func (s *TransactionService) GetTransactionsSummary(ctx context.Context, address string, chainID int, days int, alchemyAPIKey, coinGeckoAPIKey string) (*TransactionSummary, error) {
	logger.Info("Fetching transaction summary", "address", address, "chainID", chainID, "days", days)

	// Create blockchain service with dynamic API keys
	blockchainService := blockchain.NewBlockchainServiceWithDynamicKeys(alchemyAPIKey, coinGeckoAPIKey)

	transactions, err := blockchainService.GetTransactionHistory(ctx, address, chainID, 100) // Get recent transactions
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %w", err)
	}

	// Calculate summary
	cutoff := time.Now().AddDate(0, 0, -days)
	summary := &TransactionSummary{
		Period:            days,
		TotalTransactions: 0,
		TypeCounts:        make(map[string]int),
		TotalGasFeesUSD:   0.0,
	}

	for _, tx := range transactions {
		if tx.Timestamp.After(cutoff) {
			summary.TotalTransactions++
			summary.TypeCounts[tx.Type]++
			
			if tx.GasFeeUSD != nil {
				summary.TotalGasFeesUSD += *tx.GasFeeUSD
			}
		}
	}

	return summary, nil
}

// Helper types and functions

type TransactionResponse struct {
	Data []*models.Transaction `json:"data"`
	Meta *PaginationMeta       `json:"meta"`
}

type PaginationMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type TokenApproval struct {
	ID              uuid.UUID     `json:"id"`
	Token           *models.Token `json:"token"`
	SpenderAddress  string        `json:"spender_address"`
	SpenderName     *string       `json:"spender_name,omitempty"`
	Allowance       string        `json:"allowance"`
	AllowanceUSD    *float64      `json:"allowance_usd,omitempty"`
	LastUpdated     time.Time     `json:"last_updated"`
}

type TransactionSummary struct {
	Period            int            `json:"period_days"`
	TotalTransactions int            `json:"total_transactions"`
	TypeCounts        map[string]int `json:"type_counts"`
	TotalGasFeesUSD   float64        `json:"total_gas_fees_usd"`
}