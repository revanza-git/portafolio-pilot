package swap

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/defi-dashboard/backend/internal/clients"
)

// ZeroXClient implements SwapClient for 0x API
type ZeroXClient struct {
	httpClient clients.HTTPClient
	baseURL    string
	apiKey     string
}

// 0x API types
type zeroXQuoteRequest struct {
	SellToken            string  `json:"sellToken"`
	BuyToken             string  `json:"buyToken"`
	SellAmount           string  `json:"sellAmount,omitempty"`
	BuyAmount            string  `json:"buyAmount,omitempty"`
	TakerAddress         string  `json:"takerAddress,omitempty"`
	SlippagePercentage   float64 `json:"slippagePercentage,omitempty"`
	GasPrice             string  `json:"gasPrice,omitempty"`
	SkipValidation       bool    `json:"skipValidation,omitempty"`
	FeeRecipient         string  `json:"feeRecipient,omitempty"`
	BuyTokenPercentageFee float64 `json:"buyTokenPercentageFee,omitempty"`
	Affiliate            string  `json:"affiliate,omitempty"`
}

type zeroXQuoteResponse struct {
	ChainId              int                    `json:"chainId"`
	Price                string                 `json:"price"`
	GuaranteedPrice      string                 `json:"guaranteedPrice"`
	EstimatedPriceImpact string                 `json:"estimatedPriceImpact"`
	To                   string                 `json:"to"`
	Data                 string                 `json:"data"`
	Value                string                 `json:"value"`
	Gas                  string                 `json:"gas"`
	EstimatedGas         string                 `json:"estimatedGas"`
	GasPrice             string                 `json:"gasPrice"`
	ProtocolFee          string                 `json:"protocolFee"`
	MinimumProtocolFee   string                 `json:"minimumProtocolFee"`
	BuyTokenAddress      string                 `json:"buyTokenAddress"`
	SellTokenAddress     string                 `json:"sellTokenAddress"`
	BuyAmount            string                 `json:"buyAmount"`
	SellAmount           string                 `json:"sellAmount"`
	Sources              []zeroXSource          `json:"sources"`
	Orders               []zeroXOrder           `json:"orders"`
	AllowanceTarget      string                 `json:"allowanceTarget"`
	DecodedUniqueId      string                 `json:"decodedUniqueId"`
	SellTokenToEthRate   string                 `json:"sellTokenToEthRate"`
	BuyTokenToEthRate    string                 `json:"buyTokenToEthRate"`
	ExpectedSlippage     string                 `json:"expectedSlippage"`
}

type zeroXSource struct {
	Name       string `json:"name"`
	Proportion string `json:"proportion"`
}

type zeroXOrder struct {
	MakerToken       string `json:"makerToken"`
	TakerToken       string `json:"takerToken"`
	MakerAmount      string `json:"makerAmount"`
	TakerAmount      string `json:"takerAmount"`
	FillData         struct {
		TokenAddressPath []string `json:"tokenAddressPath"`
		Router           string   `json:"router"`
	} `json:"fillData"`
	Source          string `json:"source"`
	SourcePathId    string `json:"sourcePathId"`
	Type            int    `json:"type"`
}

type zeroXTokensResponse struct {
	Records []zeroXToken `json:"records"`
}

type zeroXToken struct {
	Address     string   `json:"address"`
	ChainId     int      `json:"chainId"`
	Name        string   `json:"name"`
	Symbol      string   `json:"symbol"`
	Decimals    int      `json:"decimals"`
	LogoURI     string   `json:"logoURI"`
	Tags        []string `json:"tags"`
}

type zeroXErrorResponse struct {
	Code           int    `json:"code"`
	Reason         string `json:"reason"`
	ValidationErrors []struct {
		Field  string `json:"field"`
		Code   int    `json:"code"`
		Reason string `json:"reason"`
	} `json:"validationErrors"`
}

// NewZeroXClient creates a new 0x swap client
func NewZeroXClient(config clients.ClientConfig) *ZeroXClient {
	httpClient := clients.NewBaseHTTPClient(config)
	
	return &ZeroXClient{
		httpClient: httpClient,
		baseURL:    config.BaseURL,
		apiKey:     config.APIKey,
	}
}

// GetQuote fetches a swap quote from 0x
func (c *ZeroXClient) GetQuote(ctx context.Context, req clients.QuoteRequest) (*clients.Quote, error) {
	// Determine chain-specific endpoint
	chainId, err := strconv.Atoi(req.FromChainID)
	if err != nil {
		return nil, fmt.Errorf("invalid chain ID: %w", err)
	}

	chainName := c.getChainName(chainId)
	if chainName == "" {
		return nil, fmt.Errorf("unsupported chain ID: %d", chainId)
	}

	url := fmt.Sprintf("%s/%s/swap/v1/quote", c.baseURL, chainName)
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add query parameters
	q := httpReq.URL.Query()
	q.Add("sellToken", req.FromToken)
	q.Add("buyToken", req.ToToken)
	q.Add("sellAmount", req.Amount)
	if req.UserAddress != "" {
		q.Add("takerAddress", req.UserAddress)
	}
	if req.Slippage > 0 {
		q.Add("slippagePercentage", fmt.Sprintf("%.4f", req.Slippage))
	}
	q.Add("skipValidation", "true")
	httpReq.URL.RawQuery = q.Encode()

	// Add headers
	if c.apiKey != "" {
		httpReq.Header.Set("0x-api-key", c.apiKey)
	}
	httpReq.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}

	var zeroXResp zeroXQuoteResponse
	if err := clients.ParseResponse(resp, &zeroXResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Convert to unified format
	return c.convertToUnifiedQuote(zeroXResp, req), nil
}

