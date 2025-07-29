package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/defi-dashboard/backend/internal/clients"
	"github.com/defi-dashboard/backend/internal/clients/bridge"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLiFiClient_GetQuote(t *testing.T) {
	// Mock LI.FI API response
	mockResponse := map[string]interface{}{
		"routes": []map[string]interface{}{
			{
				"id":           "test-route-123",
				"fromChainId":  1,
				"toChainId":    137,
				"fromToken": map[string]interface{}{
					"address":  "0xA0b86a33E6441e6e80fd4e3Cd9Cc5F7b",
					"symbol":   "USDC",
					"name":     "USD Coin",
					"decimals": 6,
					"chainId":  1,
					"logoURI":  "https://example.com/usdc.png",
				},
				"toToken": map[string]interface{}{
					"address":  "0x2791Bca1f2de4661ED88A30C99A7a9449Aa84174",
					"symbol":   "USDC",
					"name":     "USD Coin",
					"decimals": 6,
					"chainId":  137,
					"logoURI":  "https://example.com/usdc-polygon.png",
				},
				"fromAmount": "1000000",
				"toAmount":   "990000",
				"steps": []map[string]interface{}{
					{
						"id":   "step-1",
						"type": "bridge",
						"tool": "stargate",
						"action": map[string]interface{}{
							"fromChainId": 1,
							"toChainId":   137,
							"fromToken": map[string]interface{}{
								"address":  "0xA0b86a33E6441e6e80fd4e3Cd9Cc5F7b",
								"symbol":   "USDC",
								"name":     "USD Coin",
								"decimals": 6,
								"chainId":  1,
								"logoURI":  "https://example.com/usdc.png",
							},
							"toToken": map[string]interface{}{
								"address":  "0x2791Bca1f2de4661ED88A30C99A7a9449Aa84174",
								"symbol":   "USDC",
								"name":     "USD Coin",
								"decimals": 6,
								"chainId":  137,
								"logoURI":  "https://example.com/usdc-polygon.png",
							},
							"fromAmount": "1000000",
							"toAmount":   "990000",
							"slippage":   0.5,
						},
						"estimate": map[string]interface{}{
							"fromAmount":        "1000000",
							"toAmount":          "990000",
							"executionDuration": 300,
							"feeCosts":         []interface{}{},
							"gasCosts":         []interface{}{},
						},
						"transactionRequest": map[string]interface{}{
							"data":     "0x1234567890abcdef",
							"to":       "0x1234567890123456789012345678901234567890",
							"value":    "0",
							"chainId":  1,
							"gasLimit": "150000",
							"gasPrice": "20000000000",
						},
					},
				},
				"tags":       []string{"RECOMMENDED"},
				"gasCostUSD": "5.50",
			},
		},
	}

	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/quote", r.URL.Path)
		
		// Check query parameters
		query := r.URL.Query()
		assert.Equal(t, "1", query.Get("fromChain"))
		assert.Equal(t, "137", query.Get("toChain"))
		assert.Equal(t, "0xA0b86a33E6441e6e80fd4e3Cd9Cc5F7b", query.Get("fromToken"))
		assert.Equal(t, "0x2791Bca1f2de4661ED88A30C99A7a9449Aa84174", query.Get("toToken"))
		assert.Equal(t, "1000000", query.Get("fromAmount"))
		assert.Equal(t, "0x1234567890123456789012345678901234567890", query.Get("fromAddress"))
		assert.Equal(t, "0x1234567890123456789012345678901234567890", query.Get("toAddress"))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Create client with mock server URL
	config := clients.ClientConfig{
		BaseURL:    server.URL,
		APIKey:     "test-api-key",
		Timeout:    5 * time.Second,
		MaxRetries: 1,
		RetryDelay: 100 * time.Millisecond,
		RateLimit: clients.RateLimitConfig{
			RequestsPerSecond: 10,
			BurstSize:         20,
		},
	}

	client := bridge.NewLiFiClient(config)

	// Test request
	req := clients.QuoteRequest{
		FromChainID: "1",
		ToChainID:   "137",
		FromToken:   "0xA0b86a33E6441e6e80fd4e3Cd9Cc5F7b",
		ToToken:     "0x2791Bca1f2de4661ED88A30C99A7a9449Aa84174",
		Amount:      "1000000",
		UserAddress: "0x1234567890123456789012345678901234567890",
		Slippage:    0.5,
	}

	ctx := context.Background()
	quote, err := client.GetQuote(ctx, req)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, quote)

	assert.Equal(t, "test-route-123", quote.ID)
	assert.Equal(t, "bridge", quote.Type)
	assert.Equal(t, "LI.FI", quote.Provider)
	assert.Equal(t, "1", quote.FromChainID)
	assert.Equal(t, "137", quote.ToChainID)
	assert.Equal(t, "1000000", quote.FromAmount)
	assert.Equal(t, "990000", quote.ToAmount)
	
	// Check transaction data
	require.NotNil(t, quote.TransactionData)
	assert.Equal(t, "0x1234567890123456789012345678901234567890", quote.TransactionData.To)
	assert.Equal(t, "0x1234567890abcdef", quote.TransactionData.Data)
	assert.Equal(t, "0", quote.TransactionData.Value)
	assert.Equal(t, "150000", quote.TransactionData.GasLimit)
	assert.Equal(t, "1", quote.TransactionData.ChainID)

	// Check route
	require.Len(t, quote.Route, 1)
	assert.Equal(t, "stargate", quote.Route[0].Protocol)
	assert.Equal(t, "bridge", quote.Route[0].Type)
}

