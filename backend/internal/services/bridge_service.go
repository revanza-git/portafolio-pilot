package services

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/defi-dashboard/backend/internal/clients"
	"github.com/defi-dashboard/backend/internal/clients/bridge"
	"github.com/defi-dashboard/backend/pkg/errors"
)

type BridgeService struct {
	lifiClient   clients.BridgeClient
	socketClient clients.BridgeClient
	cache        clients.Cache
}

func NewBridgeService(lifiConfig, socketConfig clients.ClientConfig) *BridgeService {
	return &BridgeService{
		lifiClient:   bridge.NewLiFiClient(lifiConfig),
		socketClient: bridge.NewSocketClient(socketConfig),
		cache:        clients.NewMemoryCache(),
	}
}

type BridgeRouteRequest struct {
	FromChain   int    `json:"fromChain"`
	ToChain     int    `json:"toChain"`
	FromToken   string `json:"fromToken"`
	ToToken     string `json:"toToken"`
	FromAmount  string `json:"fromAmount"`
	UserAddress string `json:"userAddress"`
	Slippage    float64 `json:"slippage"`
}

type BridgeRoute struct {
	ID            string       `json:"id"`
	FromChain     int          `json:"fromChain"`
	ToChain       int          `json:"toChain"`
	FromToken     string       `json:"fromToken"`
	ToToken       string       `json:"toToken"`
	FromAmount    string       `json:"fromAmount"`
	ToAmount      string       `json:"toAmount"`
	EstimatedGas  string       `json:"estimatedGas"`
	EstimatedTime int          `json:"estimatedTime"`
	Fees          BridgeFees   `json:"fees"`
	Steps         []BridgeStep `json:"steps"`
	Provider      string       `json:"provider"`
}

type BridgeFees struct {
	BridgeFee string `json:"bridgeFee"`
	GasFee    string `json:"gasFee"`
	Total     string `json:"total"`
}

type BridgeStep struct {
	Type      string `json:"type"`
	Protocol  string `json:"protocol"`
	FromChain int    `json:"fromChain"`
	ToChain   int    `json:"toChain"`
	FromToken string `json:"fromToken"`
	ToToken   string `json:"toToken"`
	FromAmount string `json:"fromAmount"`
	ToAmount  string `json:"toAmount"`
	Data      string `json:"data"`
	Value     string `json:"value"`
	GasLimit  string `json:"gasLimit"`
}

func (s *BridgeService) GetRoutes(ctx context.Context, req BridgeRouteRequest) ([]BridgeRoute, error) {
	// Convert request to unified format
	quoteReq := clients.QuoteRequest{
		FromChainID: strconv.Itoa(req.FromChain),
		ToChainID:   strconv.Itoa(req.ToChain),
		FromToken:   req.FromToken,
		ToToken:     req.ToToken,
		Amount:      req.FromAmount,
		UserAddress: req.UserAddress,
		Slippage:    req.Slippage,
	}

	// Generate cache keys
	lifiCacheKey := clients.CacheKey{
		Provider:    "lifi",
		FromChain:   quoteReq.FromChainID,
		ToChain:     quoteReq.ToChainID,
		FromToken:   quoteReq.FromToken,
		ToToken:     quoteReq.ToToken,
		Amount:      quoteReq.Amount,
		UserAddress: quoteReq.UserAddress,
	}.String()

	socketCacheKey := clients.CacheKey{
		Provider:    "socket",
		FromChain:   quoteReq.FromChainID,
		ToChain:     quoteReq.ToChainID,
		FromToken:   quoteReq.FromToken,
		ToToken:     quoteReq.ToToken,
		Amount:      quoteReq.Amount,
		UserAddress: quoteReq.UserAddress,
	}.String()

	var routes []BridgeRoute
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Fetch LiFi quote
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Check cache first
		if cachedQuote, found := s.cache.Get(lifiCacheKey); found {
			mu.Lock()
			routes = append(routes, s.convertQuoteToBridgeRoute(*cachedQuote))
			mu.Unlock()
			return
		}

		// Fetch from API
		quote, err := s.lifiClient.GetQuote(ctx, quoteReq)
		if err == nil {
			// Cache the quote
			s.cache.Set(lifiCacheKey, quote, 30*time.Second)
			
			mu.Lock()
			routes = append(routes, s.convertQuoteToBridgeRoute(*quote))
			mu.Unlock()
		}
	}()

	// Fetch Socket quote
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Check cache first
		if cachedQuote, found := s.cache.Get(socketCacheKey); found {
			mu.Lock()
			routes = append(routes, s.convertQuoteToBridgeRoute(*cachedQuote))
			mu.Unlock()
			return
		}

		// Fetch from API
		quote, err := s.socketClient.GetQuote(ctx, quoteReq)
		if err == nil {
			// Cache the quote
			s.cache.Set(socketCacheKey, quote, 60*time.Second)
			
			mu.Lock()
			routes = append(routes, s.convertQuoteToBridgeRoute(*quote))
			mu.Unlock()
		}
	}()

	wg.Wait()

	if len(routes) == 0 {
		return nil, errors.BadRequest("No bridge routes found")
	}

	return routes, nil
}

// convertQuoteToBridgeRoute converts a unified quote to the legacy BridgeRoute format
func (s *BridgeService) convertQuoteToBridgeRoute(quote clients.Quote) BridgeRoute {
	fromChain, _ := strconv.Atoi(quote.FromChainID)
	toChain, _ := strconv.Atoi(quote.ToChainID)
	
	// Convert fees
	var bridgeFeeTotal float64 = 0
	var gasFeeTotal float64 = 0
	
	for _, fee := range quote.Fees {
		if feeAmount, err := strconv.ParseFloat(fee.AmountUSD, 64); err == nil {
			switch fee.Type {
			case "bridge", "protocol":
				bridgeFeeTotal += feeAmount
			case "gas":
				gasFeeTotal += feeAmount
			}
		}
	}

	// Convert route steps
	var steps []BridgeStep
	for _, step := range quote.Route {
		fromChainInt, _ := strconv.Atoi(step.FromToken.ChainID)
		toChainInt, _ := strconv.Atoi(step.ToToken.ChainID)
		
		bridgeStep := BridgeStep{
			Type:       step.Type,
			Protocol:   step.Protocol,
			FromChain:  fromChainInt,
			ToChain:    toChainInt,
			FromToken:  step.FromToken.Address,
			ToToken:    step.ToToken.Address,
			FromAmount: step.FromAmount,
			ToAmount:   step.ToAmount,
		}

		// Add transaction data if available
		if quote.TransactionData != nil {
			bridgeStep.Data = quote.TransactionData.Data
			bridgeStep.Value = quote.TransactionData.Value
			bridgeStep.GasLimit = quote.TransactionData.GasLimit
		}

		steps = append(steps, bridgeStep)
	}

	return BridgeRoute{
		ID:            quote.ID,
		FromChain:     fromChain,
		ToChain:       toChain,
		FromToken:     quote.FromToken.Address,
		ToToken:       quote.ToToken.Address,
		FromAmount:    quote.FromAmount,
		ToAmount:      quote.ToAmount,
		EstimatedGas:  quote.EstimatedGas,
		EstimatedTime: int(quote.EstimatedTime.Seconds()),
		Fees: BridgeFees{
			BridgeFee: fmt.Sprintf("%.6f", bridgeFeeTotal),
			GasFee:    fmt.Sprintf("%.6f", gasFeeTotal),
			Total:     fmt.Sprintf("%.6f", bridgeFeeTotal+gasFeeTotal),
		},
		Steps:    steps,
		Provider: quote.Provider,
	}
}