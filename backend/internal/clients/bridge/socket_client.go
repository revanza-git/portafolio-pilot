package bridge

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/defi-dashboard/backend/internal/clients"
)

// SocketClient implements BridgeClient for Socket API
type SocketClient struct {
	httpClient clients.HTTPClient
	baseURL    string
	apiKey     string
}

// Socket API types
type socketQuoteRequest struct {
	FromChainId     string  `json:"fromChainId"`
	FromTokenAddress string  `json:"fromTokenAddress"`
	ToChainId       string  `json:"toChainId"`
	ToTokenAddress  string  `json:"toTokenAddress"`
	FromAmount      string  `json:"fromAmount"`
	UserAddress     string  `json:"userAddress"`
	UniqueRoutesPerBridge bool `json:"uniqueRoutesPerBridge,omitempty"`
	Sort            string  `json:"sort,omitempty"`
	Singletxn       bool    `json:"singleTxn,omitempty"`
}

type socketQuoteResponse struct {
	Success bool           `json:"success"`
	Result  socketResult   `json:"result"`
}

type socketResult struct {
	Routes                []socketRoute `json:"routes"`
	FromChainId           int           `json:"fromChainId"`
	FromAsset             socketToken   `json:"fromAsset"`
	ToChainId             int           `json:"toChainId"`
	ToAsset               socketToken   `json:"toAsset"`
	FromAmount            string        `json:"fromAmount"`
}

type socketRoute struct {
	RouteId              string               `json:"routeId"`
	IsOnlySwapRoute      bool                 `json:"isOnlySwapRoute"`
	FromAmount           string               `json:"fromAmount"`
	ToAmount             string               `json:"toAmount"`
	UsedBridgeNames      []string             `json:"usedBridgeNames"`
	TotalUserTx          int                  `json:"totalUserTx"`
	Sender               string               `json:"sender"`
	Recipient            string               `json:"recipient"`
	TotalGasFeesInUsd    float64              `json:"totalGasFeesInUsd"`
	ReceiveValueInUsd    float64              `json:"receiveValueInUsd"`
	InputValueInUsd      float64              `json:"inputValueInUsd"`
	OutputValueInUsd     float64              `json:"outputValueInUsd"`
	UserTxs              []socketUserTx       `json:"userTxs"`
	ServiceTime          int                  `json:"serviceTime"`
	MaxServiceTime       int                  `json:"maxServiceTime"`
	IntegratorFee        socketIntegratorFee  `json:"integratorFee"`
}

type socketUserTx struct {
	UserTxType        string               `json:"userTxType"`
	TxType            string               `json:"txType"`
	ChainId           int                  `json:"chainId"`
	ToAmount          string               `json:"toAmount"`
	ToAsset           socketToken          `json:"toAsset"`
	StepCount         int                  `json:"stepCount"`
	RoutePath         string               `json:"routePath"`
	Sender            string               `json:"sender"`
	ApprovalData      *socketApprovalData  `json:"approvalData,omitempty"`
	Steps             []socketStep         `json:"steps"`
	GasFees           socketGasFees        `json:"gasFees"`
}

type socketStep struct {
	Type             string              `json:"type"`
	Protocol         socketProtocol      `json:"protocol"`
	FromChainId      int                 `json:"fromChainId"`
	FromAsset        socketToken         `json:"fromAsset"`
	FromAmount       string              `json:"fromAmount"`
	ToChainId        int                 `json:"toChainId"`
	ToAsset          socketToken         `json:"toAsset"`
	ToAmount         string              `json:"toAmount"`
	BridgeSlippage   float64             `json:"bridgeSlippage,omitempty"`
	SwapSlippage     float64             `json:"swapSlippage,omitempty"`
	ServiceTime      int                 `json:"serviceTime"`
	MaxServiceTime   int                 `json:"maxServiceTime"`
}

type socketProtocol struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Icon        string `json:"icon"`
}

type socketToken struct {
	ChainId      int     `json:"chainId"`
	Address      string  `json:"address"`
	Name         string  `json:"name"`
	Symbol       string  `json:"symbol"`
	Decimals     int     `json:"decimals"`
	Icon         string  `json:"icon"`
	LogoURI      string  `json:"logoURI"`
	ChainAgnosticId string `json:"chainAgnosticId,omitempty"`
}

type socketApprovalData struct {
	MinimumApprovalAmount string `json:"minimumApprovalAmount"`
	ApprovalTokenAddress  string `json:"approvalTokenAddress"`
	AllowanceTarget       string `json:"allowanceTarget"`
	Owner                 string `json:"owner"`
}

type socketGasFees struct {
	GasAmount string      `json:"gasAmount"`
	GasLimit  int         `json:"gasLimit"`
	Asset     socketToken `json:"asset"`
	FeesInUsd float64     `json:"feesInUsd"`
}

type socketIntegratorFee struct {
	Amount string      `json:"amount"`
	Asset  socketToken `json:"asset"`
}

