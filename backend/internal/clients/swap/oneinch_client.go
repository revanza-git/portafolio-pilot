package swap

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/defi-dashboard/backend/internal/clients"
)

// OneInchClient implements SwapClient for 1inch API
type OneInchClient struct {
	httpClient clients.HTTPClient
	baseURL    string
	apiKey     string
}

// 1inch API types
type oneInchQuoteRequest struct {
	FromTokenAddress string  `json:"fromTokenAddress"`
	ToTokenAddress   string  `json:"toTokenAddress"`
	Amount           string  `json:"amount"`
	FromAddress      string  `json:"fromAddress,omitempty"`
	Slippage         float64 `json:"slippage,omitempty"`
	Protocols        string  `json:"protocols,omitempty"`
	Fee              float64 `json:"fee,omitempty"`
	GasLimit         int     `json:"gasLimit,omitempty"`
	ConnectorTokens  string  `json:"connectorTokens,omitempty"`
	ComplexityLevel  int     `json:"complexityLevel,omitempty"`
	MainRouteParts   int     `json:"mainRouteParts,omitempty"`
	Parts            int     `json:"parts,omitempty"`
	GasPrice         string  `json:"gasPrice,omitempty"`
}

type oneInchQuoteResponse struct {
	FromToken         oneInchToken    `json:"fromToken"`
	ToToken           oneInchToken    `json:"toToken"`
	ToTokenAmount     string          `json:"toTokenAmount"`
	FromTokenAmount   string          `json:"fromTokenAmount"`
	Protocols         [][]oneInchProtocolEntry `json:"protocols"`
	EstimatedGas      int             `json:"estimatedGas"`
	Tx                oneInchTx       `json:"tx,omitempty"`
}

type oneInchToken struct {
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Decimals int    `json:"decimals"`
	LogoURI  string `json:"logoURI"`
}

type oneInchProtocolEntry struct {
	Name         string                `json:"name"`
	Part         float64               `json:"part"`
	FromTokenAddress string            `json:"fromTokenAddress"`
	ToTokenAddress   string            `json:"toTokenAddress"`
}

type oneInchTx struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Data     string `json:"data"`
	Value    string `json:"value"`
	GasPrice string `json:"gasPrice"`
	Gas      int    `json:"gas"`
}

type oneInchTokensResponse struct {
	Tokens map[string]oneInchTokenInfo `json:"tokens"`
}

type oneInchTokenInfo struct {
	Symbol      string   `json:"symbol"`
	Name        string   `json:"name"`
	Address     string   `json:"address"`
	Decimals    int      `json:"decimals"`
	LogoURI     string   `json:"logoURI"`
	Tags        []string `json:"tags"`
	Providers   []string `json:"providers"`
	EIP2612     bool     `json:"eip2612"`
	IsFoT       bool     `json:"isFoT"`
}

type oneInchProtocolsResponse struct {
	Protocols []oneInchProtocolInfo `json:"protocols"`
}

type oneInchProtocolInfo struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Img   string `json:"img"`
}

type oneInchHealthResponse struct {
	Status string `json:"status"`
}

// NewOneInchClient creates a new 1inch swap client
func NewOneInchClient(config clients.ClientConfig) *OneInchClient {
	httpClient := clients.NewBaseHTTPClient(config)
	
	return &OneInchClient{
		httpClient: httpClient,
		baseURL:    config.BaseURL,
		apiKey:     config.APIKey,
	}
}

// GetQuote fetches a swap quote from 1inch
func (c *OneInchClient) GetQuote(ctx context.Context, req clients.QuoteRequest) (*clients.Quote, error) {
	chainId, err := strconv.Atoi(req.FromChainID)
	if err != nil {
		return nil, fmt.Errorf("invalid chain ID: %w", err)
	}

	if !c.isChainSupported(chainId) {
		return nil, fmt.Errorf("unsupported chain ID: %d", chainId)
	}

	url := fmt.Sprintf("%s/v5.0/%d/quote", c.baseURL, chainId)
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add query parameters
	q := httpReq.URL.Query()
	q.Add("fromTokenAddress", req.FromToken)
	q.Add("toTokenAddress", req.ToToken)
	q.Add("amount", req.Amount)
	if req.UserAddress != "" {
		q.Add("fromAddress", req.UserAddress)
	}
	if req.Slippage > 0 {
		q.Add("slippage", fmt.Sprintf("%.1f", req.Slippage))
	}
	httpReq.URL.RawQuery = q.Encode()

	// Add headers
	if c.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	httpReq.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}

	var oneInchResp oneInchQuoteResponse
	if err := clients.ParseResponse(resp, &oneInchResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Convert to unified format
	return c.convertToUnifiedQuote(oneInchResp, req), nil
}