// GetSupportedTokens returns supported tokens for a specific chain
func (c *ZeroXClient) GetSupportedTokens(ctx context.Context, chainID string) ([]clients.Token, error) {
	chainId, err := strconv.Atoi(chainID)
	if err != nil {
		return nil, fmt.Errorf("invalid chain ID: %w", err)
	}

	chainName := c.getChainName(chainId)
	if chainName == "" {
		return nil, fmt.Errorf("unsupported chain ID: %d", chainId)
	}

	url := fmt.Sprintf("%s/%s/swap/v1/tokens", c.baseURL, chainName)
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.apiKey != "" {
		httpReq.Header.Set("0x-api-key", c.apiKey)
	}
	httpReq.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}

	var zeroXResp zeroXTokensResponse
	if err := clients.ParseResponse(resp, &zeroXResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	tokens := make([]clients.Token, len(zeroXResp.Records))
	for i, token := range zeroXResp.Records {
		tokens[i] = clients.Token{
			Address:  token.Address,
			Symbol:   token.Symbol,
			Name:     token.Name,
			Decimals: token.Decimals,
			ChainID:  strconv.Itoa(token.ChainId),
			LogoURI:  token.LogoURI,
		}
	}

	return tokens, nil
}

// GetProviderName returns the name of the swap provider
func (c *ZeroXClient) GetProviderName() string {
	return "0x"
}

// IsHealthy checks if the provider API is responding
func (c *ZeroXClient) IsHealthy(ctx context.Context) bool {
	// Use Ethereum mainnet for health check
	url := fmt.Sprintf("%s/ethereum/swap/v1/tokens", c.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return false
	}

	if c.apiKey != "" {
		httpReq.Header.Set("0x-api-key", c.apiKey)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// getChainName returns the 0x API chain name for a given chain ID
func (c *ZeroXClient) getChainName(chainID int) string {
	switch chainID {
	case 1:
		return "ethereum"
	case 10:
		return "optimism"
	case 56:
		return "bsc"
	case 137:
		return "polygon"
	case 250:
		return "fantom"
	case 42161:
		return "arbitrum"
	case 43114:
		return "avalanche"
	case 8453:
		return "base"
	default:
		return ""
	}
}

// convertToUnifiedQuote converts 0x response to unified quote format
func (c *ZeroXClient) convertToUnifiedQuote(zeroXResp zeroXQuoteResponse, req clients.QuoteRequest) *clients.Quote {
	// Parse price impact
	priceImpact := 0.0
	if zeroXResp.EstimatedPriceImpact != "" {
		if pi, err := strconv.ParseFloat(zeroXResp.EstimatedPriceImpact, 64); err == nil {
			priceImpact = pi
		}
	}

	quote := &clients.Quote{
		ID:            zeroXResp.DecodedUniqueId,
		Type:          "swap",
		Provider:      "0x",
		FromChainID:   req.FromChainID,
		ToChainID:     req.FromChainID, // Same chain for swaps
		FromAmount:    zeroXResp.SellAmount,
		ToAmount:      zeroXResp.BuyAmount,
		ExchangeRate:  zeroXResp.Price,
		PriceImpact:   priceImpact,
		EstimatedGas:  zeroXResp.EstimatedGas,
		GasPriceWei:   zeroXResp.GasPrice,
		ExpiresAt:     time.Now().Add(30 * time.Second), // 0x quotes expire quickly
	}

	// Add transaction data
	quote.TransactionData = &clients.TransactionData{
		To:       zeroXResp.To,
		Data:     zeroXResp.Data,
		Value:    zeroXResp.Value,
		GasLimit: zeroXResp.Gas,
		ChainID:  req.FromChainID,
	}

	// Convert sources to route steps
	var routeSteps []clients.RouteStep
	for _, source := range zeroXResp.Sources {
		if strings.TrimSpace(source.Proportion) == "0" {
			continue // Skip sources with 0% proportion
		}
		
		routeSteps = append(routeSteps, clients.RouteStep{
			Protocol:   source.Name,
			Type:       "swap",
			FromAmount: zeroXResp.SellAmount,
			ToAmount:   zeroXResp.BuyAmount,
			Percentage: source.Proportion,
		})
	}
	quote.Route = routeSteps

	// Add protocol fee
	if zeroXResp.ProtocolFee != "" && zeroXResp.ProtocolFee != "0" {
		quote.Fees = append(quote.Fees, clients.Fee{
			Type:        "protocol",
			Amount:      zeroXResp.ProtocolFee,
			Description: "0x Protocol Fee",
		})
	}

	// We don't have detailed token info in quote response, so we'll use basic info
	quote.FromToken = clients.Token{
		Address: zeroXResp.SellTokenAddress,
		ChainID: req.FromChainID,
	}
	quote.ToToken = clients.Token{
		Address: zeroXResp.BuyTokenAddress,
		ChainID: req.FromChainID,
	}

	return quote
}