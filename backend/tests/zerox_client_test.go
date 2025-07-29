package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/defi-dashboard/backend/internal/clients"
	"github.com/defi-dashboard/backend/internal/clients/swap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestZeroXClient_GetQuote(t *testing.T) {
	// Mock 0x API response
	mockResponse := map[string]interface{}{
		"chainId":              1,
		"price":                "1.002345",
		"guaranteedPrice":      "1.000000",
		"estimatedPriceImpact": "0.0234",
		"to":                   "0x1234567890123456789012345678901234567890",
		"data":                 "0xabcdef1234567890",
		"value":                "0",
		"gas":                  "120000",
		"estimatedGas":         "115000",
		"gasPrice":             "20000000000",
		"protocolFee":          "1000",
		"minimumProtocolFee":   "500",
		"buyTokenAddress":      "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		"sellTokenAddress":     "0xA0b86a33E6441e6e80fd4e3Cd9Cc5F7b",
		"buyAmount":            "1002345000000000000",
		"sellAmount":           "1000000",
		"sources": []map[string]interface{}{
			{
				"name":       "Uniswap_V3",
				"proportion": "0.8",
			},
			{
				"name":       "SushiSwap",
				"proportion": "0.2",
			},
		},
		"orders":               []interface{}{},
		"allowanceTarget":      "0x1234567890123456789012345678901234567890",
		"decodedUniqueId":      "quote-123456789",
		"sellTokenToEthRate":   "1.0",
		"buyTokenToEthRate":    "1.0",
		"expectedSlippage":     "0.01",
	}

	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/ethereum/swap/v1/quote", r.URL.Path)
		
		// Check query parameters
		query := r.URL.Query()
		assert.Equal(t, "0xA0b86a33E6441e6e80fd4e3Cd9Cc5F7b", query.Get("sellToken"))
		assert.Equal(t, "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2", query.Get("buyToken"))
		assert.Equal(t, "1000000", query.Get("sellAmount"))
		assert.Equal(t, "0x1234567890123456789012345678901234567890", query.Get("takerAddress"))
		assert.Equal(t, "0.5000", query.Get("slippagePercentage"))

		// Check headers
		assert.Equal(t, "test-api-key", r.Header.Get("0x-api-key"))
		assert.Equal(t, "application/json", r.Header.Get("Accept"))

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

	client := swap.NewZeroXClient(config)

	// Test request
	req := clients.QuoteRequest{
		FromChainID: "1",
		FromToken:   "0xA0b86a33E6441e6e80fd4e3Cd9Cc5F7b",
		ToToken:     "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		Amount:      "1000000",
		UserAddress: "0x1234567890123456789012345678901234567890",
		Slippage:    0.5,
	}

	ctx := context.Background()
	quote, err := client.GetQuote(ctx, req)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, quote)

	assert.Equal(t, "quote-123456789", quote.ID)
	assert.Equal(t, "swap", quote.Type)
	assert.Equal(t, "0x", quote.Provider)
	assert.Equal(t, "1", quote.FromChainID)
	assert.Equal(t, "1", quote.ToChainID) // Same chain for swaps
	assert.Equal(t, "1000000", quote.FromAmount)
	assert.Equal(t, "1002345000000000000", quote.ToAmount)
	assert.Equal(t, "1.002345", quote.ExchangeRate)
	assert.Equal(t, 0.0234, quote.PriceImpact)
	assert.Equal(t, "115000", quote.EstimatedGas)
	assert.Equal(t, "20000000000", quote.GasPriceWei)
	
	// Check transaction data
	require.NotNil(t, quote.TransactionData)
	assert.Equal(t, "0x1234567890123456789012345678901234567890", quote.TransactionData.To)
	assert.Equal(t, "0xabcdef1234567890", quote.TransactionData.Data)
	assert.Equal(t, "0", quote.TransactionData.Value)
	assert.Equal(t, "120000", quote.TransactionData.GasLimit)
	assert.Equal(t, "1", quote.TransactionData.ChainID)

	// Check route steps
	require.Len(t, quote.Route, 2)
	assert.Equal(t, "Uniswap_V3", quote.Route[0].Protocol)
	assert.Equal(t, "swap", quote.Route[0].Type)
	assert.Equal(t, "0.8", quote.Route[0].Percentage)
	
	assert.Equal(t, "SushiSwap", quote.Route[1].Protocol)
	assert.Equal(t, "swap", quote.Route[1].Type)
	assert.Equal(t, "0.2", quote.Route[1].Percentage)

	// Check protocol fee
	require.Len(t, quote.Fees, 1)
	assert.Equal(t, "protocol", quote.Fees[0].Type)
	assert.Equal(t, "1000", quote.Fees[0].Amount)
	assert.Equal(t, "0x Protocol Fee", quote.Fees[0].Description)
}

