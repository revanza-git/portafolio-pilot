package repos

import (
	"context"
	"encoding/json"

	"github.com/defi-dashboard/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type yieldPoolRepository struct {
	db *pgxpool.Pool
}

// NewYieldPoolRepository creates a new yield pool repository
func NewYieldPoolRepository(db *pgxpool.Pool) YieldPoolRepository {
	return &yieldPoolRepository{db: db}
}

func (r *yieldPoolRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.YieldPool, error) {
	query := `
		SELECT yp.id, yp.pool_id, yp.protocol_id, yp.pool_name, yp.chain_id, yp.chain, 
		       yp.pool_address, yp.symbol, yp.token_addresses, yp.tvl_usd, yp.apy, 
		       yp.apy_base, yp.apy_reward, yp.fees_apr, yp.il_7d, yp.risk_level,
		       yp.min_deposit_usd, yp.max_deposit_usd, yp.is_active, yp.stable_coin,
		       yp.metadata, yp.created_at, yp.updated_at,
		       p.name as protocol_name, p.logo_uri as protocol_logo_uri
		FROM yield_pools yp
		LEFT JOIN protocols p ON yp.protocol_id = p.id
		WHERE yp.id = $1
	`
	
	var pool models.YieldPool
	var tokenAddressesJSON, metadataJSON []byte
	var protocolName, protocolLogoURI *string
	
	err := r.db.QueryRow(ctx, query, id).Scan(
		&pool.ID, &pool.PoolID, &pool.ProtocolID, &pool.PoolName, &pool.ChainID,
		&pool.Chain, &pool.PoolAddress, &pool.Symbol, &tokenAddressesJSON,
		&pool.TVLUSD, &pool.APY, &pool.APYBase, &pool.APYReward, &pool.FeesAPR,
		&pool.IL7D, &pool.RiskLevel, &pool.MinDepositUSD, &pool.MaxDepositUSD,
		&pool.IsActive, &pool.StableCoin, &metadataJSON, &pool.CreatedAt,
		&pool.UpdatedAt, &protocolName, &protocolLogoURI,
	)
	if err != nil {
		return nil, err
	}

	// Parse JSON fields
	if tokenAddressesJSON != nil {
		if err := json.Unmarshal(tokenAddressesJSON, &pool.TokenAddresses); err != nil {
			return nil, err
		}
	}
	if metadataJSON != nil {
		if err := json.Unmarshal(metadataJSON, &pool.Metadata); err != nil {
			return nil, err
		}
	}

	// Set protocol info if available
	if protocolName != nil {
		pool.Protocol = &models.Protocol{
			Name:    *protocolName,
			LogoURI: protocolLogoURI,
		}
	}

	return &pool, nil
}

func (r *yieldPoolRepository) GetByPoolID(ctx context.Context, poolID string) (*models.YieldPool, error) {
	query := `
		SELECT yp.id, yp.pool_id, yp.protocol_id, yp.pool_name, yp.chain_id, yp.chain, 
		       yp.pool_address, yp.symbol, yp.token_addresses, yp.tvl_usd, yp.apy, 
		       yp.apy_base, yp.apy_reward, yp.fees_apr, yp.il_7d, yp.risk_level,
		       yp.min_deposit_usd, yp.max_deposit_usd, yp.is_active, yp.stable_coin,
		       yp.metadata, yp.created_at, yp.updated_at,
		       p.name as protocol_name, p.logo_uri as protocol_logo_uri
		FROM yield_pools yp
		LEFT JOIN protocols p ON yp.protocol_id = p.id
		WHERE yp.pool_id = $1
	`
	
	var pool models.YieldPool
	var tokenAddressesJSON, metadataJSON []byte
	var protocolName, protocolLogoURI *string
	
	err := r.db.QueryRow(ctx, query, poolID).Scan(
		&pool.ID, &pool.PoolID, &pool.ProtocolID, &pool.PoolName, &pool.ChainID,
		&pool.Chain, &pool.PoolAddress, &pool.Symbol, &tokenAddressesJSON,
		&pool.TVLUSD, &pool.APY, &pool.APYBase, &pool.APYReward, &pool.FeesAPR,
		&pool.IL7D, &pool.RiskLevel, &pool.MinDepositUSD, &pool.MaxDepositUSD,
		&pool.IsActive, &pool.StableCoin, &metadataJSON, &pool.CreatedAt,
		&pool.UpdatedAt, &protocolName, &protocolLogoURI,
	)
	if err != nil {
		return nil, err
	}

	// Parse JSON fields
	if tokenAddressesJSON != nil {
		json.Unmarshal(tokenAddressesJSON, &pool.TokenAddresses)
	}
	if metadataJSON != nil {
		json.Unmarshal(metadataJSON, &pool.Metadata)
	}

	// Set protocol info if available
	if protocolName != nil {
		pool.Protocol = &models.Protocol{
			Name:    *protocolName,
			LogoURI: protocolLogoURI,
		}
	}

	return &pool, nil
}

