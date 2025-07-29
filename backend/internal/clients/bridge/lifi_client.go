package bridge

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/defi-dashboard/backend/internal/clients"
)

// LiFiClient implements BridgeClient for LI.FI API
type LiFiClient struct {
	httpClient clients.HTTPClient
	baseURL    string
	apiKey     string
}

// LiFi API types
type lifiQuoteRequest struct {
	FromChain   string `json:"fromChain"`
	ToChain     string `json:"toChain"`
	FromToken   string `json:"fromToken"`
	ToToken     string `json:"toToken"`
	FromAmount  string `json:"fromAmount"`
	FromAddress string `json:"fromAddress"`
	ToAddress   string `json:"toAddress"`
	Options     struct {
		Slippage        float64  `json:"slippage,omitempty"`
		AllowBridges    []string `json:"allowBridges,omitempty"`
		DenyBridges     []string `json:"denyBridges,omitempty"`
		PreferBridges   []string `json:"preferBridges,omitempty"`
		AllowExchanges  []string `json:"allowExchanges,omitempty"`
		DenyExchanges   []string `json:"denyExchanges,omitempty"`
		PreferExchanges []string `json:"preferExchanges,omitempty"`
	} `json:"options,omitempty"`
}

type lifiQuoteResponse struct {
	Routes []lifiRoute `json:"routes"`
}

type lifiRoute struct {
	ID           string       `json:"id"`
	FromChainID  int          `json:"fromChainId"`
	ToChainID    int          `json:"toChainId"`
	FromToken    lifiToken    `json:"fromToken"`
	ToToken      lifiToken    `json:"toToken"`
	FromAmount   string       `json:"fromAmount"`
	ToAmount     string       `json:"toAmount"`
	Steps        []lifiStep   `json:"steps"`
	Tags         []string     `json:"tags"`
	GasCostUSD   string       `json:"gasCostUSD"`
	Insurance    *lifiInsurance `json:"insurance,omitempty"`
}

type lifiStep struct {
	ID                string          `json:"id"`
	Type              string          `json:"type"`
	Tool              string          `json:"tool"`
	Action            lifiAction      `json:"action"`
	Estimate          lifiEstimate    `json:"estimate"`
	IncludedSteps     []lifiStep      `json:"includedSteps,omitempty"`
	TransactionRequest *lifiTxRequest `json:"transactionRequest,omitempty"`
}

type lifiAction struct {
	FromChainID   int       `json:"fromChainId"`
	ToChainID     int       `json:"toChainId"`
	FromToken     lifiToken `json:"fromToken"`
	ToToken       lifiToken `json:"toToken"`
	FromAmount    string    `json:"fromAmount"`
	ToAmount      string    `json:"toAmount"`
	Slippage      float64   `json:"slippage"`
	FromAddress   string    `json:"fromAddress"`
	ToAddress     string    `json:"toAddress"`
}

type lifiEstimate struct {
	FromAmount         string         `json:"fromAmount"`
	ToAmount           string         `json:"toAmount"`
	ExecutionDuration  int            `json:"executionDuration"`
	FeeCosts           []lifiFee      `json:"feeCosts"`
	GasCosts           []lifiGasCost  `json:"gasCosts"`
	ToAmountMin        string         `json:"toAmountMin"`
	DataGasLimit       string         `json:"dataGasLimit,omitempty"`
	ApprovalAddress    string         `json:"approvalAddress,omitempty"`
}

type lifiFee struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Token       lifiToken `json:"token"`
	Amount      string    `json:"amount"`
	AmountUSD   string    `json:"amountUSD"`
	Percentage  string    `json:"percentage"`
	Included    bool      `json:"included"`
}

type lifiGasCost struct {
	Type     string    `json:"type"`
	Price    string    `json:"price"`
	Estimate string    `json:"estimate"`
	Limit    string    `json:"limit"`
	Amount   string    `json:"amount"`
	AmountUSD string   `json:"amountUSD"`
	Token    lifiToken `json:"token"`
}

type lifiToken struct {
	Address  string `json:"address"`
	ChainID  int    `json:"chainId"`
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Decimals int    `json:"decimals"`
	LogoURI  string `json:"logoURI"`
	PriceUSD string `json:"priceUSD"`
}

type lifiTxRequest struct {
	Data     string `json:"data"`
	To       string `json:"to"`
	Value    string `json:"value"`
	From     string `json:"from"`
	ChainID  int    `json:"chainId"`
	GasLimit string `json:"gasLimit"`
	GasPrice string `json:"gasPrice"`
}

type lifiInsurance struct {
	State string `json:"state"`
	FeeAmountUsd string `json:"feeAmountUsd"`
}

type lifiChainsResponse struct {
	Chains []lifiChainInfo `json:"chains"`
}