func TestZeroXClient_GetSupportedTokens(t *testing.T) {
	// Mock tokens response
	mockResponse := map[string]interface{}{
		"records": []map[string]interface{}{
			{
				"address":  "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
				"chainId":  1,
				"name":     "Wrapped Ether",
				"symbol":   "WETH",
				"decimals": 18,
				"logoURI":  "https://example.com/weth.png",
				"tags":     []string{"native"},
			},
			{
				"address":  "0xA0b86a33E6441e6e80fd4e3Cd9Cc5F7b",
				"chainId":  1,
				"name":     "USD Coin",
				"symbol":   "USDC",
				"decimals": 6,
				"logoURI":  "https://example.com/usdc.png",
				"tags":     []string{"stablecoin"},
			},
		},
	}

	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/ethereum/swap/v1/tokens", r.URL.Path)

		// Check headers
		assert.Equal(t, "test-api-key", r.Header.Get("0x-api-key"))
		assert.Equal(t, "application/json", r.Header.Get("Accept"))

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

	client := swap.NewZeroXClient(config)

	// Test
	ctx := context.Background()
	tokens, err := client.GetSupportedTokens(ctx, "1")

	// Assertions
	require.NoError(t, err)
	require.Len(t, tokens, 2)

	// Check WETH
	weth := tokens[0]
	assert.Equal(t, "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2", weth.Address)
	assert.Equal(t, "WETH", weth.Symbol)
	assert.Equal(t, "Wrapped Ether", weth.Name)
	assert.Equal(t, 18, weth.Decimals)
	assert.Equal(t, "1", weth.ChainID)
	assert.Equal(t, "https://example.com/weth.png", weth.LogoURI)

	// Check USDC
	usdc := tokens[1]
	assert.Equal(t, "0xA0b86a33E6441e6e80fd4e3Cd9Cc5F7b", usdc.Address)
	assert.Equal(t, "USDC", usdc.Symbol)
	assert.Equal(t, "USD Coin", usdc.Name)
	assert.Equal(t, 6, usdc.Decimals)
	assert.Equal(t, "1", usdc.ChainID)
	assert.Equal(t, "https://example.com/usdc.png", usdc.LogoURI)
}

func TestZeroXClient_IsHealthy(t *testing.T) {
	// Test healthy response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/ethereum/swap/v1/tokens", r.URL.Path)
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

	client := swap.NewZeroXClient(config)

	ctx := context.Background()
	healthy := client.IsHealthy(ctx)

	assert.True(t, healthy)
}

func TestZeroXClient_GetProviderName(t *testing.T) {
	config := clients.ClientConfig{
		BaseURL: "https://api.test.com",
		APIKey:  "test-key",
	}

	client := swap.NewZeroXClient(config)
	assert.Equal(t, "0x", client.GetProviderName())
}

func TestZeroXClient_UnsupportedChain(t *testing.T) {
	config := clients.ClientConfig{
		BaseURL:    "https://api.test.com",
		APIKey:     "test-api-key",
		Timeout:    5 * time.Second,
		MaxRetries: 1,
		RetryDelay: 100 * time.Millisecond,
		RateLimit: clients.RateLimitConfig{
			RequestsPerSecond: 10,
			BurstSize:         20,
		},
	}

	client := swap.NewZeroXClient(config)

	// Test unsupported chain
	req := clients.QuoteRequest{
		FromChainID: "999", // Unsupported chain
		FromToken:   "0xA0b86a33E6441e6e80fd4e3Cd9Cc5F7b",
		ToToken:     "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		Amount:      "1000000",
		UserAddress: "0x1234567890123456789012345678901234567890",
	}

	ctx := context.Background()
	quote, err := client.GetQuote(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, quote)
	assert.Contains(t, err.Error(), "unsupported chain ID")
}