package services

import (
	"context"
	"testing"
	"time"

	"github.com/defi-dashboard/backend/internal/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repositories
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByAddress(ctx context.Context, address string) (*models.User, error) {
	args := m.Called(ctx, address)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, address, nonce string) (*models.User, error) {
	args := m.Called(ctx, address, nonce)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) UpdateNonce(ctx context.Context, address, nonce string) (*models.User, error) {
	args := m.Called(ctx, address, nonce)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID, lastLogin time.Time) error {
	args := m.Called(ctx, id, lastLogin)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateEmail(ctx context.Context, id uuid.UUID, email string) (*models.User, error) {
	args := m.Called(ctx, id, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockNonceRepository struct {
	mock.Mock
}

func (m *MockNonceRepository) Store(ctx context.Context, address, nonce string, expiresAt time.Time) error {
	args := m.Called(ctx, address, nonce, expiresAt)
	return args.Error(0)
}

func (m *MockNonceRepository) ValidateAndUse(ctx context.Context, address, nonce string) (bool, error) {
	args := m.Called(ctx, address, nonce)
	return args.Bool(0), args.Error(1)
}

func (m *MockNonceRepository) CleanupExpired(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockNonceRepository) GetByAddressAndNonce(ctx context.Context, address, nonce string) (*models.NonceStorage, error) {
	args := m.Called(ctx, address, nonce)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.NonceStorage), args.Error(1)
}

func TestSIWEService_GenerateNonce(t *testing.T) {
	userRepo := new(MockUserRepository)
	nonceRepo := new(MockNonceRepository)
	service := NewSIWEService(userRepo, nonceRepo, "localhost")

	ctx := context.Background()
	address := "0x742d35Cc6573C42c8Ee90b4E43e04c1Fe9E2395d"

	// Mock successful nonce storage
	nonceRepo.On("Store", ctx, address, mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).Return(nil)

	nonce, err := service.GenerateNonce(ctx, address)

	assert.NoError(t, err)
	assert.NotEmpty(t, nonce)
	assert.Len(t, nonce, 32) // 16 bytes hex encoded = 32 characters
	userRepo.AssertExpectations(t)
	nonceRepo.AssertExpectations(t)
}

func TestSIWEService_GenerateNonce_InvalidAddress(t *testing.T) {
	userRepo := new(MockUserRepository)
	nonceRepo := new(MockNonceRepository)
	service := NewSIWEService(userRepo, nonceRepo, "localhost")

	ctx := context.Background()
	invalidAddress := "invalid-address"

	_, err := service.GenerateNonce(ctx, invalidAddress)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid Ethereum address")
}

func TestSIWEService_ValidateAddress(t *testing.T) {
	userRepo := new(MockUserRepository)
	nonceRepo := new(MockNonceRepository)
	service := NewSIWEService(userRepo, nonceRepo, "localhost")

	tests := []struct {
		name     string
		address  string
		expected bool
	}{
		{"Valid address", "0x742d35Cc6573C42c8Ee90b4E43e04c1Fe9E2395d", true},
		{"Valid address lowercase", "0x742d35cc6573c42c8ee90b4e43e04c1fe9e2395d", true},
		{"Invalid address", "invalid-address", false},
		{"Empty address", "", false},
		{"Address without 0x", "742d35Cc6573C42c8Ee90b4E43e04c1Fe9E2395d", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.ValidateAddress(tt.address)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSIWEService_GenerateSIWEMessage(t *testing.T) {
	userRepo := new(MockUserRepository)
	nonceRepo := new(MockNonceRepository)
	service := NewSIWEService(userRepo, nonceRepo, "localhost")

	address := "0x742d35Cc6573C42c8Ee90b4E43e04c1Fe9E2395d"
	nonce := "test-nonce-123"

	message, err := service.GenerateSIWEMessage(address, nonce)

	assert.NoError(t, err)
	assert.NotEmpty(t, message)
	assert.Contains(t, message, "localhost")
	assert.Contains(t, message, address)
	assert.Contains(t, message, nonce)
	assert.Contains(t, message, "Sign in to DeFi Portfolio Dashboard")
}

func TestSIWEService_VerifyEthereumSignature(t *testing.T) {
	userRepo := new(MockUserRepository)
	nonceRepo := new(MockNonceRepository)
	service := NewSIWEService(userRepo, nonceRepo, "localhost")

	// Create a test private key and address
	privateKey, err := crypto.GenerateKey()
	assert.NoError(t, err)
	
	testAddress := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	message := "Test message for signing"

	// Sign the message
	messageHash := crypto.Keccak256Hash([]byte("\\x19Ethereum Signed Message:\\n" + "25" + message))
	signature, err := crypto.Sign(messageHash.Bytes(), privateKey)
	assert.NoError(t, err)
	
	// Convert to hex with recovery ID adjustment for MetaMask compatibility
	signature[64] += 27
	signatureHex := common.Bytes2Hex(signature)

	// Test signature verification
	err = service.verifyEthereumSignature(message, signatureHex, testAddress)
	assert.NoError(t, err)

	// Test with wrong address
	wrongAddress := "0x742d35Cc6573C42c8Ee90b4E43e04c1Fe9E2395d"
	err = service.verifyEthereumSignature(message, signatureHex, wrongAddress)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "signature address mismatch")

	// Test with invalid signature
	err = service.verifyEthereumSignature(message, "invalid-signature", testAddress)
	assert.Error(t, err)
}