type socketSupportedResponse struct {
	Success bool                   `json:"success"`
	Result  socketSupportedResult  `json:"result"`
}

type socketSupportedResult struct {
	FromChainId int                        `json:"fromChainId"`
	ToChainId   int                        `json:"toChainId"`
	Result      []socketSupportedToken     `json:"result"`
}

type socketSupportedToken struct {
	Address     string `json:"address"`
	ChainId     int    `json:"chainId"`
	Currency    string `json:"currency"`
	Decimals    int    `json:"decimals"`
	Icon        string `json:"icon"`
	LogoURI     string `json:"logoURI"`
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
}

type socketChainResponse struct {
	Success bool               `json:"success"`
	Result  []socketChainInfo  `json:"result"`
}

type socketChainInfo struct {
	ChainId      int                 `json:"chainId"`
	Name         string              `json:"name"`
	IsL1         bool                `json:"isL1"`
	SendingEnabled bool              `json:"sendingEnabled"`
	ReceivingEnabled bool            `json:"receivingEnabled"`
	RefuelEnabled bool               `json:"refuelEnabled"`
	Icon         string              `json:"icon"`
	LogoURI      string              `json:"logoURI"`
	Currency     socketChainCurrency `json:"currency"`
	Rpcs         []string            `json:"rpcs"`
	Explorers    []string            `json:"explorers"`
}

type socketChainCurrency struct {
	Address     string `json:"address"`
	Icon        string `json:"icon"`
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Decimals    int    `json:"decimals"`
	ChainId     int    `json:"chainId"`
}

// NewSocketClient creates a new Socket bridge client
func NewSocketClient(config clients.ClientConfig) *SocketClient {
	httpClient := clients.NewBaseHTTPClient(config)
	
	return &SocketClient{
		httpClient: httpClient,
		baseURL:    config.BaseURL,
		apiKey:     config.APIKey,
	}
}

