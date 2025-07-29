package repos

import (
	"context"

	"github.com/defi-dashboard/backend/internal/models"
	"github.com/defi-dashboard/backend/pkg/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type tokenRepository struct {
	db *pgxpool.Pool
}

// NewTokenRepository creates a new token repository
func NewTokenRepository(db *pgxpool.Pool) TokenRepository {
	return &tokenRepository{db: db}
}

func (r *tokenRepository) GetByAddress(ctx context.Context, address string, chainID int) (*models.Token, error) {
	// TODO: Implement actual database query
	// Mock data for common tokens
	mockTokens := map[string]*models.Token{
		"0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48": {
			ID:       uuid.New(),
			Address:  "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
			ChainID:  1,
			Symbol:   "USDC",
			Name:     "USD Coin",
			Decimals: 6,
			LogoURI:  utils.StrPtr("https://assets.coingecko.com/coins/images/6319/small/USD_Coin_icon.png"),
			PriceUSD: utils.Float64Ptr(1.0),
			PriceChange24h: utils.Float64Ptr(0.01),
			MarketCap: utils.Float64Ptr(25000000000),
		},
		"0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2": {
			ID:       uuid.New(),
			Address:  "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
			ChainID:  1,
			Symbol:   "WETH",
			Name:     "Wrapped Ether",
			Decimals: 18,
			LogoURI:  utils.StrPtr("https://assets.coingecko.com/coins/images/2518/small/weth.png"),
			PriceUSD: utils.Float64Ptr(2345.67),
			PriceChange24h: utils.Float64Ptr(2.34),
			MarketCap: utils.Float64Ptr(8000000000),
		},
	}

	if token, ok := mockTokens[address]; ok && token.ChainID == chainID {
		return token, nil
	}

	// Default mock token
	return &models.Token{
		ID:       uuid.New(),
		Address:  address,
		ChainID:  chainID,
		Symbol:   "MOCK",
		Name:     "Mock Token",
		Decimals: 18,
		PriceUSD: utils.Float64Ptr(10.0),
	}, nil
}

func (r *tokenRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Token, error) {
	// TODO: Implement actual database query
	return &models.Token{
		ID:       id,
		Address:  "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
		ChainID:  1,
		Symbol:   "USDC",
		Name:     "USD Coin",
		Decimals: 6,
		PriceUSD: utils.Float64Ptr(1.0),
	}, nil
}

func (r *tokenRepository) GetByChainID(ctx context.Context, chainID int, limit, offset int) ([]*models.Token, error) {
	// TODO: Implement actual database query
	// Return mock tokens
	return []*models.Token{
		{
			ID:       uuid.New(),
			Address:  "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
			ChainID:  chainID,
			Symbol:   "USDC",
			Name:     "USD Coin",
			Decimals: 6,
			PriceUSD: utils.Float64Ptr(1.0),
			MarketCap: utils.Float64Ptr(25000000000),
		},
		{
			ID:       uuid.New(),
			Address:  "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
			ChainID:  chainID,
			Symbol:   "WETH",
			Name:     "Wrapped Ether",
			Decimals: 18,
			PriceUSD: utils.Float64Ptr(2345.67),
			MarketCap: utils.Float64Ptr(8000000000),
		},
		{
			ID:       uuid.New(),
			Address:  "0x6B175474E89094C44Da98b954EedeAC495271d0F",
			ChainID:  chainID,
			Symbol:   "DAI",
			Name:     "Dai Stablecoin",
			Decimals: 18,
			PriceUSD: utils.Float64Ptr(1.0),
			MarketCap: utils.Float64Ptr(5000000000),
		},
	}, nil
}

func (r *tokenRepository) Search(ctx context.Context, query string, chainID *int) ([]*models.Token, error) {
	// TODO: Implement actual database search
	return r.GetByChainID(ctx, 1, 20, 0)
}

func (r *tokenRepository) Create(ctx context.Context, token *models.Token) (*models.Token, error) {
	// TODO: Implement actual database insert
	token.ID = uuid.New()
	return token, nil
}

func (r *tokenRepository) UpdatePrice(ctx context.Context, address string, chainID int, priceUSD, priceChange24h, marketCap float64) (*models.Token, error) {
	// TODO: Implement actual database update
	token, _ := r.GetByAddress(ctx, address, chainID)
	token.PriceUSD = &priceUSD
	token.PriceChange24h = &priceChange24h
	token.MarketCap = &marketCap
	return token, nil
}

