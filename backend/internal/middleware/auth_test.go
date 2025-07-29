package middleware

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupTestApp() *fiber.App {
	app := fiber.New()
	return app
}

func generateTestJWT(secret, address string, expiry time.Duration) (string, error) {
	claims := &Claims{
		Address: address,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func TestJWTAuth_ValidToken(t *testing.T) {
	app := setupTestApp()
	secret := "test-secret"
	address := "0x742d35Cc6573C42c8Ee90b4E43e04c1Fe9E2395d"

	// Create test endpoint
	app.Get("/protected", JWTAuth(secret), func(c *fiber.Ctx) error {
		userAddress := c.Locals("address").(string)
		return c.JSON(fiber.Map{"address": userAddress})
	})

	// Generate valid token
	token, err := generateTestJWT(secret, address, time.Hour)
	assert.NoError(t, err)

	// Test with valid token
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestJWTAuth_MissingAuthHeader(t *testing.T) {
	app := setupTestApp()
	secret := "test-secret"

	// Create test endpoint
	app.Get("/protected", JWTAuth(secret), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})

	// Test without Authorization header
	req := httptest.NewRequest("GET", "/protected", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestJWTAuth_InvalidAuthHeaderFormat(t *testing.T) {
	app := setupTestApp()
	secret := "test-secret"

	// Create test endpoint
	app.Get("/protected", JWTAuth(secret), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})

	testCases := []struct {
		name   string
		header string
	}{
		{"Missing Bearer prefix", "token123"},
		{"Wrong prefix", "Basic token123"},
		{"No token", "Bearer"},
		{"Extra parts", "Bearer token123 extra"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/protected", nil)
			req.Header.Set("Authorization", tc.header)

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, 401, resp.StatusCode)
		})
	}
}

func TestJWTAuth_ExpiredToken(t *testing.T) {
	app := setupTestApp()
	secret := "test-secret"
	address := "0x742d35Cc6573C42c8Ee90b4E43e04c1Fe9E2395d"

	// Create test endpoint
	app.Get("/protected", JWTAuth(secret), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})

	// Generate expired token (expired 1 hour ago)
	token, err := generateTestJWT(secret, address, -time.Hour)
	assert.NoError(t, err)

	// Test with expired token
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestJWTAuth_InvalidSecret(t *testing.T) {
	app := setupTestApp()
	secret := "test-secret"
	wrongSecret := "wrong-secret"
	address := "0x742d35Cc6573C42c8Ee90b4E43e04c1Fe9E2395d"

	// Create test endpoint
	app.Get("/protected", JWTAuth(secret), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})

	// Generate token with wrong secret
	token, err := generateTestJWT(wrongSecret, address, time.Hour)
	assert.NoError(t, err)

	// Test with token signed with wrong secret
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestJWTAuth_MalformedToken(t *testing.T) {
	app := setupTestApp()
	secret := "test-secret"

	// Create test endpoint
	app.Get("/protected", JWTAuth(secret), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})

	testCases := []string{
		"invalid.token",
		"invalid",
		"",
		"header.payload", // Missing signature
	}

	for _, token := range testCases {
		t.Run("Token: "+token, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/protected", nil)
			req.Header.Set("Authorization", "Bearer "+token)

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, 401, resp.StatusCode)
		})
	}
}

func TestOptionalAuth_ValidToken(t *testing.T) {
	app := setupTestApp()
	secret := "test-secret"
	address := "0x742d35Cc6573C42c8Ee90b4E43e04c1Fe9E2395d"

	// Create test endpoint
	app.Get("/optional", OptionalAuth(secret), func(c *fiber.Ctx) error {
		userAddress := c.Locals("address")
		if userAddress != nil {
			return c.JSON(fiber.Map{"authenticated": true, "address": userAddress})
		}
		return c.JSON(fiber.Map{"authenticated": false})
	})

	// Generate valid token
	token, err := generateTestJWT(secret, address, time.Hour)
	assert.NoError(t, err)

	// Test with valid token
	req := httptest.NewRequest("GET", "/optional", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestOptionalAuth_NoToken(t *testing.T) {
	app := setupTestApp()
	secret := "test-secret"

	// Create test endpoint
	app.Get("/optional", OptionalAuth(secret), func(c *fiber.Ctx) error {
		userAddress := c.Locals("address")
		if userAddress != nil {
			return c.JSON(fiber.Map{"authenticated": true, "address": userAddress})
		}
		return c.JSON(fiber.Map{"authenticated": false})
	})

	// Test without token
	req := httptest.NewRequest("GET", "/optional", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestOptionalAuth_InvalidToken(t *testing.T) {
	app := setupTestApp()
	secret := "test-secret"

	// Create test endpoint
	app.Get("/optional", OptionalAuth(secret), func(c *fiber.Ctx) error {
		userAddress := c.Locals("address")
		if userAddress != nil {
			return c.JSON(fiber.Map{"authenticated": true, "address": userAddress})
		}
		return c.JSON(fiber.Map{"authenticated": false})
	})

	// Test with invalid token - should continue without authentication
	req := httptest.NewRequest("GET", "/optional", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestClaims_Structure(t *testing.T) {
	claims := &Claims{
		Address: "0x742d35Cc6573C42c8Ee90b4E43e04c1Fe9E2395d",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
	}

	// Test that claims can be used to generate a token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("test-secret"))
	
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	// Test parsing the token back
	parsedToken, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("test-secret"), nil
	})

	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	parsedClaims, ok := parsedToken.Claims.(*Claims)
	assert.True(t, ok)
	assert.Equal(t, claims.Address, parsedClaims.Address)
}