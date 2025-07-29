package pnl

import (
	"context"
	"fmt"
	"time"

	"github.com/defi-dashboard/backend/internal/models"
	"github.com/defi-dashboard/backend/internal/repos"
	"github.com/google/uuid"
)

type Service interface {
	CalculatePnL(ctx context.Context, walletAddress string, from, to time.Time, method CalculationMethod) (*models.PnLCalculation, error)
	CalculatePnLByToken(ctx context.Context, walletAddress, tokenAddress string, from, to time.Time, method CalculationMethod) (*models.PnLCalculation, error)
	CreateLotFromTransaction(ctx context.Context, transaction *models.Transaction, tokenID uuid.UUID, quantity, priceUSD string) error
	GetPnLExportData(ctx context.Context, walletAddress string, from, to time.Time, method CalculationMethod) ([]models.PnLExportData, error)
}

type service struct {
	pnlRepo     Repository
	walletRepo  repos.WalletRepository
	tokenRepo   repos.TokenRepository
	calculator  *Calculator
}

func NewService(pnlRepo Repository, walletRepo repos.WalletRepository, tokenRepo repos.TokenRepository) Service {
	return &service{
		pnlRepo:    pnlRepo,
		walletRepo: walletRepo,
		tokenRepo:  tokenRepo,
		calculator: NewCalculator(FIFO), // default calculator
	}
}

func (s *service) CalculatePnL(ctx context.Context, walletAddress string, from, to time.Time, method CalculationMethod) (*models.PnLCalculation, error) {
	// Get wallet ID
	wallet, err := s.walletRepo.GetByAddress(ctx, walletAddress, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	// Get all tokens for the wallet
	tokenIDs, err := s.pnlRepo.GetWalletTokens(ctx, wallet.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet tokens: %w", err)
	}

	if len(tokenIDs) == 0 {
		return nil, fmt.Errorf("no tokens found for wallet")
	}

	// For now, calculate for the first token
	// In a full implementation, this would aggregate across all tokens
	return s.CalculatePnLByToken(ctx, walletAddress, "", from, to, method)
}

func (s *service) CalculatePnLByToken(ctx context.Context, walletAddress, tokenAddress string, from, to time.Time, method CalculationMethod) (*models.PnLCalculation, error) {
	// Get wallet
	wallet, err := s.walletRepo.GetByAddress(ctx, walletAddress, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	// Get token
	var token *models.Token
	if tokenAddress != "" {
		token, err = s.tokenRepo.GetByAddress(ctx, tokenAddress, wallet.ChainID)
		if err != nil {
			return nil, fmt.Errorf("failed to get token: %w", err)
		}
	} else {
		// Get first token for wallet
		tokenIDs, err := s.pnlRepo.GetWalletTokens(ctx, wallet.ID)
		if err != nil || len(tokenIDs) == 0 {
			return nil, fmt.Errorf("no tokens found for wallet")
		}
		token, err = s.tokenRepo.GetByID(ctx, tokenIDs[0])
		if err != nil {
			return nil, fmt.Errorf("failed to get token: %w", err)
		}
	}

	// Get lots for the time period
	lots, err := s.pnlRepo.GetLotsByWallet(ctx, wallet.ID, token.ID, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to get lots: %w", err)
	}

	if len(lots) == 0 {
		return nil, fmt.Errorf("no lots found for the specified period")
	}

	// Create calculator with specified method
	calculator := NewCalculator(method)

	// Get current price (assuming it's stored in the token)
	currentPriceUSD := "0"
	if token.PriceUSD != nil {
		currentPriceUSD = fmt.Sprintf("%.10f", *token.PriceUSD)
	}

	// Calculate PnL
	calculation, err := calculator.CalculatePnL(lots, currentPriceUSD)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate PnL: %w", err)
	}

	// Populate missing fields
	calculation.WalletAddress = walletAddress
	calculation.TokenAddress = token.Address
	calculation.TokenSymbol = token.Symbol

	return calculation, nil
}

func (s *service) CreateLotFromTransaction(ctx context.Context, transaction *models.Transaction, tokenID uuid.UUID, quantity, priceUSD string) error {
	// Get wallet ID from transaction
	wallet, err := s.walletRepo.GetByAddress(ctx, transaction.FromAddress, 1)
	if err != nil {
		return fmt.Errorf("failed to get wallet: %w", err)
	}

	// Determine transaction type for PnL purposes
	lotType := s.determineLotType(transaction)

	lot := &models.PnLLot{
		ID:                uuid.New(),
		WalletID:          wallet.ID,
		TokenID:           tokenID,
		TransactionHash:   transaction.Hash,
		ChainID:           transaction.ChainID,
		Type:              lotType,
		Quantity:          quantity,
		PriceUSD:          priceUSD,
		RemainingQuantity: quantity, // Initially, all quantity is remaining
		BlockNumber:       *transaction.BlockNumber,
		Timestamp:         transaction.Timestamp,
	}

	return s.pnlRepo.CreateLot(ctx, lot)
}

func (s *service) GetPnLExportData(ctx context.Context, walletAddress string, from, to time.Time, method CalculationMethod) ([]models.PnLExportData, error) {
	// Get wallet
	wallet, err := s.walletRepo.GetByAddress(ctx, walletAddress, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	// Get all tokens for the wallet
	tokenIDs, err := s.pnlRepo.GetWalletTokens(ctx, wallet.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet tokens: %w", err)
	}

	var exportData []models.PnLExportData

	for _, tokenID := range tokenIDs {
		// Get token info
		token, err := s.tokenRepo.GetByID(ctx, tokenID)
		if err != nil {
			continue // Skip this token if we can't get its info
		}

		// Get lots for this token
		lots, err := s.pnlRepo.GetLotsByWallet(ctx, wallet.ID, tokenID, from, to)
		if err != nil {
			continue // Skip this token if we can't get its lots
		}

		// Create calculator with specified method
		calculator := NewCalculator(method)
		
		// Get current price
		currentPriceUSD := "0"
		if token.PriceUSD != nil {
			currentPriceUSD = fmt.Sprintf("%.10f", *token.PriceUSD)
		}

		// Calculate PnL to get processed lots
		calculation, err := calculator.CalculatePnL(lots, currentPriceUSD)
		if err != nil {
			continue // Skip this token if calculation fails
		}

		// Convert lots to export data
		for _, lot := range calculation.Lots {
			exportData = append(exportData, models.PnLExportData{
				WalletAddress:     walletAddress,
				TokenSymbol:       token.Symbol,
				TokenAddress:      token.Address,
				TransactionHash:   lot.TransactionHash,
				Type:              lot.Type,
				Quantity:          lot.Quantity,
				PriceUSD:          lot.PriceUSD,
				RemainingQuantity: lot.RemainingQuantity,
				RealizedPnLUSD:    "0", // This would need to be calculated per lot
				Timestamp:         lot.Timestamp,
				BlockNumber:       lot.BlockNumber,
			})
		}
	}

	return exportData, nil
}

// determineLotType maps transaction types to PnL lot types
func (s *service) determineLotType(transaction *models.Transaction) string {
	switch transaction.Type {
	case "receive", "swap":
		return "buy"
	case "send":
		return "sell"
	default:
		return "buy" // Default to buy for other transaction types
	}
}