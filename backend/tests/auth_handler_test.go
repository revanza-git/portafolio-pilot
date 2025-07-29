package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/defi-dashboard/backend/internal/handlers"
	"github.com/defi-dashboard/backend/internal/models"
	"github.com/defi-dashboard/backend/internal/repos"
	"github.com/defi-dashboard/backend/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repository
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

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestGetNonce(t *testing.T) {
	// Setup
	app := fiber.New()
	mockUserRepo := new(MockUserRepository)
	mockWalletRepo := repos.NewWalletRepository(nil) // Use real repo with nil db for this test
	
	authService := services.NewAuthService(mockUserRepo, mockWalletRepo, "test-secret", 24)
	authHandler := handlers.NewAuthHandler(authService)

	app.Get("/auth/nonce", authHandler.GetNonce)

	tests := []struct {
		name           string
		address        string
		mockSetup      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:    "Valid address - existing user",
			address: "0x1234567890123456789012345678901234567890",
			mockSetup: func() {
				mockUserRepo.On("GetByAddress", mock.Anything, "0x1234567890123456789012345678901234567890").
					Return(&models.User{
						ID:      uuid.New(),
						Address: "0x1234567890123456789012345678901234567890",
						Nonce:   "old-nonce",
					}, nil)
				mockUserRepo.On("UpdateNonce", mock.Anything, "0x1234567890123456789012345678901234567890", mock.Anything).
					Return(&models.User{
						ID:      uuid.New(),
						Address: "0x1234567890123456789012345678901234567890",
						Nonce:   "new-nonce",
					}, nil)
			},
			expectedStatus: 200,
			expectedBody: map[string]interface{}{
				"nonce": "mock-nonce", // Will be checked for existence only
			},
		},
		{
			name:           "Missing address",
			address:        "",
			mockSetup:      func() {},
			expectedStatus: 400,
			expectedBody: map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "Address is required",
			},
		},
		{
			name:           "Invalid address format",
			address:        "invalid-address",
			mockSetup:      func() {},
			expectedStatus: 400,
			expectedBody: map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "Invalid Ethereum address",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mock
			mockUserRepo.ExpectedCalls = nil
			mockUserRepo.Calls = nil
			
			// Setup mock
			tt.mockSetup()

			// Make request
			req := httptest.NewRequest("GET", "/auth/nonce?address="+tt.address, nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)

			// Check status
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			// Check response body
			var body map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&body)
			assert.NoError(t, err)

			if tt.expectedStatus == 200 {
				// For success, just check that nonce exists and is not empty
				assert.NotEmpty(t, body["nonce"])
			} else {
				// For errors, check exact response
				assert.Equal(t, tt.expectedBody["code"], body["code"])
				assert.Equal(t, tt.expectedBody["message"], body["message"])
			}

			// Verify mock calls
			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestVerify(t *testing.T) {
	// Setup
	app := fiber.New()
	mockUserRepo := new(MockUserRepository)
	mockWalletRepo := repos.NewWalletRepository(nil)
	
	authService := services.NewAuthService(mockUserRepo, mockWalletRepo, "test-secret", 24)
	authHandler := handlers.NewAuthHandler(authService)

	app.Post("/auth/verify", authHandler.Verify)

	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func()
		expectedStatus int
		checkResponse  func(t *testing.T, body map[string]interface{})
	}{
		{
			name: "Valid verification",
			requestBody: map[string]interface{}{
				"message": map[string]interface{}{
					"address": "0x1234567890123456789012345678901234567890",
					"nonce":   "test-nonce",
				},
				"signature": "0xabcdef...",
			},
			mockSetup: func() {
				mockUserRepo.On("GetByAddress", mock.Anything, "0x1234567890123456789012345678901234567890").
					Return(&models.User{
						ID:      uuid.New(),
						Address: "0x1234567890123456789012345678901234567890",
						Nonce:   "test-nonce",
					}, nil)
				mockUserRepo.On("UpdateNonce", mock.Anything, "0x1234567890123456789012345678901234567890", mock.Anything).
					Return(&models.User{
						ID:      uuid.New(),
						Address: "0x1234567890123456789012345678901234567890",
						Nonce:   "new-nonce",
					}, nil)
			},
			expectedStatus: 200,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.NotEmpty(t, body["token"])
				assert.Equal(t, float64(86400), body["expiresIn"])
				assert.Equal(t, "0x1234567890123456789012345678901234567890", body["address"])
			},
		},
		{
			name:           "Invalid request body",
			requestBody:    "invalid",
			mockSetup:      func() {},
			expectedStatus: 400,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "BAD_REQUEST", body["code"])
			},
		},
		{
			name: "Missing signature",
			requestBody: map[string]interface{}{
				"message": map[string]interface{}{
					"address": "0x1234567890123456789012345678901234567890",
					"nonce":   "test-nonce",
				},
			},
			mockSetup:      func() {},
			expectedStatus: 400,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "BAD_REQUEST", body["code"])
				assert.Equal(t, "Message and signature are required", body["message"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mock
			mockUserRepo.ExpectedCalls = nil
			mockUserRepo.Calls = nil
			
			// Setup mock
			tt.mockSetup()

			// Prepare request body
			bodyBytes, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/auth/verify", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			// Make request
			resp, err := app.Test(req)
			assert.NoError(t, err)

			// Check status
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			// Check response
			var body map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&body)
			assert.NoError(t, err)
			
			tt.checkResponse(t, body)

			// Verify mock calls
			mockUserRepo.AssertExpectations(t)
		})
	}
}