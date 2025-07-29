package external

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	CoinGeckoAPIBase = "https://api.coingecko.com/api/v3"
	RateLimitPerMin  = 50 // Free tier limit
)

type CoinGeckoClient struct {
	httpClient *http.Client
	apiKey     string
	rateLimiter *RateLimiter
}

func NewCoinGeckoClient(apiKey string) *CoinGeckoClient {
	return &CoinGeckoClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		apiKey:      apiKey,
		rateLimiter: NewRateLimiter(RateLimitPerMin, time.Minute),
	}
}

type TokenPrice struct {
	USD         float64 `json:"usd"`
	USD24hChange float64 `json:"usd_24h_change"`
}

type PriceResponse map[string]TokenPrice

// GetTokenPrices fetches current prices for multiple tokens
func (c *CoinGeckoClient) GetTokenPrices(ctx context.Context, tokenIDs []string) (PriceResponse, error) {
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, err
	}

	ids := strings.Join(tokenIDs, ",")
	url := fmt.Sprintf("%s/simple/price?ids=%s&vs_currencies=usd&include_24hr_change=true", CoinGeckoAPIBase, ids)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	if c.apiKey != "" {
		req.Header.Set("x-cg-pro-api-key", c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("CoinGecko API error: %d", resp.StatusCode)
	}

	var prices PriceResponse
	if err := json.NewDecoder(resp.Body).Decode(&prices); err != nil {
		return nil, err
	}

	return prices, nil
}

// GetPriceHistory fetches historical price data
func (c *CoinGeckoClient) GetPriceHistory(ctx context.Context, tokenID string, days int) ([][]float64, error) {
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/coins/%s/market_chart?vs_currency=usd&days=%d", CoinGeckoAPIBase, tokenID, days)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	if c.apiKey != "" {
		req.Header.Set("x-cg-pro-api-key", c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("CoinGecko API error: %d", resp.StatusCode)
	}

	var data struct {
		Prices [][]float64 `json:"prices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data.Prices, nil
}

// Token ID mappings
var TokenIDMappings = map[string]string{
	"eth":  "ethereum",
	"weth": "weth",
	"usdc": "usd-coin",
	"usdt": "tether",
	"dai":  "dai",
	"wbtc": "wrapped-bitcoin",
	"uni":  "uniswap",
	"aave": "aave",
	"link": "chainlink",
	"matic": "matic-network",
}