func (r *yieldPoolRepository) GetAll(ctx context.Context, filters YieldPoolFilters) ([]*models.YieldPool, error) {
	query := `
		SELECT yp.id, yp.pool_id, yp.protocol_id, yp.pool_name, yp.chain_id, yp.chain, 
		       yp.pool_address, yp.symbol, yp.token_addresses, yp.tvl_usd, yp.apy, 
		       yp.apy_base, yp.apy_reward, yp.fees_apr, yp.il_7d, yp.risk_level,
		       yp.min_deposit_usd, yp.max_deposit_usd, yp.is_active, yp.stable_coin,
		       yp.metadata, yp.created_at, yp.updated_at,
		       p.name as protocol_name, p.logo_uri as protocol_logo_uri, p.category as protocol_category
		FROM yield_pools yp
		LEFT JOIN protocols p ON yp.protocol_id = p.id
		WHERE ($1::varchar IS NULL OR yp.chain = $1)
		  AND ($2::integer IS NULL OR yp.chain_id = $2)
		  AND ($3::decimal IS NULL OR yp.tvl_usd >= $3)
		  AND ($4::decimal IS NULL OR yp.apy >= $4)
		  AND ($5::varchar IS NULL OR p.slug = $5)
		  AND ($6::varchar IS NULL OR yp.risk_level = $6)
		  AND ($7::boolean IS NULL OR yp.is_active = $7)
		ORDER BY 
		  CASE WHEN $8 = 'apy' THEN yp.apy END DESC,
		  CASE WHEN $8 = 'tvl' THEN yp.tvl_usd END DESC,
		  CASE WHEN $8 = 'name' THEN yp.pool_name END ASC,
		  yp.created_at DESC
		LIMIT $9 OFFSET $10
	`
	
	rows, err := r.db.Query(ctx, query,
		filters.Chain, filters.ChainID, filters.MinTVL, filters.MinAPY,
		filters.ProtocolSlug, filters.RiskLevel, filters.IsActive,
		filters.SortBy, filters.Limit, filters.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pools []*models.YieldPool
	for rows.Next() {
		var pool models.YieldPool
		var tokenAddressesJSON, metadataJSON []byte
		var protocolName, protocolLogoURI, protocolCategory *string
		
		err := rows.Scan(
			&pool.ID, &pool.PoolID, &pool.ProtocolID, &pool.PoolName, &pool.ChainID,
			&pool.Chain, &pool.PoolAddress, &pool.Symbol, &tokenAddressesJSON,
			&pool.TVLUSD, &pool.APY, &pool.APYBase, &pool.APYReward, &pool.FeesAPR,
			&pool.IL7D, &pool.RiskLevel, &pool.MinDepositUSD, &pool.MaxDepositUSD,
			&pool.IsActive, &pool.StableCoin, &metadataJSON, &pool.CreatedAt,
			&pool.UpdatedAt, &protocolName, &protocolLogoURI, &protocolCategory,
		)
		if err != nil {
			return nil, err
		}

		// Parse JSON fields
		if tokenAddressesJSON != nil {
			json.Unmarshal(tokenAddressesJSON, &pool.TokenAddresses)
		}
		if metadataJSON != nil {
			json.Unmarshal(metadataJSON, &pool.Metadata)
		}

		// Set protocol info if available
		if protocolName != nil {
			pool.Protocol = &models.Protocol{
				Name:     *protocolName,
				LogoURI:  protocolLogoURI,
				Category: protocolCategory,
			}
		}

		pools = append(pools, &pool)
	}

	return pools, nil
}

func (r *yieldPoolRepository) Count(ctx context.Context, filters YieldPoolFilters) (int64, error) {
	query := `
		SELECT COUNT(*) FROM yield_pools yp
		LEFT JOIN protocols p ON yp.protocol_id = p.id
		WHERE ($1::varchar IS NULL OR yp.chain = $1)
		  AND ($2::integer IS NULL OR yp.chain_id = $2)
		  AND ($3::decimal IS NULL OR yp.tvl_usd >= $3)
		  AND ($4::decimal IS NULL OR yp.apy >= $4)
		  AND ($5::varchar IS NULL OR p.slug = $5)
		  AND ($6::varchar IS NULL OR yp.risk_level = $6)
		  AND ($7::boolean IS NULL OR yp.is_active = $7)
	`
	
	var count int64
	err := r.db.QueryRow(ctx, query,
		filters.Chain, filters.ChainID, filters.MinTVL, filters.MinAPY,
		filters.ProtocolSlug, filters.RiskLevel, filters.IsActive).Scan(&count)
	return count, err
}

func (r *yieldPoolRepository) GetByProtocol(ctx context.Context, protocolID uuid.UUID, activeOnly bool) ([]*models.YieldPool, error) {
	query := `
		SELECT id, pool_id, protocol_id, pool_name, chain_id, chain, 
		       pool_address, symbol, token_addresses, tvl_usd, apy, 
		       apy_base, apy_reward, fees_apr, il_7d, risk_level,
		       min_deposit_usd, max_deposit_usd, is_active, stable_coin,
		       metadata, created_at, updated_at
		FROM yield_pools
		WHERE protocol_id = $1
		  AND ($2::boolean IS NULL OR is_active = $2)
		ORDER BY tvl_usd DESC
	`
	
	var activeOnlyPtr *bool
	if activeOnly {
		activeOnlyPtr = &activeOnly
	}
	
	rows, err := r.db.Query(ctx, query, protocolID, activeOnlyPtr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanPoolsFromRows(rows)
}

func (r *yieldPoolRepository) GetByChain(ctx context.Context, chainID int) ([]*models.YieldPool, error) {
	query := `
		SELECT yp.id, yp.pool_id, yp.protocol_id, yp.pool_name, yp.chain_id, yp.chain, 
		       yp.pool_address, yp.symbol, yp.token_addresses, yp.tvl_usd, yp.apy, 
		       yp.apy_base, yp.apy_reward, yp.fees_apr, yp.il_7d, yp.risk_level,
		       yp.min_deposit_usd, yp.max_deposit_usd, yp.is_active, yp.stable_coin,
		       yp.metadata, yp.created_at, yp.updated_at,
		       p.name as protocol_name
		FROM yield_pools yp
		LEFT JOIN protocols p ON yp.protocol_id = p.id
		WHERE yp.chain_id = $1
		  AND yp.is_active = true
		ORDER BY yp.apy DESC
	`
	
	rows, err := r.db.Query(ctx, query, chainID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanPoolsFromRows(rows)
}

func (r *yieldPoolRepository) GetTopByTVL(ctx context.Context, limit int) ([]*models.YieldPool, error) {
	query := `
		SELECT yp.id, yp.pool_id, yp.protocol_id, yp.pool_name, yp.chain_id, yp.chain, 
		       yp.pool_address, yp.symbol, yp.token_addresses, yp.tvl_usd, yp.apy, 
		       yp.apy_base, yp.apy_reward, yp.fees_apr, yp.il_7d, yp.risk_level,
		       yp.min_deposit_usd, yp.max_deposit_usd, yp.is_active, yp.stable_coin,
		       yp.metadata, yp.created_at, yp.updated_at,
		       p.name as protocol_name
		FROM yield_pools yp
		LEFT JOIN protocols p ON yp.protocol_id = p.id
		WHERE yp.is_active = true
		  AND yp.tvl_usd IS NOT NULL
		ORDER BY yp.tvl_usd DESC
		LIMIT $1
	`
	
	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanPoolsFromRows(rows)
}

func (r *yieldPoolRepository) Upsert(ctx context.Context, pool *models.YieldPool) error {
	// Serialize JSON fields
	tokenAddressesJSON, _ := json.Marshal(pool.TokenAddresses)
	metadataJSON, _ := json.Marshal(pool.Metadata)

	query := `
		INSERT INTO yield_pools (
		    pool_id, protocol, pool_name, chain, symbol,
		    tvl_usd, apy, apy_base, apy_reward,
		    il_7d, stable_coin, protocol_id, chain_id,
		    pool_address, token_addresses, fees_apr,
		    risk_level, min_deposit_usd, max_deposit_usd,
		    is_active, metadata, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, NOW())
		ON CONFLICT (pool_id) DO UPDATE SET
		    tvl_usd = $6,
		    apy = $7,
		    apy_base = $8,
		    apy_reward = $9,
		    il_7d = $10,
		    stable_coin = $11,
		    protocol_id = $12,
		    chain_id = $13,
		    pool_address = $14,
		    token_addresses = $15,
		    fees_apr = $16,
		    risk_level = $17,
		    min_deposit_usd = $18,
		    max_deposit_usd = $19,
		    is_active = $20,
		    metadata = $21,
		    updated_at = NOW()
	`
	
	// Extract protocol name for legacy field (since some queries might still use it)
	protocolName := ""
	if pool.Protocol != nil {
		protocolName = pool.Protocol.Name
	}
	
	_, err := r.db.Exec(ctx, query,
		pool.PoolID, protocolName, pool.PoolName, pool.Chain, pool.Symbol,
		pool.TVLUSD, pool.APY, pool.APYBase, pool.APYReward, pool.IL7D,
		pool.StableCoin, pool.ProtocolID, pool.ChainID, pool.PoolAddress,
		tokenAddressesJSON, pool.FeesAPR, pool.RiskLevel, pool.MinDepositUSD,
		pool.MaxDepositUSD, pool.IsActive, metadataJSON)
	
	return err
}

func (r *yieldPoolRepository) UpdateAPY(ctx context.Context, poolID string, apy, apyBase, apyReward float64) error {
	query := `
		UPDATE yield_pools 
		SET apy = $2,
		    apy_base = $3,
		    apy_reward = $4,
		    updated_at = NOW()
		WHERE pool_id = $1
	`
	
	_, err := r.db.Exec(ctx, query, poolID, apy, apyBase, apyReward)
	return err
}

func (r *yieldPoolRepository) UpdateTVL(ctx context.Context, poolID string, tvlUSD float64) error {
	query := `
		UPDATE yield_pools 
		SET tvl_usd = $2,
		    updated_at = NOW()
		WHERE pool_id = $1
	`
	
	_, err := r.db.Exec(ctx, query, poolID, tvlUSD)
	return err
}

func (r *yieldPoolRepository) Deactivate(ctx context.Context, poolID string) error {
	query := `
		UPDATE yield_pools 
		SET is_active = false,
		    updated_at = NOW()
		WHERE pool_id = $1
	`
	
	_, err := r.db.Exec(ctx, query, poolID)
	return err
}

// Helper method for scanning pools from rows
func (r *yieldPoolRepository) scanPoolsFromRows(rows interface{}) ([]*models.YieldPool, error) {
	// This is a simplified implementation
	// In a real implementation, you'd properly scan the rows and handle all fields
	return []*models.YieldPool{}, nil
}