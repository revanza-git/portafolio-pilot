package jobs

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/defi-dashboard/backend/pkg/external"
	"github.com/defi-dashboard/backend/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PriceRefreshJob struct {
	db              *pgxpool.Pool
	coinGeckoClient *external.CoinGeckoClient
	defiLlamaClient *external.DefiLlamaClient
}

func NewPriceRefreshJob(db *pgxpool.Pool, cgClient *external.CoinGeckoClient, dlClient *external.DefiLlamaClient) *PriceRefreshJob {
	return &PriceRefreshJob{
		db:              db,
		coinGeckoClient: cgClient,
		defiLlamaClient: dlClient,
	}
}

// Run executes the price refresh job
func (j *PriceRefreshJob) Run(ctx context.Context) error {
	logger.Info("Starting price refresh job")

	// Run price updates and yield updates concurrently
	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	// Update token prices
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := j.updateTokenPrices(ctx); err != nil {
			errChan <- fmt.Errorf("token price update failed: %w", err)
		}
	}()

	// Update yield pool APRs
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := j.updateYieldPools(ctx); err != nil {
			errChan <- fmt.Errorf("yield pool update failed: %w", err)
		}
	}()

	// Wait for all updates to complete
	wg.Wait()
	close(errChan)

	// Check for errors
	var errors []string
	for err := range errChan {
		errors = append(errors, err.Error())
		logger.Error("Job error", "error", err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("price refresh job completed with errors: %s", strings.Join(errors, "; "))
	}

	logger.Info("Price refresh job completed successfully")
	return nil
}

