package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/defi-dashboard/backend/internal/clients"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	// Server
	Port     string
	LogLevel string

	// Database
	DatabaseURL string

	// JWT
	JWTSecret string
	JWTExpiry int // hours

	// API
	APIVersion   string
	AllowOrigins string

	// External Services
	AlchemyAPIKey   string
	InfuraAPIKey    string
	EtherscanAPIKey string
	CoinGeckoAPIKey string
	DefiLlamaEnabled bool

	// Bridge Clients
	LiFiAPIKey   string
	LiFiBaseURL  string
	SocketAPIKey string
	SocketBaseURL string

	// Swap Clients
	ZeroXAPIKey   string
	ZeroXBaseURL  string
	OneInchAPIKey string
	OneInchBaseURL string

	// External API Settings
	ExternalAPITimeout     int
	ExternalAPIMaxRetries  int
	ExternalAPIRetryDelay  int
	ExternalAPIRateLimitRPS int
	ExternalAPIRateLimitBurst int

	// Redis (optional)
	RedisURL string
}

func Load() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// Not an error if .env doesn't exist
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("error loading .env file: %w", err)
		}
	}

	// Set up viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(filepath.Join("..", "config"))

	// Read config file if exists
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Allow environment variables to override
	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("PORT", "3000")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("API_VERSION", "v1")
	viper.SetDefault("JWT_EXPIRY", 24)
	viper.SetDefault("ALLOW_ORIGINS", "*")
	viper.SetDefault("DEFILLAMA_ENABLED", true)
	
	// External API defaults
	viper.SetDefault("LIFI_BASE_URL", "https://li.quest/v1")
	viper.SetDefault("SOCKET_BASE_URL", "https://api.socket.tech/v2")
	viper.SetDefault("ZEROX_BASE_URL", "https://api.0x.org")
	viper.SetDefault("ONEINCH_BASE_URL", "https://api.1inch.io")
	viper.SetDefault("EXTERNAL_API_TIMEOUT", 30000)
	viper.SetDefault("EXTERNAL_API_MAX_RETRIES", 3)
	viper.SetDefault("EXTERNAL_API_RETRY_DELAY", 1000)
	viper.SetDefault("EXTERNAL_API_RATE_LIMIT_RPS", 10)
	viper.SetDefault("EXTERNAL_API_RATE_LIMIT_BURST", 20)

	cfg := &Config{
		Port:            viper.GetString("PORT"),
		LogLevel:        viper.GetString("LOG_LEVEL"),
		DatabaseURL:     viper.GetString("DATABASE_URL"),
		JWTSecret:       viper.GetString("JWT_SECRET"),
		JWTExpiry:       viper.GetInt("JWT_EXPIRY"),
		APIVersion:      viper.GetString("API_VERSION"),
		AllowOrigins:    viper.GetString("ALLOW_ORIGINS"),
		AlchemyAPIKey:   viper.GetString("ALCHEMY_API_KEY"),
		InfuraAPIKey:    viper.GetString("INFURA_API_KEY"),
		EtherscanAPIKey: viper.GetString("ETHERSCAN_API_KEY"),
		CoinGeckoAPIKey: viper.GetString("COINGECKO_API_KEY"),
		DefiLlamaEnabled: viper.GetBool("DEFILLAMA_ENABLED"),
		
		// Bridge Clients
		LiFiAPIKey:      viper.GetString("LIFI_API_KEY"),
		LiFiBaseURL:     viper.GetString("LIFI_BASE_URL"),
		SocketAPIKey:    viper.GetString("SOCKET_API_KEY"),
		SocketBaseURL:   viper.GetString("SOCKET_BASE_URL"),
		
		// Swap Clients
		ZeroXAPIKey:     viper.GetString("ZEROX_API_KEY"),
		ZeroXBaseURL:    viper.GetString("ZEROX_BASE_URL"),
		OneInchAPIKey:   viper.GetString("ONEINCH_API_KEY"),
		OneInchBaseURL:  viper.GetString("ONEINCH_BASE_URL"),
		
		// External API Settings
		ExternalAPITimeout:        viper.GetInt("EXTERNAL_API_TIMEOUT"),
		ExternalAPIMaxRetries:     viper.GetInt("EXTERNAL_API_MAX_RETRIES"),
		ExternalAPIRetryDelay:     viper.GetInt("EXTERNAL_API_RETRY_DELAY"),
		ExternalAPIRateLimitRPS:   viper.GetInt("EXTERNAL_API_RATE_LIMIT_RPS"),
		ExternalAPIRateLimitBurst: viper.GetInt("EXTERNAL_API_RATE_LIMIT_BURST"),
		
		RedisURL:        viper.GetString("REDIS_URL"),
	}

	// Validate required fields
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

// GetLiFiClientConfig returns client configuration for LI.FI
func (c *Config) GetLiFiClientConfig() clients.ClientConfig {
	return clients.ClientConfig{
		BaseURL:    c.LiFiBaseURL,
		APIKey:     c.LiFiAPIKey,
		Timeout:    time.Duration(c.ExternalAPITimeout) * time.Millisecond,
		MaxRetries: c.ExternalAPIMaxRetries,
		RetryDelay: time.Duration(c.ExternalAPIRetryDelay) * time.Millisecond,
		RateLimit: clients.RateLimitConfig{
			RequestsPerSecond: c.ExternalAPIRateLimitRPS,
			BurstSize:         c.ExternalAPIRateLimitBurst,
		},
	}
}

// GetSocketClientConfig returns client configuration for Socket
func (c *Config) GetSocketClientConfig() clients.ClientConfig {
	return clients.ClientConfig{
		BaseURL:    c.SocketBaseURL,
		APIKey:     c.SocketAPIKey,
		Timeout:    time.Duration(c.ExternalAPITimeout) * time.Millisecond,
		MaxRetries: c.ExternalAPIMaxRetries,
		RetryDelay: time.Duration(c.ExternalAPIRetryDelay) * time.Millisecond,
		RateLimit: clients.RateLimitConfig{
			RequestsPerSecond: c.ExternalAPIRateLimitRPS,
			BurstSize:         c.ExternalAPIRateLimitBurst,
		},
	}
}

// GetZeroXClientConfig returns client configuration for 0x
func (c *Config) GetZeroXClientConfig() clients.ClientConfig {
	return clients.ClientConfig{
		BaseURL:    c.ZeroXBaseURL,
		APIKey:     c.ZeroXAPIKey,
		Timeout:    time.Duration(c.ExternalAPITimeout) * time.Millisecond,
		MaxRetries: c.ExternalAPIMaxRetries,
		RetryDelay: time.Duration(c.ExternalAPIRetryDelay) * time.Millisecond,
		RateLimit: clients.RateLimitConfig{
			RequestsPerSecond: c.ExternalAPIRateLimitRPS,
			BurstSize:         c.ExternalAPIRateLimitBurst,
		},
	}
}

// GetOneInchClientConfig returns client configuration for 1inch
func (c *Config) GetOneInchClientConfig() clients.ClientConfig {
	return clients.ClientConfig{
		BaseURL:    c.OneInchBaseURL,
		APIKey:     c.OneInchAPIKey,
		Timeout:    time.Duration(c.ExternalAPITimeout) * time.Millisecond,
		MaxRetries: c.ExternalAPIMaxRetries,
		RetryDelay: time.Duration(c.ExternalAPIRetryDelay) * time.Millisecond,
		RateLimit: clients.RateLimitConfig{
			RequestsPerSecond: c.ExternalAPIRateLimitRPS,
			BurstSize:         c.ExternalAPIRateLimitBurst,
		},
	}
}