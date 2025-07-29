package services

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/defi-dashboard/backend/internal/clients"
	"github.com/defi-dashboard/backend/internal/clients/swap"
	"github.com/defi-dashboard/backend/pkg/errors"
)

type SwapService struct {
	zeroXClient   clients.SwapClient
	oneInchClient clients.SwapClient
	cache         clients.Cache
}

func NewSwapService(zeroXConfig, oneInchConfig clients.ClientConfig) *SwapService {
	return &SwapService{
		zeroXClient:   swap.NewZeroXClient(zeroXConfig),
		oneInchClient: swap.NewOneInchClient(oneInchConfig),
		cache:         clients.NewMemoryCache(),
	}
}

type SwapQuoteRequest struct {
	ChainID     int     `json:"chainId"`
	FromToken   string  `json:"fromToken"`
	ToToken     string  `json:"toToken"`
	FromAmount  string  `json:"fromAmount"`
	UserAddress string  `json:"userAddress"`
	Slippage    float64 `json:"slippage"`
	GasPrice    string  `json:"gasPrice,omitempty"`
}

type SwapRoute struct {
	ID           string   `json:"id"`
	FromToken    string   `json:"fromToken"`
	ToToken      string   `json:"toToken"`
	FromAmount   string   `json:"fromAmount"`
	ToAmount     string   `json:"toAmount"`
	EstimatedGas string   `json:"estimatedGas"`
	GasPrice     string   `json:"gasPrice"`
	PriceImpact  float64  `json:"priceImpact"`
	Fees         SwapFees `json:"fees"`
	Path         []string `json:"path"`
	Provider     string   `json:"provider"`
	Dex          string   `json:"dex"`
	Calldata     string   `json:"calldata"`
	Value        string   `json:"value"`
}

type SwapFees struct {
	ProtocolFee string `json:"protocolFee"`
	GasFee      string `json:"gasFee"`
	Total       string `json:"total"`
}

func (s *SwapService) GetQuotes(ctx context.Context, req SwapQuoteRequest) ([]SwapRoute, error) {
	// Convert request to unified format
	quoteReq := clients.QuoteRequest{
		FromChainID: strconv.Itoa(req.ChainID),
		FromToken:   req.FromToken,
		ToToken:     req.ToToken,
		Amount:      req.FromAmount,
		UserAddress: req.UserAddress,
		Slippage:    req.Slippage,
	}

	// Generate cache keys
	zeroXCacheKey := clients.CacheKey{
		Provider:    "0x",
		FromChain:   quoteReq.FromChainID,
		FromToken:   quoteReq.FromToken,
		ToToken:     quoteReq.ToToken,
		Amount:      quoteReq.Amount,
		UserAddress: quoteReq.UserAddress,
	}.String()

	oneInchCacheKey := clients.CacheKey{
		Provider:    "1inch",
		FromChain:   quoteReq.FromChainID,
		FromToken:   quoteReq.FromToken,
		ToToken:     quoteReq.ToToken,
		Amount:      quoteReq.Amount,
		UserAddress: quoteReq.UserAddress,
	}.String()

	var routes []SwapRoute
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Fetch 0x quote
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Check cache first
		if cachedQuote, found := s.cache.Get(zeroXCacheKey); found {
			mu.Lock()
			routes = append(routes, s.convertQuoteToSwapRoute(*cachedQuote, req.GasPrice))
			mu.Unlock()
			return
		}

		// Fetch from API
		quote, err := s.zeroXClient.GetQuote(ctx, quoteReq)
		if err == nil {
			// Cache the quote
			s.cache.Set(zeroXCacheKey, quote, 30*time.Second)
			
			mu.Lock()
			routes = append(routes, s.convertQuoteToSwapRoute(*quote, req.GasPrice))
			mu.Unlock()
		}
	}()

	// Fetch 1inch quote
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Check cache first
		if cachedQuote, found := s.cache.Get(oneInchCacheKey); found {
			mu.Lock()
			routes = append(routes, s.convertQuoteToSwapRoute(*cachedQuote, req.GasPrice))
			mu.Unlock()
			return
		}

		// Fetch from API
		quote, err := s.oneInchClient.GetQuote(ctx, quoteReq)
		if err == nil {
			// Cache the quote
			s.cache.Set(oneInchCacheKey, quote, 60*time.Second)
			
			mu.Lock()
			routes = append(routes, s.convertQuoteToSwapRoute(*quote, req.GasPrice))
			mu.Unlock()
		}
	}()

	wg.Wait()

	if len(routes) == 0 {
		return nil, errors.BadRequest("No swap quotes found")
	}

	return routes, nil
}

// convertQuoteToSwapRoute converts a unified quote to the legacy SwapRoute format
func (s *SwapService) convertQuoteToSwapRoute(quote clients.Quote, gasPrice string) SwapRoute {
	// Use provided gas price or fall back to quote gas price
	finalGasPrice := gasPrice
	if finalGasPrice == "" {
		finalGasPrice = quote.GasPriceWei
		if finalGasPrice == "" {
			finalGasPrice = "20000000000" // 20 gwei default
		}
	}

	// Convert fees
	var protocolFeeTotal float64 = 0
	var gasFeeTotal float64 = 0
	
	for _, fee := range quote.Fees {
		if feeAmount, err := strconv.ParseFloat(fee.AmountUSD, 64); err == nil {
			switch fee.Type {
			case "protocol":
				protocolFeeTotal += feeAmount
			case "gas":
				gasFeeTotal += feeAmount
			}
		}
	}

	// Build path from route steps
	var path []string
	if len(quote.Route) > 0 {
		path = append(path, quote.FromToken.Address)
		for _, step := range quote.Route {
			if step.ToToken.Address != path[len(path)-1] {
				path = append(path, step.ToToken.Address)
			}
		}
	} else {
		path = []string{quote.FromToken.Address, quote.ToToken.Address}
	}

	// Determine primary DEX from route
	dex := "Unknown"
	if len(quote.Route) > 0 {
		dex = quote.Route[0].Protocol
		if len(quote.Route) > 1 {
			dex = "Multiple DEXs"
		}
	}

	// Get transaction data
	calldata := "0x"
	value := "0"
	if quote.TransactionData != nil {
		calldata = quote.TransactionData.Data
		value = quote.TransactionData.Value
	}

	return SwapRoute{
		ID:           quote.ID,
		FromToken:    quote.FromToken.Address,
		ToToken:      quote.ToToken.Address,
		FromAmount:   quote.FromAmount,
		ToAmount:     quote.ToAmount,
		EstimatedGas: quote.EstimatedGas,
		GasPrice:     finalGasPrice,
		PriceImpact:  quote.PriceImpact,
		Fees: SwapFees{
			ProtocolFee: fmt.Sprintf("%.6f", protocolFeeTotal),
			GasFee:      fmt.Sprintf("%.6f", gasFeeTotal),
			Total:       fmt.Sprintf("%.6f", protocolFeeTotal+gasFeeTotal),
		},
		Path:     path,
		Provider: quote.Provider,
		Dex:      dex,
		Calldata: calldata,
		Value:    value,
	}
}