func TestLiFiClient_GetSupportedChains(t *testing.T) {
	// Mock chains response
	mockResponse := map[string]interface{}{
		"chains": []map[string]interface{}{
			{
				"id":        1,
				"key":       "eth",
				"name":      "Ethereum",
				"coin":      "ETH",
				"mainnetId": 1,
				"logoURI":   "https://example.com/ethereum.png",
				"nativeToken": map[string]interface{}{
					"address":  "0x0000000000000000000000000000000000000000",
					"chainId":  1,
					"symbol":   "ETH",
					"name":     "Ethereum",
					"decimals": 18,
					"logoURI":  "https://example.com/eth.png",
				},
			},
			{
				"id":        137,
				"key":       "pol",
				"name":      "Polygon",
				"coin":      "MATIC",
				"mainnetId": 137,
				"logoURI":   "https://example.com/polygon.png",
				"nativeToken": map[string]interface{}{
					"address":  "0x0000000000000000000000000000000000001010",
					"chainId":  137,
					"symbol":   "MATIC",
					"name":     "Polygon",
					"decimals": 18,
					"logoURI":  "https://example.com/matic.png",
				},
			},
		},
	}

	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/chains", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Create client
	config := clients.ClientConfig{
		BaseURL:    server.URL,
		APIKey:     "test-api-key",
		Timeout:    5 * time.Second,
		MaxRetries: 1,
		RetryDelay: 100 * time.Millisecond,
		RateLimit: clients.RateLimitConfig{
			RequestsPerSecond: 10,
			BurstSize:         20,
		},
	}

	client := bridge.NewLiFiClient(config)

	// Test
	ctx := context.Background()
	chains, err := client.GetSupportedChains(ctx)

	// Assertions
	require.NoError(t, err)
	require.Len(t, chains, 2)

	// Check Ethereum
	eth := chains[0]
	assert.Equal(t, "1", eth.ID)
	assert.Equal(t, "Ethereum", eth.Name)
	assert.Equal(t, "https://example.com/ethereum.png", eth.LogoURI)
	assert.False(t, eth.IsTestnet)
	assert.Equal(t, "ETH", eth.NativeCurrency.Symbol)

	// Check Polygon
	polygon := chains[1]
	assert.Equal(t, "137", polygon.ID)
	assert.Equal(t, "Polygon", polygon.Name)
	assert.Equal(t, "https://example.com/polygon.png", polygon.LogoURI)
	assert.False(t, polygon.IsTestnet)
	assert.Equal(t, "MATIC", polygon.NativeCurrency.Symbol)
}

func TestLiFiClient_IsHealthy(t *testing.T) {
	// Test healthy response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/status", r.URL.Path)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := clients.ClientConfig{
		BaseURL:    server.URL,
		APIKey:     "test-api-key",
		Timeout:    5 * time.Second,
		MaxRetries: 1,
		RetryDelay: 100 * time.Millisecond,
		RateLimit: clients.RateLimitConfig{
			RequestsPerSecond: 10,
			BurstSize:         20,
		},
	}

	client := bridge.NewLiFiClient(config)

	ctx := context.Background()
	healthy := client.IsHealthy(ctx)

	assert.True(t, healthy)
}

func TestLiFiClient_GetProviderName(t *testing.T) {
	config := clients.ClientConfig{
		BaseURL: "https://api.test.com",
		APIKey:  "test-key",
	}

	client := bridge.NewLiFiClient(config)
	assert.Equal(t, "LI.FI", client.GetProviderName())
}