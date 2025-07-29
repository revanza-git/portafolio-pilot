package repos

import (
	"context"
	"time"

	"github.com/defi-dashboard/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type transactionRepository struct {
	db *pgxpool.Pool
}

// NewTransactionRepository creates a new transaction repository
func NewTransactionRepository(db *pgxpool.Pool) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) GetByHash(ctx context.Context, hash string) (*models.Transaction, error) {
	// TODO: Implement actual database query
	toAddr := "0x0987654321098765432109876543210987654321"
	value := "1000000000000000000"
	gasUsed := int64(21000)
	gasPrice := "20000000000"
	gasFeeUSD := float64(2.50)
	blockNumber := int64(18500000)
	
	return &models.Transaction{
		ID:          uuid.New(),
		Hash:        hash,
		ChainID:     1,
		FromAddress: "0x1234567890123456789012345678901234567890",
		ToAddress:   &toAddr,
		Value:       &value,
		GasUsed:     &gasUsed,
		GasPrice:    &gasPrice,
		GasFeeUSD:   &gasFeeUSD,
		BlockNumber: &blockNumber,
		Timestamp:   time.Now().Add(-1 * time.Hour),
		Status:      "confirmed",
		Type:        "send",
	}, nil
}

func (r *transactionRepository) GetUserTransactions(ctx context.Context, userID uuid.UUID, filters TransactionFilters) ([]*models.Transaction, error) {
	// TODO: Implement actual database query
	// Return mock transactions
	return r.getMockTransactions(), nil
}

func (r *transactionRepository) GetWalletTransactions(ctx context.Context, address string, chainID int, txType *string, limit, offset int) ([]*models.Transaction, error) {
	// TODO: Implement actual database query
	return r.getMockTransactions(), nil
}

func (r *transactionRepository) Create(ctx context.Context, tx *models.Transaction) (*models.Transaction, error) {
	// TODO: Implement actual database insert
	tx.ID = uuid.New()
	return tx, nil
}

func (r *transactionRepository) UpdateStatus(ctx context.Context, hash, status string, blockNumber, gasUsed int64, gasFeeUSD float64) (*models.Transaction, error) {
	// TODO: Implement actual database update
	tx, _ := r.GetByHash(ctx, hash)
	tx.Status = status
	tx.BlockNumber = &blockNumber
	tx.GasUsed = &gasUsed
	tx.GasFeeUSD = &gasFeeUSD
	return tx, nil
}

func (r *transactionRepository) LinkToUser(ctx context.Context, userID, transactionID, walletID uuid.UUID) error {
	// TODO: Implement actual database insert
	return nil
}

func (r *transactionRepository) getMockTransactions() []*models.Transaction {
	// Mock transaction data
	toAddr1 := "0x0987654321098765432109876543210987654321"
	value1 := "1000000000000000000"
	gasUsed1 := int64(21000)
	gasPrice1 := "20000000000"
	gasFeeUSD1 := float64(2.50)
	blockNumber1 := int64(18500000)

	toAddr2 := "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"
	value2 := "50000000" // 50 USDC
	gasUsed2 := int64(65000)
	gasPrice2 := "25000000000"
	gasFeeUSD2 := float64(8.12)
	blockNumber2 := int64(18499900)

	transactions := []*models.Transaction{
		{
			ID:          uuid.New(),
			Hash:        "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			ChainID:     1,
			FromAddress: "0x1234567890123456789012345678901234567890",
			ToAddress:   &toAddr1,
			Value:       &value1,
			GasUsed:     &gasUsed1,
			GasPrice:    &gasPrice1,
			GasFeeUSD:   &gasFeeUSD1,
			BlockNumber: &blockNumber1,
			Timestamp:   time.Now().Add(-1 * time.Hour),
			Status:      "confirmed",
			Type:        "send",
		},
		{
			ID:          uuid.New(),
			Hash:        "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
			ChainID:     1,
			FromAddress: "0x0987654321098765432109876543210987654321",
			ToAddress:   &toAddr2,
			Value:       &value2,
			GasUsed:     &gasUsed2,
			GasPrice:    &gasPrice2,
			GasFeeUSD:   &gasFeeUSD2,
			BlockNumber: &blockNumber2,
			Timestamp:   time.Now().Add(-2 * time.Hour),
			Status:      "confirmed",
			Type:        "approve",
			Metadata: map[string]interface{}{
				"token_address": "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
				"spender": "0x68b3465833fb72A70ecDF485E0e4C7bD8665Fc45",
				"amount": "115792089237316195423570985008687907853269984665640564039457584007913129639935",
			},
		},
		{
			ID:          uuid.New(),
			Hash:        "0xfedcba0987654321fedcba0987654321fedcba0987654321fedcba0987654321",
			ChainID:     1,
			FromAddress: "0x1234567890123456789012345678901234567890",
			ToAddress:   nil,
			Value:       nil,
			GasUsed:     &gasUsed2,
			GasPrice:    &gasPrice2,
			GasFeeUSD:   &gasFeeUSD2,
			BlockNumber: &blockNumber2,
			Timestamp:   time.Now().Add(-3 * time.Hour),
			Status:      "confirmed",
			Type:        "swap",
			Metadata: map[string]interface{}{
				"token_in": "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
				"token_out": "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48", 
				"amount_in": "1000000000000000000",
				"amount_out": "2345670000",
				"router": "0x68b3465833fb72A70ecDF485E0e4C7bD8665Fc45",
			},
		},
	}

	return transactions
}