package clients

import (
	"context"
	"net/http"
	"time"
)

// BridgeClient defines the interface for bridge providers
type BridgeClient interface {
	// GetQuote fetches a bridge quote from the provider
	GetQuote(ctx context.Context, req QuoteRequest) (*Quote, error)
	
	// GetSupportedChains returns supported chains for bridging
	GetSupportedChains(ctx context.Context) ([]Chain, error)
	
	// GetSupportedTokens returns supported tokens for a specific chain
	GetSupportedTokens(ctx context.Context, chainID string) ([]Token, error)
	
	// GetProviderName returns the name of the bridge provider
	GetProviderName() string
	
	// IsHealthy checks if the provider API is responding
	IsHealthy(ctx context.Context) bool
}

// SwapClient defines the interface for swap providers
type SwapClient interface {
	// GetQuote fetches a swap quote from the provider
	GetQuote(ctx context.Context, req QuoteRequest) (*Quote, error)
	
	// GetSupportedTokens returns supported tokens for a specific chain
	GetSupportedTokens(ctx context.Context, chainID string) ([]Token, error)
	
	// GetProviderName returns the name of the swap provider
	GetProviderName() string
	
	// IsHealthy checks if the provider API is responding
	IsHealthy(ctx context.Context) bool
}

// HTTPClient defines the interface for HTTP operations with retry logic
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
	Get(url string) (*http.Response, error)
	Post(url, contentType string, body interface{}) (*http.Response, error)
}

// Cache defines the interface for caching quotes
type Cache interface {
	Get(key string) (*Quote, bool)
	Set(key string, quote *Quote, ttl time.Duration)
	Delete(key string)
	Clear()
}

// Chain represents blockchain information
type Chain struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	NativeCurrency Token `json:"nativeCurrency"`
	LogoURI  string `json:"logoUri,omitempty"`
	IsTestnet bool  `json:"isTestnet"`
}

// ClientConfig holds configuration for external API clients
type ClientConfig struct {
	BaseURL     string
	APIKey      string
	Timeout     time.Duration
	MaxRetries  int
	RetryDelay  time.Duration
	RateLimit   RateLimitConfig
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	RequestsPerSecond int
	BurstSize         int
}

// HealthCheckResponse represents a health check response
type HealthCheckResponse struct {
	Provider string `json:"provider"`
	Status   string `json:"status"` // "healthy", "degraded", "unhealthy"
	Latency  time.Duration `json:"latency"`
	Error    string `json:"error,omitempty"`
}