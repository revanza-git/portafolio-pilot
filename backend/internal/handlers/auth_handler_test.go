package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/defi-dashboard/backend/internal/models"
	"github.com/defi-dashboard/backend/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock services
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) GenerateNonce(ctx context.Context, address string) (string, error) {
	args := m.Called(ctx, address)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) VerifySignature(ctx context.Context, message map[string]interface{}, signature string) (string, error) {
	args := m.Called(ctx, message, signature)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) ValidateToken(tokenString string) (*models.User, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockAuthService) GetUserByAddress(ctx context.Context, address string) (*models.User, error) {
	args := m.Called(ctx, address)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

type MockSIWEService struct {
	mock.Mock
}

func (m *MockSIWEService) GenerateNonce(ctx context.Context, address string) (string, error) {
	args := m.Called(ctx, address)
	return args.String(0), args.Error(1)
}

func (m *MockSIWEService) VerifySignature(ctx context.Context, message, signature string) (*models.User, error) {
	args := m.Called(ctx, message, signature)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockSIWEService) GenerateSIWEMessage(address, nonce string) (string, error) {
	args := m.Called(address, nonce)
	return args.String(0), args.Error(1)
}

func (m *MockSIWEService) ValidateAddress(address string) bool {
	args := m.Called(address)
	return args.Bool(0)
}

func setupTestApp() *fiber.App {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Next()
	})
	return app
}

func TestAuthHandler_GetNonce(t *testing.T) {
	app := setupTestApp()
	authService := new(MockAuthService)
	siweService := new(MockSIWEService)
	
	handler := NewAuthHandler(authService, siweService, "test-secret", 24)
	app.Post("/nonce", handler.GetNonce)

	t.Run("Valid address", func(t *testing.T) {
		address := "0x742d35Cc6573C42c8Ee90b4E43e04c1Fe9E2395d"
		testNonce := "test-nonce-123"
		testMessage := "localhost wants you to sign in with your Ethereum account"

		siweService.On("ValidateAddress", address).Return(true)
		siweService.On("GenerateNonce", mock.Anything, address).Return(testNonce, nil)
		siweService.On("GenerateSIWEMessage", address, testNonce).Return(testMessage, nil)

		reqBody := NonceRequest{Address: address}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/nonce", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		var response NonceResponse
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, testNonce, response.Nonce)
		assert.Equal(t, testMessage, response.Message)

		siweService.AssertExpectations(t)
	})

	t.Run("Invalid address", func(t *testing.T) {
		address := "invalid-address"

		siweService.On("ValidateAddress", address).Return(false)

		reqBody := NonceRequest{Address: address}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/nonce", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("Missing address", func(t *testing.T) {
		reqBody := NonceRequest{}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/nonce", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})
}

func TestAuthHandler_Verify(t *testing.T) {
	app := setupTestApp()
	authService := new(MockAuthService)
	siweService := new(MockSIWEService)
	
	handler := NewAuthHandler(authService, siweService, "test-secret", 24)
	app.Post("/verify", handler.Verify)

	t.Run("Valid verification", func(t *testing.T) {
		testMessage := "localhost wants you to sign in with your Ethereum account"
		testSignature := "0x1234567890abcdef"
		testUser := &models.User{
			ID:      uuid.New(),
			Address: "0x742d35Cc6573C42c8Ee90b4E43e04c1Fe9E2395d",
		}

		siweService.On("VerifySignature", mock.Anything, testMessage, testSignature).Return(testUser, nil)

		reqBody := VerifyRequest{
			Message:   testMessage,
			Signature: testSignature,
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/verify", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		var response AuthResponse
		json.NewDecoder(resp.Body).Decode(&response)
		assert.NotEmpty(t, response.Token)
		assert.Equal(t, testUser.Address, response.Address)
		assert.Equal(t, 24*3600, response.ExpiresIn) // 24 hours in seconds

		siweService.AssertExpectations(t)
	})

	t.Run("Missing message", func(t *testing.T) {
		reqBody := VerifyRequest{
			Signature: "0x1234567890abcdef",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/verify", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("Missing signature", func(t *testing.T) {
		reqBody := VerifyRequest{
			Message: "test message",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/verify", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})
}

func TestAuthHandler_SendMagicLink(t *testing.T) {
	app := setupTestApp()
	authService := new(MockAuthService)
	siweService := new(MockSIWEService)
	
	handler := NewAuthHandler(authService, siweService, "test-secret", 24)
	app.Post("/magic-link", handler.SendMagicLink)

	t.Run("Valid email", func(t *testing.T) {
		reqBody := MagicLinkRequest{
			Email: "test@example.com",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/magic-link", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		var response MagicLinkResponse
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Magic link sent to your email", response.Message)
		assert.NotEmpty(t, response.MagicLink)
	})

	t.Run("Missing email", func(t *testing.T) {
		reqBody := MagicLinkRequest{}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/magic-link", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})
}

func TestAuthHandler_GenerateJWT(t *testing.T) {
	authService := new(MockAuthService)
	siweService := new(MockSIWEService)
	
	handler := NewAuthHandler(authService, siweService, "test-secret", 24)

	testUser := &models.User{
		ID:        uuid.New(),
		Address:   "0x742d35Cc6573C42c8Ee90b4E43e04c1Fe9E2395d",
		CreatedAt: time.Now(),
	}

	token, err := handler.generateJWT(testUser)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Token should be a valid JWT format (3 parts separated by dots)
	parts := len(bytes.Split([]byte(token), []byte(".")))
	assert.Equal(t, 3, parts)
}