type lifiChainInfo struct {
	ID             int       `json:"id"`
	Key            string    `json:"key"`
	Name           string    `json:"name"`
	Coin           string    `json:"coin"`
	MainnetId      int       `json:"mainnetId"`
	LogoURI        string    `json:"logoURI"`
	TokenlistUrl   string    `json:"tokenlistUrl"`
	FaucetUrls     []string  `json:"faucetUrls"`
	NativeToken    lifiToken `json:"nativeToken"`
}

type lifiTokensResponse struct {
	Tokens map[string][]lifiToken `json:"tokens"`
}

// NewLiFiClient creates a new LI.FI bridge client
func NewLiFiClient(config clients.ClientConfig) *LiFiClient {
	httpClient := clients.NewBaseHTTPClient(config)
	
	return &LiFiClient{
		httpClient: httpClient,
		baseURL:    config.BaseURL,
		apiKey:     config.APIKey,
	}
}

// GetQuote fetches a bridge quote from LI.FI
func (c *LiFiClient) GetQuote(ctx context.Context, req clients.QuoteRequest) (*clients.Quote, error) {
	// Convert to LI.FI request format
	lifiReq := lifiQuoteRequest{
		FromChain:   req.FromChainID,
		ToChain:     req.ToChainID,
		FromToken:   req.FromToken,
		ToToken:     req.ToToken,
		FromAmount:  req.Amount,
		FromAddress: req.UserAddress,
		ToAddress:   req.UserAddress,
	}

	if req.Slippage > 0 {
		lifiReq.Options.Slippage = req.Slippage
	}

	// Make API request
	url := fmt.Sprintf("%s/quote", c.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add query parameters
	q := httpReq.URL.Query()
	q.Add("fromChain", lifiReq.FromChain)
	q.Add("toChain", lifiReq.ToChain)
	q.Add("fromToken", lifiReq.FromToken)
	q.Add("toToken", lifiReq.ToToken)
	q.Add("fromAmount", lifiReq.FromAmount)
	q.Add("fromAddress", lifiReq.FromAddress)
	q.Add("toAddress", lifiReq.ToAddress)
	if lifiReq.Options.Slippage > 0 {
		q.Add("options.slippage", fmt.Sprintf("%.4f", lifiReq.Options.Slippage))
	}
	httpReq.URL.RawQuery = q.Encode()

	// Add headers
	if c.apiKey != "" {
		httpReq.Header.Set("x-lifi-api-key", c.apiKey)
	}
	httpReq.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}

	var lifiResp lifiQuoteResponse
	if err := clients.ParseResponse(resp, &lifiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Convert to unified format
	if len(lifiResp.Routes) == 0 {
		return nil, fmt.Errorf("no routes found")
	}

	// Use the first route (best route)
	route := lifiResp.Routes[0]
	return c.convertToUnifiedQuote(route), nil
}

// GetSupportedChains returns supported chains for bridging
func (c *LiFiClient) GetSupportedChains(ctx context.Context) ([]clients.Chain, error) {
	url := fmt.Sprintf("%s/chains", c.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.apiKey != "" {
		httpReq.Header.Set("x-lifi-api-key", c.apiKey)
	}
	httpReq.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}

	var lifiResp lifiChainsResponse
	if err := clients.ParseResponse(resp, &lifiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	chains := make([]clients.Chain, len(lifiResp.Chains))
	for i, chain := range lifiResp.Chains {
		chains[i] = clients.Chain{
			ID:       strconv.Itoa(chain.ID),
			Name:     chain.Name,
			LogoURI:  chain.LogoURI,
			IsTestnet: chain.MainnetId != chain.ID,
			NativeCurrency: clients.Token{
				Address:  chain.NativeToken.Address,
				Symbol:   chain.NativeToken.Symbol,
				Name:     chain.NativeToken.Name,
				Decimals: chain.NativeToken.Decimals,
				ChainID:  strconv.Itoa(chain.NativeToken.ChainID),
				LogoURI:  chain.NativeToken.LogoURI,
			},
		}
	}

	return chains, nil
}

// GetSupportedTokens returns supported tokens for a specific chain
func (c *LiFiClient) GetSupportedTokens(ctx context.Context, chainID string) ([]clients.Token, error) {
	url := fmt.Sprintf("%s/tokens", c.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add chain filter
	q := httpReq.URL.Query()
	q.Add("chains", chainID)
	httpReq.URL.RawQuery = q.Encode()

	if c.apiKey != "" {
		httpReq.Header.Set("x-lifi-api-key", c.apiKey)
	}
	httpReq.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}

	var lifiResp lifiTokensResponse
	if err := clients.ParseResponse(resp, &lifiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	var tokens []clients.Token
	for _, chainTokens := range lifiResp.Tokens {
		for _, token := range chainTokens {
			tokens = append(tokens, clients.Token{
				Address:  token.Address,
				Symbol:   token.Symbol,
				Name:     token.Name,
				Decimals: token.Decimals,
				ChainID:  strconv.Itoa(token.ChainID),
				LogoURI:  token.LogoURI,
			})
		}
	}

	return tokens, nil
}

// GetProviderName returns the name of the bridge provider
func (c *LiFiClient) GetProviderName() string {
	return "LI.FI"
}

// IsHealthy checks if the provider API is responding
func (c *LiFiClient) IsHealthy(ctx context.Context) bool {
	url := fmt.Sprintf("%s/status", c.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return false
	}

	if c.apiKey != "" {
		httpReq.Header.Set("x-lifi-api-key", c.apiKey)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// convertToUnifiedQuote converts LI.FI route to unified quote format
func (c *LiFiClient) convertToUnifiedQuote(route lifiRoute) *clients.Quote {
	quote := &clients.Quote{
		ID:          route.ID,
		Type:        "bridge",
		Provider:    "LI.FI",
		FromChainID: strconv.Itoa(route.FromChainID),
		ToChainID:   strconv.Itoa(route.ToChainID),
		FromToken: clients.Token{
			Address:  route.FromToken.Address,
			Symbol:   route.FromToken.Symbol,
			Name:     route.FromToken.Name,
			Decimals: route.FromToken.Decimals,
			ChainID:  strconv.Itoa(route.FromToken.ChainID),
			LogoURI:  route.FromToken.LogoURI,
		},
		ToToken: clients.Token{
			Address:  route.ToToken.Address,
			Symbol:   route.ToToken.Symbol,
			Name:     route.ToToken.Name,
			Decimals: route.ToToken.Decimals,
			ChainID:  strconv.Itoa(route.ToToken.ChainID),
			LogoURI:  route.ToToken.LogoURI,
		},
		FromAmount: route.FromAmount,
		ToAmount:   route.ToAmount,
		ExpiresAt:  time.Now().Add(30 * time.Second), // LI.FI quotes expire after 30 seconds
	}

	// Convert steps to route steps
	var routeSteps []clients.RouteStep
	for _, step := range route.Steps {
		routeSteps = append(routeSteps, clients.RouteStep{
			Protocol: step.Tool,
			Type:     step.Type,
			FromToken: clients.Token{
				Address:  step.Action.FromToken.Address,
				Symbol:   step.Action.FromToken.Symbol,
				Name:     step.Action.FromToken.Name,
				Decimals: step.Action.FromToken.Decimals,
				ChainID:  strconv.Itoa(step.Action.FromToken.ChainID),
				LogoURI:  step.Action.FromToken.LogoURI,
			},
			ToToken: clients.Token{
				Address:  step.Action.ToToken.Address,
				Symbol:   step.Action.ToToken.Symbol,
				Name:     step.Action.ToToken.Name,
				Decimals: step.Action.ToToken.Decimals,
				ChainID:  strconv.Itoa(step.Action.ToToken.ChainID),
				LogoURI:  step.Action.ToToken.LogoURI,
			},
			FromAmount: step.Action.FromAmount,
			ToAmount:   step.Action.ToAmount,
		})

		// Add estimated time from first step
		if len(routeSteps) == 1 {
			quote.EstimatedTime = time.Duration(step.Estimate.ExecutionDuration) * time.Second
		}

		// Add transaction data from first step
		if step.TransactionRequest != nil && quote.TransactionData == nil {
			quote.TransactionData = &clients.TransactionData{
				To:       step.TransactionRequest.To,
				Data:     step.TransactionRequest.Data,
				Value:    step.TransactionRequest.Value,
				GasLimit: step.TransactionRequest.GasLimit,
				ChainID:  strconv.Itoa(step.TransactionRequest.ChainID),
			}
			quote.EstimatedGas = step.TransactionRequest.GasLimit
			quote.GasPriceWei = step.TransactionRequest.GasPrice
		}

		// Convert fees
		for _, fee := range step.Estimate.FeeCosts {
			quote.Fees = append(quote.Fees, clients.Fee{
				Type:        fee.Name,
				Amount:      fee.Amount,
				AmountUSD:   fee.AmountUSD,
				Percentage:  fee.Percentage,
				Description: fee.Description,
				Token: clients.Token{
					Address:  fee.Token.Address,
					Symbol:   fee.Token.Symbol,
					Name:     fee.Token.Name,
					Decimals: fee.Token.Decimals,
					ChainID:  strconv.Itoa(fee.Token.ChainID),
					LogoURI:  fee.Token.LogoURI,
				},
			})
		}
	}

	quote.Route = routeSteps
	return quote
}