// GetSupportedTokens returns supported tokens for a specific chain
func (c *OneInchClient) GetSupportedTokens(ctx context.Context, chainID string) ([]clients.Token, error) {
	chainId, err := strconv.Atoi(chainID)
	if err != nil {
		return nil, fmt.Errorf("invalid chain ID: %w", err)
	}

	if !c.isChainSupported(chainId) {
		return nil, fmt.Errorf("unsupported chain ID: %d", chainId)
	}

	url := fmt.Sprintf("%s/v5.0/%d/tokens", c.baseURL, chainId)
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	httpReq.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}

	var oneInchResp oneInchTokensResponse
	if err := clients.ParseResponse(resp, &oneInchResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	var tokens []clients.Token
	for address, token := range oneInchResp.Tokens {
		tokens = append(tokens, clients.Token{
			Address:  address,
			Symbol:   token.Symbol,
			Name:     token.Name,
			Decimals: token.Decimals,
			ChainID:  chainID,
			LogoURI:  token.LogoURI,
		})
	}

	return tokens, nil
}

// GetProviderName returns the name of the swap provider
func (c *OneInchClient) GetProviderName() string {
	return "1inch"
}

// IsHealthy checks if the provider API is responding
func (c *OneInchClient) IsHealthy(ctx context.Context) bool {
	// Use Ethereum mainnet for health check
	url := fmt.Sprintf("%s/v5.0/1/healthcheck", c.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return false
	}

	if c.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false
	}

	var healthResp oneInchHealthResponse
	if err := clients.ParseResponse(resp, &healthResp); err != nil {
		return false
	}

	return healthResp.Status == "OK"
}

// isChainSupported checks if a chain ID is supported by 1inch
func (c *OneInchClient) isChainSupported(chainID int) bool {
	supportedChains := []int{
		1,     // Ethereum
		10,    // Optimism
		56,    // BSC
		100,   // Gnosis
		137,   // Polygon
		250,   // Fantom
		8217,  // Klaytn
		42161, // Arbitrum
		43114, // Avalanche
		1313161554, // Aurora
	}

	for _, supported := range supportedChains {
		if chainID == supported {
			return true
		}
	}
	return false
}

// convertToUnifiedQuote converts 1inch response to unified quote format
func (c *OneInchClient) convertToUnifiedQuote(oneInchResp oneInchQuoteResponse, req clients.QuoteRequest) *clients.Quote {
	quote := &clients.Quote{
		ID:          fmt.Sprintf("1inch-%d", time.Now().Unix()),
		Type:        "swap",
		Provider:    "1inch",
		FromChainID: req.FromChainID,
		ToChainID:   req.FromChainID, // Same chain for swaps
		FromToken: clients.Token{
			Address:  oneInchResp.FromToken.Address,
			Symbol:   oneInchResp.FromToken.Symbol,
			Name:     oneInchResp.FromToken.Name,
			Decimals: oneInchResp.FromToken.Decimals,
			ChainID:  req.FromChainID,
			LogoURI:  oneInchResp.FromToken.LogoURI,
		},
		ToToken: clients.Token{
			Address:  oneInchResp.ToToken.Address,
			Symbol:   oneInchResp.ToToken.Symbol,
			Name:     oneInchResp.ToToken.Name,
			Decimals: oneInchResp.ToToken.Decimals,
			ChainID:  req.FromChainID,
			LogoURI:  oneInchResp.ToToken.LogoURI,
		},
		FromAmount:   oneInchResp.FromTokenAmount,
		ToAmount:     oneInchResp.ToTokenAmount,
		EstimatedGas: strconv.Itoa(oneInchResp.EstimatedGas),
		ExpiresAt:    time.Now().Add(60 * time.Second), // 1inch quotes expire after 60 seconds
	}

	// Calculate exchange rate
	if oneInchResp.FromTokenAmount != "0" && oneInchResp.ToTokenAmount != "0" {
		fromAmt, _ := strconv.ParseFloat(oneInchResp.FromTokenAmount, 64)
		toAmt, _ := strconv.ParseFloat(oneInchResp.ToTokenAmount, 64)
		if fromAmt > 0 {
			rate := toAmt / fromAmt
			quote.ExchangeRate = fmt.Sprintf("%.18f", rate)
		}
	}

	// Add transaction data if present
	if oneInchResp.Tx.To != "" {
		quote.TransactionData = &clients.TransactionData{
			To:       oneInchResp.Tx.To,
			Data:     oneInchResp.Tx.Data,
			Value:    oneInchResp.Tx.Value,
			GasLimit: strconv.Itoa(oneInchResp.Tx.Gas),
			ChainID:  req.FromChainID,
		}
		quote.GasPriceWei = oneInchResp.Tx.GasPrice
	}

	// Convert protocols to route steps
	var routeSteps []clients.RouteStep
	for _, protocolGroup := range oneInchResp.Protocols {
		for _, protocol := range protocolGroup {
			if protocol.Part > 0 {
				routeSteps = append(routeSteps, clients.RouteStep{
					Protocol:     protocol.Name,
					Type:         "swap",
					FromAmount:   oneInchResp.FromTokenAmount,
					ToAmount:     oneInchResp.ToTokenAmount,
					Percentage:   fmt.Sprintf("%.2f", protocol.Part*100),
					FromToken: clients.Token{
						Address: protocol.FromTokenAddress,
						ChainID: req.FromChainID,
					},
					ToToken: clients.Token{
						Address: protocol.ToTokenAddress,
						ChainID: req.FromChainID,
					},
				})
			}
		}
	}
	quote.Route = routeSteps

	return quote
}