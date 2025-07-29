package clients

import (
	"time"
)

// Common types for bridge and swap operations

// QuoteRequest represents a unified request for quotes
type QuoteRequest struct {
	FromChainID string  `json:"fromChainId"`
	ToChainID   string  `json:"toChainId,omitempty"` // Optional for swaps
	FromToken   string  `json:"fromToken"`
	ToToken     string  `json:"toToken"`
	Amount      string  `json:"amount"`
	UserAddress string  `json:"userAddress"`
	Slippage    float64 `json:"slippage,omitempty"` // Optional, defaults to 0.5%
}

// Quote represents a unified response for quotes
type Quote struct {
	ID                string                 `json:"id"`
	Type              string                 `json:"type"` // "bridge" or "swap"
	Provider          string                 `json:"provider"`
	FromChainID       string                 `json:"fromChainId"`
	ToChainID         string                 `json:"toChainId,omitempty"`
	FromToken         Token                  `json:"fromToken"`
	ToToken           Token                  `json:"toToken"`
	FromAmount        string                 `json:"fromAmount"`
	ToAmount          string                 `json:"toAmount"`
	EstimatedGas      string                 `json:"estimatedGas"`
	GasPriceWei       string                 `json:"gasPriceWei"`
	EstimatedTime     time.Duration          `json:"estimatedTime"`
	ExchangeRate      string                 `json:"exchangeRate"`
	PriceImpact       float64                `json:"priceImpact"`
	Fees              []Fee                  `json:"fees"`
	Route             []RouteStep            `json:"route"`
	TransactionData   *TransactionData       `json:"transactionData,omitempty"`
	ExpiresAt         time.Time              `json:"expiresAt"`
	AdditionalData    map[string]interface{} `json:"additionalData,omitempty"`
}

// Token represents token information
type Token struct {
	Address  string `json:"address"`
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Decimals int    `json:"decimals"`
	ChainID  string `json:"chainId"`
	LogoURI  string `json:"logoUri,omitempty"`
}

// Fee represents various fees
type Fee struct {
	Type        string `json:"type"` // "bridge", "gas", "protocol", "slippage"
	Amount      string `json:"amount"`
	Token       Token  `json:"token"`
	AmountUSD   string `json:"amountUsd,omitempty"`
	Percentage  string `json:"percentage,omitempty"`
	Description string `json:"description,omitempty"`
}

// RouteStep represents a step in the route
type RouteStep struct {
	Protocol     string `json:"protocol"`
	Type         string `json:"type"` // "swap", "bridge"
	FromToken    Token  `json:"fromToken"`
	ToToken      Token  `json:"toToken"`
	FromAmount   string `json:"fromAmount"`
	ToAmount     string `json:"toAmount"`
	ExchangeRate string `json:"exchangeRate"`
	Percentage   string `json:"percentage,omitempty"`
}

// TransactionData represents transaction data for execution
type TransactionData struct {
	To       string `json:"to"`
	Data     string `json:"data"`
	Value    string `json:"value"`
	GasLimit string `json:"gasLimit"`
	ChainID  string `json:"chainId"`
}

// ErrorResponse represents an error from external APIs
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// CacheKey represents a cache key for quotes
type CacheKey struct {
	Provider    string
	FromChain   string
	ToChain     string
	FromToken   string
	ToToken     string
	Amount      string
	UserAddress string
}

// String returns a string representation of the cache key
func (c CacheKey) String() string {
	if c.ToChain == "" {
		// Swap key
		return c.Provider + ":" + c.FromChain + ":" + c.FromToken + ":" + c.ToToken + ":" + c.Amount + ":" + c.UserAddress
	}
	// Bridge key
	return c.Provider + ":" + c.FromChain + ":" + c.ToChain + ":" + c.FromToken + ":" + c.ToToken + ":" + c.Amount + ":" + c.UserAddress
}