package services

import (
	"context"
	"time"

	"github.com/defi-dashboard/backend/internal/models"
	"github.com/defi-dashboard/backend/internal/repos"
	"github.com/defi-dashboard/backend/pkg/errors"
	"github.com/defi-dashboard/backend/pkg/utils"
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

// GetTransactions returns transactions for an address
func (s *TransactionService) GetTransactions(ctx context.Context, address string, chainID *int, txType *string, page, limit int) (*TransactionResponse, error) {
	// Calculate offset
	offset := (page - 1) * limit

	// TODO: Implement actual transaction fetching from blockchain/database
	// For now, return mock data
	transactions, err := s.transactionRepo.GetWalletTransactions(ctx, address, 1, txType, limit, offset)
	if err != nil {
		return nil, errors.Internal("Failed to fetch transactions")
	}

	// Mock total count
	total := 100

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

// GetApprovals returns token approvals for an address
func (s *TransactionService) GetApprovals(ctx context.Context, address string, chainID *int, activeOnly bool) ([]*TokenApproval, error) {
	// TODO: Fetch actual approvals from blockchain
	// For now, return mock data

	approvals := []*TokenApproval{
		{
			ID: uuid.New(),
			Token: &models.Token{
				ID:       uuid.New(),
				Address:  "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
				ChainID:  1,
				Symbol:   "USDC",
				Name:     "USD Coin",
				Decimals: 6,
				LogoURI:  utils.StrPtr("https://assets.coingecko.com/coins/images/6319/small/USD_Coin_icon.png"),
			},
			SpenderAddress: "0x68b3465833fb72A70ecDF485E0e4C7bD8665Fc45",
			SpenderName:    utils.StrPtr("Uniswap V3 Router"),
			Allowance:      "115792089237316195423570985008687907853269984665640564039457584007913129639935", // Max uint256
			AllowanceUSD:   nil, // Unlimited
			LastUpdated:    time.Now().Add(-30 * 24 * time.Hour),
		},
		{
			ID: uuid.New(),
			Token: &models.Token{
				ID:       uuid.New(),
				Address:  "0x6B175474E89094C44Da98b954EedeAC495271d0F",
				ChainID:  1,
				Symbol:   "DAI",
				Name:     "Dai Stablecoin",
				Decimals: 18,
				LogoURI:  utils.StrPtr("https://assets.coingecko.com/coins/images/9956/small/4943.png"),
			},
			SpenderAddress: "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D",
			SpenderName:    utils.StrPtr("Uniswap V2 Router"),
			Allowance:      "1000000000000000000000", // 1000 DAI
			AllowanceUSD:   utils.Float64Ptr(1000.0),
			LastUpdated:    time.Now().Add(-60 * 24 * time.Hour),
		},
	}

	if activeOnly {
		// Filter out zero approvals
		activeApprovals := make([]*TokenApproval, 0)
		for _, approval := range approvals {
			if approval.Allowance != "0" {
				activeApprovals = append(activeApprovals, approval)
			}
		}
		return activeApprovals, nil
	}

	return approvals, nil
}

// RevokeApproval revokes a token approval
func (s *TransactionService) RevokeApproval(ctx context.Context, address, token, spender string) (string, error) {
	// TODO: Implement actual blockchain transaction to revoke approval
	// For now, return mock transaction hash

	// This would typically:
	// 1. Create a transaction to set allowance to 0
	// 2. Sign and broadcast the transaction
	// 3. Return the transaction hash

	mockTxHash := "0x" + generateMockHash()
	return mockTxHash, nil
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

func generateMockHash() string {
	// Generate a mock 64-character hex string
	const hexChars = "0123456789abcdef"
	hash := make([]byte, 64)
	for i := range hash {
		hash[i] = hexChars[time.Now().UnixNano()%16]
	}
	return string(hash)
}