// updateTokenPrices fetches and updates token prices from CoinGecko
func (j *PriceRefreshJob) updateTokenPrices(ctx context.Context) error {
	// Get list of tokens to update from database
	tokens, err := j.getActiveTokens(ctx)
	if err != nil {
		return fmt.Errorf("failed to get active tokens: %w", err)
	}

	if len(tokens) == 0 {
		logger.Warn("No active tokens found to update")
		return nil
	}

	// Map tokens to CoinGecko IDs
	var tokenIDs []string
	tokenMap := make(map[string]*TokenInfo)
	
	for _, token := range tokens {
		cgID := j.getCoinGeckoID(token.Symbol)
		if cgID != "" {
			tokenIDs = append(tokenIDs, cgID)
			tokenMap[cgID] = token
		}
	}

	if len(tokenIDs) == 0 {
		logger.Warn("No valid CoinGecko IDs found")
		return nil
	}

	// Fetch prices from CoinGecko with retry
	var prices external.PriceResponse
	for i := 0; i < 3; i++ {
		prices, err = j.coinGeckoClient.GetTokenPrices(ctx, tokenIDs)
		if err == nil {
			break
		}
		if i < 2 {
			logger.Warn("CoinGecko API call failed, retrying", 
				"attempt", i+1, 
				"error", err)
			time.Sleep(time.Duration(i+1) * time.Second)
		}
	}

	if err != nil {
		return fmt.Errorf("failed to fetch prices after retries: %w", err)
	}

	// Update prices in database
	tx, err := j.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	updated := 0
	for cgID, priceData := range prices {
		token, exists := tokenMap[cgID]
		if !exists {
			continue
		}

		// Update token price
		_, err = tx.Exec(ctx, `
			UPDATE tokens 
			SET price_usd = $1, 
				price_change_24h = $2,
				last_updated = NOW(),
				updated_at = NOW()
			WHERE address = $3 AND chain_id = $4`,
			priceData.USD, priceData.USD24hChange, token.Address, token.ChainID)
		
		if err != nil {
			logger.Error("Failed to update token price",
				"token", token.Symbol,
				"error", err)
			continue
		}

		// Insert price history record
		_, err = tx.Exec(ctx, `
			INSERT INTO price_history (token_id, price_usd, timestamp, source)
			VALUES ($1, $2, NOW(), 'coingecko')`,
			token.ID, priceData.USD)
		
		if err != nil {
			logger.Error("Failed to insert price history",
				"token", token.Symbol,
				"error", err)
		}

		updated++
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	logger.Info("Token prices updated", 
		"total", len(tokens),
		"updated", updated)

	return nil
}

// updateYieldPools fetches and updates yield pool data from DefiLlama
func (j *PriceRefreshJob) updateYieldPools(ctx context.Context) error {
	// Fetch yield pools from DefiLlama with retry
	var pools []external.YieldPool
	var err error
	
	for i := 0; i < 3; i++ {
		pools, err = j.defiLlamaClient.GetYieldPools(ctx)
		if err == nil {
			break
		}
		if i < 2 {
			logger.Warn("DefiLlama API call failed, retrying",
				"attempt", i+1,
				"error", err)
			time.Sleep(time.Duration(i+1) * time.Second)
		}
	}

	if err != nil {
		return fmt.Errorf("failed to fetch yield pools after retries: %w", err)
	}

	// Filter pools we're interested in (top protocols, supported chains)
	supportedChains := map[string]bool{
		"Ethereum": true,
		"Polygon":  true,
		"Arbitrum": true,
		"Optimism": true,
		"Base":     true,
	}

	supportedProtocols := map[string]bool{
		"aave-v3":     true,
		"aave-v2":     true,
		"compound-v3": true,
		"compound-v2": true,
		"uniswap-v3":  true,
		"curve":       true,
		"balancer-v2": true,
		"yearn":       true,
		"convex":      true,
		"stargate":    true,
	}

	// Begin transaction
	tx, err := j.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	updated := 0
	for _, pool := range pools {
		// Skip unsupported chains or protocols
		if !supportedChains[pool.Chain] || !supportedProtocols[strings.ToLower(pool.Project)] {
			continue
		}

		// Skip pools with very low TVL
		if pool.TVL < 100000 { // $100k minimum
			continue
		}

		// Upsert yield pool data
		_, err = tx.Exec(ctx, `
			INSERT INTO yield_pools (
				pool_id, protocol, pool_name, chain, symbol,
				tvl_usd, apy, apy_base, apy_reward,
				il_7d, stable_coin, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW())
			ON CONFLICT (pool_id) DO UPDATE SET
				tvl_usd = $6,
				apy = $7,
				apy_base = $8,
				apy_reward = $9,
				il_7d = $10,
				updated_at = NOW()`,
			pool.Pool, pool.Project, pool.Symbol, pool.Chain, pool.Symbol,
			pool.TVL, pool.APY, pool.APYBase, pool.APYReward,
			pool.IL7d, pool.StableCoin)

		if err != nil {
			logger.Error("Failed to upsert yield pool",
				"pool", pool.Pool,
				"error", err)
			continue
		}

		updated++
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	logger.Info("Yield pools updated",
		"total", len(pools),
		"updated", updated)

	return nil
}

// getActiveTokens retrieves tokens that need price updates
func (j *PriceRefreshJob) getActiveTokens(ctx context.Context) ([]*TokenInfo, error) {
	rows, err := j.db.Query(ctx, `
		SELECT DISTINCT t.id, t.address, t.chain_id, t.symbol, t.name
		FROM tokens t
		INNER JOIN balances b ON b.token_id = t.id
		WHERE b.balance > 0
			OR t.last_updated IS NULL
			OR t.last_updated < NOW() - INTERVAL '15 minutes'
		ORDER BY t.market_cap DESC NULLS LAST
		LIMIT 100`)
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []*TokenInfo
	for rows.Next() {
		var token TokenInfo
		err := rows.Scan(&token.ID, &token.Address, &token.ChainID, &token.Symbol, &token.Name)
		if err != nil {
			logger.Error("Failed to scan token", "error", err)
			continue
		}
		tokens = append(tokens, &token)
	}

	return tokens, rows.Err()
}

// getCoinGeckoID maps token symbols to CoinGecko IDs
func (j *PriceRefreshJob) getCoinGeckoID(symbol string) string {
	// Use the mapping from external package
	return external.TokenIDMappings[strings.ToLower(symbol)]
}

// TokenInfo represents basic token information
type TokenInfo struct {
	ID      string
	Address string
	ChainID int
	Symbol  string
	Name    string
}