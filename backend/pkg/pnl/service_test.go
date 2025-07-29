package pnl

import (
	"context"
	"testing"
	"time"

	"github.com/defi-dashboard/backend/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock implementations
type MockPnLRepository struct {
	mock.Mock
}

func (m *MockPnLRepository) CreateLot(ctx context.Context, lot *models.PnLLot) error {
	args := m.Called(ctx, lot)
	return args.Error(0)
}

func (m *MockPnLRepository) GetLotsByWallet(ctx context.Context, walletID uuid.UUID, tokenID uuid.UUID, from, to time.Time) ([]models.PnLLot, error) {
	args := m.Called(ctx, walletID, tokenID, from, to)
	return args.Get(0).([]models.PnLLot), args.Error(1)
}

func (m *MockPnLRepository) GetLotsByWalletAndToken(ctx context.Context, walletID uuid.UUID, tokenID uuid.UUID) ([]models.PnLLot, error) {
	args := m.Called(ctx, walletID, tokenID)
	return args.Get(0).([]models.PnLLot), args.Error(1)
}

func (m *MockPnLRepository) UpdateLotRemainingQuantity(ctx context.Context, lotID uuid.UUID, remainingQuantity string) error {
	args := m.Called(ctx, lotID, remainingQuantity)
	return args.Error(0)
}

func (m *MockPnLRepository) GetWalletTokens(ctx context.Context, walletID uuid.UUID) ([]uuid.UUID, error) {
	args := m.Called(ctx, walletID)
	return args.Get(0).([]uuid.UUID), args.Error(1)
}

type MockWalletRepository struct {
	mock.Mock
}

func (m *MockWalletRepository) GetByAddress(ctx context.Context, address string) (*models.Wallet, error) {
	args := m.Called(ctx, address)
	return args.Get(0).(*models.Wallet), args.Error(1)
}

func (m *MockWalletRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Wallet, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Wallet), args.Error(1)
}

func (m *MockWalletRepository) Create(ctx context.Context, wallet *models.Wallet) error {
	args := m.Called(ctx, wallet)
	return args.Error(0)
}

func (m *MockWalletRepository) Update(ctx context.Context, wallet *models.Wallet) error {
	args := m.Called(ctx, wallet)
	return args.Error(0)
}

func (m *MockWalletRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockWalletRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.Wallet, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Wallet), args.Error(1)
}

type MockTokenRepository struct {
	mock.Mock
}

func (m *MockTokenRepository) GetByAddress(ctx context.Context, address string, chainID int) (*models.Token, error) {
	args := m.Called(ctx, address, chainID)
	return args.Get(0).(*models.Token), args.Error(1)
}

func (m *MockTokenRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Token, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Token), args.Error(1)
}

func (m *MockTokenRepository) Create(ctx context.Context, token *models.Token) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockTokenRepository) Update(ctx context.Context, token *models.Token) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockTokenRepository) GetBySymbol(ctx context.Context, symbol string, chainID int) (*models.Token, error) {
	args := m.Called(ctx, symbol, chainID)
	return args.Get(0).(*models.Token), args.Error(1)
}

func (m *MockTokenRepository) List(ctx context.Context, chainID int, limit, offset int) ([]models.Token, error) {
	args := m.Called(ctx, chainID, limit, offset)
	return args.Get(0).([]models.Token), args.Error(1)
}

func (m *MockTokenRepository) UpdatePrices(ctx context.Context, prices map[string]float64) error {
	args := m.Called(ctx, prices)
	return args.Error(0)
}

func TestService_CalculatePnLByToken(t *testing.T) {
	ctx := context.Background()
	
	// Setup mocks
	mockPnLRepo := new(MockPnLRepository)
	mockWalletRepo := new(MockWalletRepository)
	mockTokenRepo := new(MockTokenRepository)
	
	service := &service{
		pnlRepo:    mockPnLRepo,
		walletRepo: mockWalletRepo,
		tokenRepo:  mockTokenRepo,
		calculator: NewCalculator(FIFO),
	}

	walletID := uuid.New()
	tokenID := uuid.New()
	walletAddress := "0x1234567890123456789012345678901234567890"
	tokenAddress := "0x0987654321098765432109876543210987654321"
	
	wallet := &models.Wallet{
		ID:      walletID,
		Address: walletAddress,
		ChainID: 1,
	}
	
	price := 15.50
	token := &models.Token{
		ID:       tokenID,
		Address:  tokenAddress,
		Symbol:   "TEST",
		PriceUSD: &price,
	}

	lots := []models.PnLLot{
		{
			ID:                uuid.New(),
			WalletID:          walletID,
			TokenID:           tokenID,
			Type:              "buy",
			Quantity:          "100",
			PriceUSD:          "10.00",
			RemainingQuantity: "100",
			Timestamp:         time.Now().Add(-time.Hour * 24),
		},
		{
			ID:                uuid.New(),
			WalletID:          walletID,
			TokenID:           tokenID,
			Type:              "sell",
			Quantity:          "50",
			PriceUSD:          "12.00",
			RemainingQuantity: "50",
			Timestamp:         time.Now().Add(-time.Hour * 12),
		},
	}

	from := time.Now().Add(-time.Hour * 48)
	to := time.Now()

	// Setup mock expectations
	mockWalletRepo.On("GetByAddress", ctx, walletAddress).Return(wallet, nil)
	mockTokenRepo.On("GetByAddress", ctx, tokenAddress, wallet.ChainID).Return(token, nil)
	mockPnLRepo.On("GetLotsByWallet", ctx, walletID, tokenID, from, to).Return(lots, nil)

	// Execute test
	result, err := service.CalculatePnLByToken(ctx, walletAddress, tokenAddress, from, to, FIFO)

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, walletAddress, result.WalletAddress)
	assert.Equal(t, tokenAddress, result.TokenAddress)
	assert.Equal(t, "TEST", result.TokenSymbol)
	assert.Equal(t, "fifo", result.Method)
	assert.Equal(t, "100", result.RealizedPnLUSD) // (12-10)*50 = 100
	assert.Equal(t, "50", result.CurrentQuantity)  // 100-50 = 50

	// Verify mock expectations
	mockWalletRepo.AssertExpectations(t)
	mockTokenRepo.AssertExpectations(t)
	mockPnLRepo.AssertExpectations(t)
}