// GetQuote fetches a bridge quote from Socket
func (c *SocketClient) GetQuote(ctx context.Context, req clients.QuoteRequest) (*clients.Quote, error) {
	url := fmt.Sprintf("%s/quote", c.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add query parameters
	q := httpReq.URL.Query()
	q.Add("fromChainId", req.FromChainID)
	q.Add("fromTokenAddress", req.FromToken)
	q.Add("toChainId", req.ToChainID)
	q.Add("toTokenAddress", req.ToToken)
	q.Add("fromAmount", req.Amount)
	q.Add("userAddress", req.UserAddress)
	q.Add("uniqueRoutesPerBridge", "true")
	q.Add("sort", "output")
	q.Add("singleTxn", "false")
	httpReq.URL.RawQuery = q.Encode()

	// Add headers
	if c.apiKey != "" {
		httpReq.Header.Set("API-KEY", c.apiKey)
	}
	httpReq.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}

	var socketResp socketQuoteResponse
	if err := clients.ParseResponse(resp, &socketResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if !socketResp.Success {
		return nil, fmt.Errorf("Socket API returned error")
	}

	// Convert to unified format
	if len(socketResp.Result.Routes) == 0 {
		return nil, fmt.Errorf("no routes found")
	}

	// Use the first route (best route by output)
	route := socketResp.Result.Routes[0]
	return c.convertToUnifiedQuote(route, socketResp.Result), nil
}

// GetSupportedChains returns supported chains for bridging
func (c *SocketClient) GetSupportedChains(ctx context.Context) ([]clients.Chain, error) {
	url := fmt.Sprintf("%s/supported/chains", c.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.apiKey != "" {
		httpReq.Header.Set("API-KEY", c.apiKey)
	}
	httpReq.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}

	var socketResp socketChainResponse
	if err := clients.ParseResponse(resp, &socketResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if !socketResp.Success {
		return nil, fmt.Errorf("Socket API returned error")
	}

	chains := make([]clients.Chain, len(socketResp.Result))
	for i, chain := range socketResp.Result {
		chains[i] = clients.Chain{
			ID:       strconv.Itoa(chain.ChainId),
			Name:     chain.Name,
			LogoURI:  chain.LogoURI,
			IsTestnet: false, // Socket doesn't distinguish testnets clearly
			NativeCurrency: clients.Token{
				Address:  chain.Currency.Address,
				Symbol:   chain.Currency.Symbol,
				Name:     chain.Currency.Name,
				Decimals: chain.Currency.Decimals,
				ChainID:  strconv.Itoa(chain.Currency.ChainId),
			},
		}
	}

	return chains, nil
}

// GetSupportedTokens returns supported tokens for a specific chain
func (c *SocketClient) GetSupportedTokens(ctx context.Context, chainID string) ([]clients.Token, error) {
	url := fmt.Sprintf("%s/token-lists/from-token-list", c.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add query parameters
	q := httpReq.URL.Query()
	q.Add("fromChainId", chainID)
	q.Add("toChainId", chainID) // Same chain for token list
	httpReq.URL.RawQuery = q.Encode()

	if c.apiKey != "" {
		httpReq.Header.Set("API-KEY", c.apiKey)
	}
	httpReq.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}

	var socketResp socketSupportedResponse
	if err := clients.ParseResponse(resp, &socketResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if !socketResp.Success {
		return nil, fmt.Errorf("Socket API returned error")
	}

	tokens := make([]clients.Token, len(socketResp.Result.Result))
	for i, token := range socketResp.Result.Result {
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

// GetProviderName returns the name of the bridge provider
func (c *SocketClient) GetProviderName() string {
	return "Socket"
}

// IsHealthy checks if the provider API is responding
func (c *SocketClient) IsHealthy(ctx context.Context) bool {
	url := fmt.Sprintf("%s/supported/chains", c.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return false
	}

	if c.apiKey != "" {
		httpReq.Header.Set("API-KEY", c.apiKey)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// convertToUnifiedQuote converts Socket route to unified quote format
func (c *SocketClient) convertToUnifiedQuote(route socketRoute, result socketResult) *clients.Quote {
	quote := &clients.Quote{
		ID:          route.RouteId,
		Type:        "bridge",
		Provider:    "Socket",
		FromChainID: strconv.Itoa(result.FromChainId),
		ToChainID:   strconv.Itoa(result.ToChainId),
		FromToken: clients.Token{
			Address:  result.FromAsset.Address,
			Symbol:   result.FromAsset.Symbol,
			Name:     result.FromAsset.Name,
			Decimals: result.FromAsset.Decimals,
			ChainID:  strconv.Itoa(result.FromAsset.ChainId),
			LogoURI:  result.FromAsset.LogoURI,
		},
		ToToken: clients.Token{
			Address:  result.ToAsset.Address,
			Symbol:   result.ToAsset.Symbol,
			Name:     result.ToAsset.Name,
			Decimals: result.ToAsset.Decimals,
			ChainID:  strconv.Itoa(result.ToAsset.ChainId),
			LogoURI:  result.ToAsset.LogoURI,
		},
		FromAmount:    route.FromAmount,
		ToAmount:      route.ToAmount,
		EstimatedTime: time.Duration(route.ServiceTime) * time.Second,
		ExpiresAt:     time.Now().Add(60 * time.Second), // Socket quotes expire after 60 seconds
	}

	// Convert user transactions to route steps
	var routeSteps []clients.RouteStep
	for _, userTx := range route.UserTxs {
		for _, step := range userTx.Steps {
			routeSteps = append(routeSteps, clients.RouteStep{
				Protocol: step.Protocol.Name,
				Type:     step.Type,
				FromToken: clients.Token{
					Address:  step.FromAsset.Address,
					Symbol:   step.FromAsset.Symbol,
					Name:     step.FromAsset.Name,
					Decimals: step.FromAsset.Decimals,
					ChainID:  strconv.Itoa(step.FromAsset.ChainId),
					LogoURI:  step.FromAsset.LogoURI,
				},
				ToToken: clients.Token{
					Address:  step.ToAsset.Address,
					Symbol:   step.ToAsset.Symbol,
					Name:     step.ToAsset.Name,
					Decimals: step.ToAsset.Decimals,
					ChainID:  strconv.Itoa(step.ToAsset.ChainId),
					LogoURI:  step.ToAsset.LogoURI,
				},
				FromAmount: step.FromAmount,
				ToAmount:   step.ToAmount,
			})
		}

		// Add gas fees
		if userTx.GasFees.GasAmount != "" {
			quote.Fees = append(quote.Fees, clients.Fee{
				Type:      "gas",
				Amount:    userTx.GasFees.GasAmount,
				AmountUSD: fmt.Sprintf("%.6f", userTx.GasFees.FeesInUsd),
				Token: clients.Token{
					Address:  userTx.GasFees.Asset.Address,
					Symbol:   userTx.GasFees.Asset.Symbol,
					Name:     userTx.GasFees.Asset.Name,
					Decimals: userTx.GasFees.Asset.Decimals,
					ChainID:  strconv.Itoa(userTx.GasFees.Asset.ChainId),
					LogoURI:  userTx.GasFees.Asset.LogoURI,
				},
			})
		}

		// Set gas estimate from first transaction
		if quote.EstimatedGas == "" {
			quote.EstimatedGas = strconv.Itoa(userTx.GasFees.GasLimit)
		}
	}

	// Add integrator fee if present
	if route.IntegratorFee.Amount != "" {
		quote.Fees = append(quote.Fees, clients.Fee{
			Type:   "protocol",
			Amount: route.IntegratorFee.Amount,
			Token: clients.Token{
				Address:  route.IntegratorFee.Asset.Address,
				Symbol:   route.IntegratorFee.Asset.Symbol,
				Name:     route.IntegratorFee.Asset.Name,
				Decimals: route.IntegratorFee.Asset.Decimals,
				ChainID:  strconv.Itoa(route.IntegratorFee.Asset.ChainId),
				LogoURI:  route.IntegratorFee.Asset.LogoURI,
			},
		})
	}

	quote.Route = routeSteps
	return quote
}