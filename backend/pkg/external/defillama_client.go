package external

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	DefiLlamaAPIBase = "https://api.llama.fi"
	DefiLlamaRateLimit = 300 // 300 requests per minute
)

type DefiLlamaClient struct {
	httpClient  *http.Client
	rateLimiter *RateLimiter
}

func NewDefiLlamaClient() *DefiLlamaClient {
	return &DefiLlamaClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		rateLimiter: NewRateLimiter(DefiLlamaRateLimit, time.Minute),
	}
}

type YieldPool struct {
	Pool        string  `json:"pool"`
	Project     string  `json:"project"`
	Symbol      string  `json:"symbol"`
	Chain       string  `json:"chain"`
	TVL         float64 `json:"tvlUsd"`
	APY         float64 `json:"apy"`
	APYBase     float64 `json:"apyBase"`
	APYReward   float64 `json:"apyReward"`
	IL7d        float64 `json:"il7d"`
	Exposure    string  `json:"exposure"`
	StableCoin  bool    `json:"stablecoin"`
}

// GetYieldPools fetches all yield pools
func (c *DefiLlamaClient) GetYieldPools(ctx context.Context) ([]YieldPool, error) {
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/pools", DefiLlamaAPIBase)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("DefiLlama API error: %d", resp.StatusCode)
	}

	var response struct {
		Status string      `json:"status"`
		Data   []YieldPool `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

// GetPoolsByChain fetches yield pools for a specific chain
func (c *DefiLlamaClient) GetPoolsByChain(ctx context.Context, chain string) ([]YieldPool, error) {
	pools, err := c.GetYieldPools(ctx)
	if err != nil {
		return nil, err
	}

	// Filter by chain
	var chainPools []YieldPool
	for _, pool := range pools {
		if pool.Chain == chain {
			chainPools = append(chainPools, pool)
		}
	}

	return chainPools, nil
}

type ProtocolTVL struct {
	Name string `json:"name"`
	TVL  float64 `json:"tvl"`
}

// GetProtocolTVL fetches TVL for a specific protocol
func (c *DefiLlamaClient) GetProtocolTVL(ctx context.Context, protocol string) (*ProtocolTVL, error) {
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/protocol/%s", DefiLlamaAPIBase, protocol)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("DefiLlama API error: %d", resp.StatusCode)
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	tvl, ok := data["tvl"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid TVL data format")
	}

	currentTVL, ok := tvl["current"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid current TVL format")
	}

	return &ProtocolTVL{
		Name: protocol,
		TVL:  currentTVL,
	}, nil
}

// Chain name mappings
var ChainMappings = map[string]string{
	"ethereum": "Ethereum",
	"eth":      "Ethereum",
	"polygon":  "Polygon",
	"matic":    "Polygon",
	"arbitrum": "Arbitrum",
	"optimism": "Optimism",
	"base":     "Base",
	"bsc":      "BSC",
	"binance":  "BSC",
}