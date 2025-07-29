package repos

import (
	"context"
	"encoding/json"

	"github.com/defi-dashboard/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type yieldPositionRepository struct {
	db *pgxpool.Pool
}

// NewYieldPositionRepository creates a new yield position repository
func NewYieldPositionRepository(db *pgxpool.Pool) YieldPositionRepository {
	return &yieldPositionRepository{db: db}
}

func (r *yieldPositionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.YieldPosition, error) {
	query := `
		SELECT id, user_id, wallet_id, pool_id, protocol_id, position_id,
		       pool_address, chain_id, balance_raw, balance_usd, balance_tokens,
		       entry_price_usd, entry_block_number, entry_transaction_hash, entry_time,
		       is_active, last_update_block, last_update_time,
		       pending_rewards, claimed_rewards, total_rewards_usd,
		       current_value_usd, unrealized_pnl_usd, realized_pnl_usd, total_fees_paid_usd,
		       metadata, created_at, updated_at
		FROM yield_positions 
		WHERE id = $1
	`
	
	var position models.YieldPosition
	var balanceTokensJSON, pendingRewardsJSON, claimedRewardsJSON, metadataJSON []byte
	
	err := r.db.QueryRow(ctx, query, id).Scan(
		&position.ID, &position.UserID, &position.WalletID, &position.PoolID, &position.ProtocolID,
		&position.PositionID, &position.PoolAddress, &position.ChainID, &position.BalanceRaw,
		&position.BalanceUSD, &balanceTokensJSON, &position.EntryPriceUSD,
		&position.EntryBlockNumber, &position.EntryTransactionHash, &position.EntryTime,
		&position.IsActive, &position.LastUpdateBlock, &position.LastUpdateTime,
		&pendingRewardsJSON, &claimedRewardsJSON, &position.TotalRewardsUSD,
		&position.CurrentValueUSD, &position.UnrealizedPnLUSD, &position.RealizedPnLUSD,
		&position.TotalFeesPaidUSD, &metadataJSON, &position.CreatedAt, &position.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Parse JSON fields
	if balanceTokensJSON != nil {
		if err := json.Unmarshal(balanceTokensJSON, &position.BalanceTokens); err != nil {
			return nil, err
		}
	}
	if pendingRewardsJSON != nil {
		if err := json.Unmarshal(pendingRewardsJSON, &position.PendingRewards); err != nil {
			return nil, err
		}
	}
	if claimedRewardsJSON != nil {
		if err := json.Unmarshal(claimedRewardsJSON, &position.ClaimedRewards); err != nil {
			return nil, err
		}
	}
	if metadataJSON != nil {
		if err := json.Unmarshal(metadataJSON, &position.Metadata); err != nil {
			return nil, err
		}
	}

	// Calculate P&L percentage if possible
	if position.EntryPriceUSD != nil && *position.EntryPriceUSD > 0 && position.UnrealizedPnLUSD != nil {
		pnlPercentage := (*position.UnrealizedPnLUSD / *position.EntryPriceUSD) * 100
		position.PnLPercentage = &pnlPercentage
	}

	return &position, nil
}

func (r *yieldPositionRepository) GetByUser(ctx context.Context, userID uuid.UUID, filters PositionFilters) ([]*models.YieldPosition, error) {
	query := `
		SELECT yp.id, yp.user_id, yp.wallet_id, yp.pool_id, yp.protocol_id, yp.position_id,
		       yp.pool_address, yp.chain_id, yp.balance_raw, yp.balance_usd, yp.balance_tokens,
		       yp.entry_price_usd, yp.entry_block_number, yp.entry_transaction_hash, yp.entry_time,
		       yp.is_active, yp.last_update_block, yp.last_update_time,
		       yp.pending_rewards, yp.claimed_rewards, yp.total_rewards_usd,
		       yp.current_value_usd, yp.unrealized_pnl_usd, yp.realized_pnl_usd, yp.total_fees_paid_usd,
		       yp.metadata, yp.created_at, yp.updated_at,
		       pools.pool_name, pools.protocol_id as pool_protocol_id, pools.apy as pool_apy, pools.tvl_usd as pool_tvl_usd,
		       protocols.name as protocol_name, protocols.logo_uri as protocol_logo_uri
		FROM yield_positions yp
		LEFT JOIN yield_pools pools ON yp.pool_id = pools.id
		LEFT JOIN protocols ON yp.protocol_id = protocols.id
		WHERE yp.user_id = $1
		  AND ($2::boolean IS NULL OR yp.is_active = $2)
		  AND ($3::integer IS NULL OR yp.chain_id = $3)
		ORDER BY yp.current_value_usd DESC NULLS LAST, yp.created_at DESC
		LIMIT $4 OFFSET $5
	`
	
	rows, err := r.db.Query(ctx, query, userID, filters.IsActive, filters.ChainID, filters.Limit, filters.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var positions []*models.YieldPosition
	for rows.Next() {
		var position models.YieldPosition
		var balanceTokensJSON, pendingRewardsJSON, claimedRewardsJSON, metadataJSON []byte
		var poolName *string
		var poolProtocolID *uuid.UUID
		var poolAPY, poolTVL *float64
		var protocolName, protocolLogoURI *string
		
		err := rows.Scan(
			&position.ID, &position.UserID, &position.WalletID, &position.PoolID, &position.ProtocolID,
			&position.PositionID, &position.PoolAddress, &position.ChainID, &position.BalanceRaw,
			&position.BalanceUSD, &balanceTokensJSON, &position.EntryPriceUSD,
			&position.EntryBlockNumber, &position.EntryTransactionHash, &position.EntryTime,
			&position.IsActive, &position.LastUpdateBlock, &position.LastUpdateTime,
			&pendingRewardsJSON, &claimedRewardsJSON, &position.TotalRewardsUSD,
			&position.CurrentValueUSD, &position.UnrealizedPnLUSD, &position.RealizedPnLUSD,
			&position.TotalFeesPaidUSD, &metadataJSON, &position.CreatedAt, &position.UpdatedAt,
			&poolName, &poolProtocolID, &poolAPY, &poolTVL, &protocolName, &protocolLogoURI,
		)
		if err != nil {
			return nil, err
		}

		// Parse JSON fields
		if balanceTokensJSON != nil {
			json.Unmarshal(balanceTokensJSON, &position.BalanceTokens)
		}
		if pendingRewardsJSON != nil {
			json.Unmarshal(pendingRewardsJSON, &position.PendingRewards)
		}
		if claimedRewardsJSON != nil {
			json.Unmarshal(claimedRewardsJSON, &position.ClaimedRewards)
		}
		if metadataJSON != nil {
			json.Unmarshal(metadataJSON, &position.Metadata)
		}

		// Set pool and protocol information if available
		if poolName != nil {
			position.Pool = &models.YieldPool{
				PoolName: *poolName,
				APY:      poolAPY,
				TVLUSD:   poolTVL,
			}
			if poolProtocolID != nil {
				position.Pool.ProtocolID = poolProtocolID
			}
		}
		if protocolName != nil {
			position.Protocol = &models.Protocol{
				Name:    *protocolName,
				LogoURI: protocolLogoURI,
			}
		}

		// Calculate P&L percentage
		if position.EntryPriceUSD != nil && *position.EntryPriceUSD > 0 && position.UnrealizedPnLUSD != nil {
			pnlPercentage := (*position.UnrealizedPnLUSD / *position.EntryPriceUSD) * 100
			position.PnLPercentage = &pnlPercentage
		}

		positions = append(positions, &position)
	}

	return positions, nil
}

func (r *yieldPositionRepository) GetByWallet(ctx context.Context, walletID uuid.UUID, activeOnly bool) ([]*models.YieldPosition, error) {
	query := `
		SELECT yp.id, yp.user_id, yp.wallet_id, yp.pool_id, yp.protocol_id, yp.position_id,
		       yp.pool_address, yp.chain_id, yp.balance_raw, yp.balance_usd, yp.balance_tokens,
		       yp.entry_price_usd, yp.entry_block_number, yp.entry_transaction_hash, yp.entry_time,
		       yp.is_active, yp.last_update_block, yp.last_update_time,
		       yp.pending_rewards, yp.claimed_rewards, yp.total_rewards_usd,
		       yp.current_value_usd, yp.unrealized_pnl_usd, yp.realized_pnl_usd, yp.total_fees_paid_usd,
		       yp.metadata, yp.created_at, yp.updated_at,
		       pools.pool_name, pools.apy as pool_apy, protocols.name as protocol_name
		FROM yield_positions yp
		LEFT JOIN yield_pools pools ON yp.pool_id = pools.id
		LEFT JOIN protocols ON yp.protocol_id = protocols.id
		WHERE yp.wallet_id = $1
		  AND ($2::boolean IS NULL OR yp.is_active = $2)
		ORDER BY yp.current_value_usd DESC NULLS LAST, yp.created_at DESC
	`
	
	var activeOnlyPtr *bool
	if activeOnly {
		activeOnlyPtr = &activeOnly
	}
	
	rows, err := r.db.Query(ctx, query, walletID, activeOnlyPtr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanPositionsFromRows(rows)
}

func (r *yieldPositionRepository) GetByPool(ctx context.Context, poolID uuid.UUID, activeOnly bool) ([]*models.YieldPosition, error) {
	query := `
		SELECT id, user_id, wallet_id, pool_id, protocol_id, position_id,
		       pool_address, chain_id, balance_raw, balance_usd, balance_tokens,
		       entry_price_usd, entry_block_number, entry_transaction_hash, entry_time,
		       is_active, last_update_block, last_update_time,
		       pending_rewards, claimed_rewards, total_rewards_usd,
		       current_value_usd, unrealized_pnl_usd, realized_pnl_usd, total_fees_paid_usd,
		       metadata, created_at, updated_at
		FROM yield_positions
		WHERE pool_id = $1
		  AND ($2::boolean IS NULL OR is_active = $2)
		ORDER BY current_value_usd DESC NULLS LAST
	`
	
	var activeOnlyPtr *bool
	if activeOnly {
		activeOnlyPtr = &activeOnly
	}
	
	rows, err := r.db.Query(ctx, query, poolID, activeOnlyPtr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanBasicPositionsFromRows(rows)
}

func (r *yieldPositionRepository) GetByProtocol(ctx context.Context, protocolID, userID uuid.UUID, activeOnly bool) ([]*models.YieldPosition, error) {
	query := `
		SELECT yp.id, yp.user_id, yp.wallet_id, yp.pool_id, yp.protocol_id, yp.position_id,
		       yp.pool_address, yp.chain_id, yp.balance_raw, yp.balance_usd, yp.balance_tokens,
		       yp.entry_price_usd, yp.entry_block_number, yp.entry_transaction_hash, yp.entry_time,
		       yp.is_active, yp.last_update_block, yp.last_update_time,
		       yp.pending_rewards, yp.claimed_rewards, yp.total_rewards_usd,
		       yp.current_value_usd, yp.unrealized_pnl_usd, yp.realized_pnl_usd, yp.total_fees_paid_usd,
		       yp.metadata, yp.created_at, yp.updated_at,
		       pools.pool_name, pools.apy as pool_apy
		FROM yield_positions yp
		LEFT JOIN yield_pools pools ON yp.pool_id = pools.id
		WHERE yp.protocol_id = $1
		  AND yp.user_id = $2
		  AND ($3::boolean IS NULL OR yp.is_active = $3)
		ORDER BY yp.current_value_usd DESC NULLS LAST
	`
	
	var activeOnlyPtr *bool
	if activeOnly {
		activeOnlyPtr = &activeOnly
	}
	
	rows, err := r.db.Query(ctx, query, protocolID, userID, activeOnlyPtr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanPositionsFromRows(rows)
}

func (r *yieldPositionRepository) GetUserSummary(ctx context.Context, userID uuid.UUID) (*models.PositionSummary, error) {
	query := `
		SELECT 
		    COALESCE(SUM(current_value_usd), 0) as total_value_usd,
		    COALESCE(SUM(unrealized_pnl_usd), 0) + COALESCE(SUM(realized_pnl_usd), 0) as total_pnl_usd,
		    COALESCE(SUM(total_rewards_usd), 0) as total_rewards_usd,
		    COUNT(*) FILTER (WHERE is_active = true) as active_positions,
		    COUNT(*) as total_positions
		FROM yield_positions 
		WHERE user_id = $1
	`
	
	var summary models.PositionSummary
	var totalPositions int
	
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&summary.TotalValueUSD, &summary.TotalPnLUSD, &summary.TotalRewardsUSD,
		&summary.ActivePositions, &totalPositions,
	)
	if err != nil {
		return nil, err
	}

	// Calculate P&L percentage
	if summary.TotalValueUSD > 0 {
		entryValue := summary.TotalValueUSD - summary.TotalPnLUSD
		if entryValue > 0 {
			summary.TotalPnLPercentage = (summary.TotalPnLUSD / entryValue) * 100
		}
	}

	return &summary, nil
}

func (r *yieldPositionRepository) GetUserPositionsWithPools(ctx context.Context, userID uuid.UUID, filters PositionFilters) ([]*models.YieldPosition, error) {
	query := `
		SELECT 
		    yp.id, yp.user_id, yp.wallet_id, yp.pool_id, yp.protocol_id, yp.position_id,
		    yp.pool_address, yp.chain_id, yp.balance_raw, yp.balance_usd, yp.balance_tokens,
		    yp.entry_price_usd, yp.entry_block_number, yp.entry_transaction_hash, yp.entry_time,
		    yp.is_active, yp.last_update_block, yp.last_update_time,
		    yp.pending_rewards, yp.claimed_rewards, yp.total_rewards_usd,
		    yp.current_value_usd, yp.unrealized_pnl_usd, yp.realized_pnl_usd, yp.total_fees_paid_usd,
		    yp.metadata, yp.created_at, yp.updated_at,
		    pools.pool_name, pools.pool_id as pool_identifier, pools.protocol_id as pool_protocol_id,
		    pools.apy as pool_apy, pools.tvl_usd as pool_tvl_usd, pools.risk_level as pool_risk_level,
		    protocols.name as protocol_name, protocols.slug as protocol_slug,
		    protocols.logo_uri as protocol_logo_uri, protocols.category as protocol_category
		FROM yield_positions yp
		JOIN yield_pools pools ON yp.pool_id = pools.id
		LEFT JOIN protocols ON pools.protocol_id = protocols.id
		WHERE yp.user_id = $1
		  AND ($2::boolean IS NULL OR yp.is_active = $2)
		  AND ($3::integer IS NULL OR yp.chain_id = $3)
		ORDER BY yp.current_value_usd DESC NULLS LAST, yp.entry_time DESC
	`
	
	rows, err := r.db.Query(ctx, query, userID, filters.IsActive, filters.ChainID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanEnhancedPositionsFromRows(rows)
}

func (r *yieldPositionRepository) GetTopByValue(ctx context.Context, limit int) ([]*models.YieldPosition, error) {
	query := `
		SELECT yp.id, yp.user_id, yp.wallet_id, yp.pool_id, yp.protocol_id, yp.position_id,
		       yp.pool_address, yp.chain_id, yp.balance_raw, yp.balance_usd, yp.balance_tokens,
		       yp.entry_price_usd, yp.entry_block_number, yp.entry_transaction_hash, yp.entry_time,
		       yp.is_active, yp.last_update_block, yp.last_update_time,
		       yp.pending_rewards, yp.claimed_rewards, yp.total_rewards_usd,
		       yp.current_value_usd, yp.unrealized_pnl_usd, yp.realized_pnl_usd, yp.total_fees_paid_usd,
		       yp.metadata, yp.created_at, yp.updated_at,
		       u.address as user_address, pools.pool_name, protocols.name as protocol_name
		FROM yield_positions yp
		JOIN users u ON yp.user_id = u.id
		LEFT JOIN yield_pools pools ON yp.pool_id = pools.id
		LEFT JOIN protocols ON yp.protocol_id = protocols.id
		WHERE yp.is_active = true
		  AND yp.current_value_usd IS NOT NULL
		ORDER BY yp.current_value_usd DESC
		LIMIT $1
	`
	
	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanPositionsFromRows(rows)
}

func (r *yieldPositionRepository) Create(ctx context.Context, position *models.YieldPosition) (*models.YieldPosition, error) {
	// Serialize JSON fields
	balanceTokensJSON, _ := json.Marshal(position.BalanceTokens)
	metadataJSON, _ := json.Marshal(position.Metadata)

	query := `
		INSERT INTO yield_positions (
		    user_id, wallet_id, pool_id, protocol_id, position_id,
		    pool_address, chain_id, balance_raw, balance_usd, balance_tokens,
		    entry_price_usd, entry_block_number, entry_transaction_hash, entry_time,
		    current_value_usd, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING id, created_at, updated_at
	`
	
	err := r.db.QueryRow(ctx, query,
		position.UserID, position.WalletID, position.PoolID, position.ProtocolID,
		position.PositionID, position.PoolAddress, position.ChainID,
		position.BalanceRaw, position.BalanceUSD, balanceTokensJSON,
		position.EntryPriceUSD, position.EntryBlockNumber, position.EntryTransactionHash,
		position.EntryTime, position.CurrentValueUSD, metadataJSON,
	).Scan(&position.ID, &position.CreatedAt, &position.UpdatedAt)
	
	return position, err
}

func (r *yieldPositionRepository) Update(ctx context.Context, position *models.YieldPosition) (*models.YieldPosition, error) {
	// Serialize JSON fields
	balanceTokensJSON, _ := json.Marshal(position.BalanceTokens)
	pendingRewardsJSON, _ := json.Marshal(position.PendingRewards)
	metadataJSON, _ := json.Marshal(position.Metadata)

	query := `
		UPDATE yield_positions 
		SET balance_raw = $2,
		    balance_usd = $3,
		    balance_tokens = $4,
		    current_value_usd = $5,
		    pending_rewards = $6,
		    total_rewards_usd = $7,
		    last_update_block = $8,
		    last_update_time = $9,
		    metadata = $10,
		    updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`
	
	err := r.db.QueryRow(ctx, query,
		position.ID, position.BalanceRaw, position.BalanceUSD, balanceTokensJSON,
		position.CurrentValueUSD, pendingRewardsJSON, position.TotalRewardsUSD,
		position.LastUpdateBlock, position.LastUpdateTime, metadataJSON,
	).Scan(&position.UpdatedAt)
	
	return position, err
}

func (r *yieldPositionRepository) UpdateBalance(ctx context.Context, id uuid.UUID, balanceRaw string, balanceUSD, currentValueUSD float64) error {
	query := `
		UPDATE yield_positions 
		SET balance_raw = $2,
		    balance_usd = $3,
		    current_value_usd = $4,
		    last_update_time = NOW(),
		    updated_at = NOW()
		WHERE id = $1
	`
	
	_, err := r.db.Exec(ctx, query, id, balanceRaw, balanceUSD, currentValueUSD)
	return err
}

func (r *yieldPositionRepository) UpdateRewards(ctx context.Context, id uuid.UUID, pendingRewards, claimedRewards interface{}, totalRewardsUSD float64) error {
	// Serialize JSON fields
	pendingRewardsJSON, _ := json.Marshal(pendingRewards)
	claimedRewardsJSON, _ := json.Marshal(claimedRewards)

	query := `
		UPDATE yield_positions 
		SET pending_rewards = $2,
		    claimed_rewards = $3,
		    total_rewards_usd = $4,
		    updated_at = NOW()
		WHERE id = $1
	`
	
	_, err := r.db.Exec(ctx, query, id, pendingRewardsJSON, claimedRewardsJSON, totalRewardsUSD)
	return err
}

func (r *yieldPositionRepository) Close(ctx context.Context, id uuid.UUID, realizedPnLUSD float64) error {
	query := `
		UPDATE yield_positions 
		SET is_active = false,
		    realized_pnl_usd = $2,
		    updated_at = NOW()
		WHERE id = $1
	`
	
	_, err := r.db.Exec(ctx, query, id, realizedPnLUSD)
	return err
}

func (r *yieldPositionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM yield_positions WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *yieldPositionRepository) UpdateAllPnL(ctx context.Context) error {
	query := `
		UPDATE yield_positions 
		SET unrealized_pnl_usd = current_value_usd - entry_price_usd,
		    updated_at = NOW()
		WHERE is_active = true 
		  AND current_value_usd IS NOT NULL 
		  AND entry_price_usd IS NOT NULL
	`
	
	_, err := r.db.Exec(ctx, query)
	return err
}

// Helper methods for scanning rows

func (r *yieldPositionRepository) scanBasicPositionsFromRows(rows interface{}) ([]*models.YieldPosition, error) {
	// Implement basic position scanning logic
	// This is a simplified version - you'd need to implement the full scanning logic
	return []*models.YieldPosition{}, nil
}

func (r *yieldPositionRepository) scanPositionsFromRows(rows interface{}) ([]*models.YieldPosition, error) {
	// Implement position scanning with basic pool/protocol info
	// This is a simplified version - you'd need to implement the full scanning logic
	return []*models.YieldPosition{}, nil
}

func (r *yieldPositionRepository) scanEnhancedPositionsFromRows(rows interface{}) ([]*models.YieldPosition, error) {
	// Implement enhanced position scanning with full pool/protocol info
	// This is a simplified version - you'd need to implement the full scanning logic
	return []*models.YieldPosition{}, nil
}