func TestService_CreateLotFromTransaction(t *testing.T) {
	ctx := context.Background()
	
	// Setup mocks
	mockPnLRepo := new(MockPnLRepository)
	mockWalletRepo := new(MockWalletRepository)
	mockTokenRepo := new(MockTokenRepository)
	
	service := &service{
		pnlRepo:    mockPnLRepo,
		walletRepo: mockWalletRepo,
		tokenRepo:  mockTokenRepo,
		calculator: NewCalculator(FIFO),
	}

	walletID := uuid.New()
	tokenID := uuid.New()
	blockNumber := int64(12345)
	
	wallet := &models.Wallet{
		ID:      walletID,
		Address: "0x1234567890123456789012345678901234567890",
		ChainID: 1,
	}
	
	transaction := &models.Transaction{
		Hash:        "0xabcdef1234567890",
		FromAddress: wallet.Address,
		ChainID:     1,
		Type:        "receive",
		BlockNumber: &blockNumber,
		Timestamp:   time.Now(),
	}

	// Setup mock expectations
	mockWalletRepo.On("GetByAddress", ctx, transaction.FromAddress).Return(wallet, nil)
	mockPnLRepo.On("CreateLot", ctx, mock.AnythingOfType("*models.PnLLot")).Return(nil)

	// Execute test
	err := service.CreateLotFromTransaction(ctx, transaction, tokenID, "100", "10.50")

	// Assertions
	require.NoError(t, err)

	// Verify mock expectations
	mockWalletRepo.AssertExpectations(t)
	mockPnLRepo.AssertExpectations(t)
}

func TestService_DetermineLotType(t *testing.T) {
	service := &service{}

	tests := []struct {
		transactionType string
		expectedLotType string
	}{
		{"receive", "buy"},
		{"swap", "buy"},
		{"send", "sell"},
		{"approve", "buy"}, // default
		{"unknown", "buy"}, // default
	}

	for _, tt := range tests {
		t.Run(tt.transactionType, func(t *testing.T) {
			transaction := &models.Transaction{Type: tt.transactionType}
			result := service.determineLotType(transaction)
			assert.Equal(t, tt.expectedLotType, result)
		})
	}
}

func TestService_EdgeCases(t *testing.T) {
	ctx := context.Background()
	
	// Setup mocks
	mockPnLRepo := new(MockPnLRepository)
	mockWalletRepo := new(MockWalletRepository)
	mockTokenRepo := new(MockTokenRepository)
	
	service := &service{
		pnlRepo:    mockPnLRepo,
		walletRepo: mockWalletRepo,
		tokenRepo:  mockTokenRepo,
		calculator: NewCalculator(FIFO),
	}

	t.Run("Wallet not found", func(t *testing.T) {
		walletAddress := "0x1234567890123456789012345678901234567890"
		from := time.Now().Add(-time.Hour * 48)
		to := time.Now()

		mockWalletRepo.On("GetByAddress", ctx, walletAddress).Return((*models.Wallet)(nil), assert.AnError)

		_, err := service.CalculatePnLByToken(ctx, walletAddress, "", from, to, FIFO)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get wallet")
	})

	t.Run("No lots found", func(t *testing.T) {
		walletID := uuid.New()
		tokenID := uuid.New()
		walletAddress := "0x1234567890123456789012345678901234567890"
		tokenAddress := "0x0987654321098765432109876543210987654321"
		
		wallet := &models.Wallet{
			ID:      walletID,
			Address: walletAddress,
			ChainID: 1,
		}
		
		token := &models.Token{
			ID:      tokenID,
			Address: tokenAddress,
			Symbol:  "TEST",
		}

		from := time.Now().Add(-time.Hour * 48)
		to := time.Now()

		mockWalletRepo.On("GetByAddress", ctx, walletAddress).Return(wallet, nil)
		mockTokenRepo.On("GetByAddress", ctx, tokenAddress, wallet.ChainID).Return(token, nil)
		mockPnLRepo.On("GetLotsByWallet", ctx, walletID, tokenID, from, to).Return([]models.PnLLot{}, nil)

		_, err := service.CalculatePnLByToken(ctx, walletAddress, tokenAddress, from, to, FIFO)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no lots found")
	})

	t.Run("Token not found", func(t *testing.T) {
		walletID := uuid.New()
		walletAddress := "0x1234567890123456789012345678901234567890"
		tokenAddress := "0x0987654321098765432109876543210987654321"
		
		wallet := &models.Wallet{
			ID:      walletID,
			Address: walletAddress,
			ChainID: 1,
		}

		from := time.Now().Add(-time.Hour * 48)
		to := time.Now()

		mockWalletRepo.On("GetByAddress", ctx, walletAddress).Return(wallet, nil)
		mockTokenRepo.On("GetByAddress", ctx, tokenAddress, wallet.ChainID).Return((*models.Token)(nil), assert.AnError)

		_, err := service.CalculatePnLByToken(ctx, walletAddress, tokenAddress, from, to, FIFO)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get token")
	})
}