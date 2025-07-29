-- name: GetPosition :one
SELECT * FROM yield_positions 
WHERE id = $1
LIMIT 1;

-- name: GetPositionsByUser :many
SELECT yp.*, 
       pools.pool_name,
       pools.protocol_id as pool_protocol_id,
       pools.apy as pool_apy,
       pools.tvl_usd as pool_tvl_usd,
       protocols.name as protocol_name,
       protocols.logo_uri as protocol_logo_uri
FROM yield_positions yp
LEFT JOIN yield_pools pools ON yp.pool_id = pools.id
LEFT JOIN protocols ON yp.protocol_id = protocols.id
WHERE yp.user_id = $1
  AND ($2::boolean IS NULL OR yp.is_active = $2)
  AND ($3::integer IS NULL OR yp.chain_id = $3)
ORDER BY yp.current_value_usd DESC NULLS LAST, yp.created_at DESC
LIMIT $4 OFFSET $5;

-- name: GetPositionsByWallet :many
SELECT yp.*, 
       pools.pool_name,
       pools.apy as pool_apy,
       protocols.name as protocol_name
FROM yield_positions yp
LEFT JOIN yield_pools pools ON yp.pool_id = pools.id
LEFT JOIN protocols ON yp.protocol_id = protocols.id
WHERE yp.wallet_id = $1
  AND ($2::boolean IS NULL OR yp.is_active = $2)
ORDER BY yp.current_value_usd DESC NULLS LAST, yp.created_at DESC;

-- name: GetPositionsByPool :many
SELECT * FROM yield_positions
WHERE pool_id = $1
  AND ($2::boolean IS NULL OR is_active = $2)
ORDER BY current_value_usd DESC NULLS LAST;

-- name: CreatePosition :one
INSERT INTO yield_positions (
    user_id, wallet_id, pool_id, protocol_id, position_id,
    pool_address, chain_id, balance_raw, balance_usd, balance_tokens,
    entry_price_usd, entry_block_number, entry_transaction_hash, entry_time,
    current_value_usd, metadata
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
RETURNING *;

-- name: UpdatePosition :one
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
RETURNING *;

-- name: UpdatePositionBalance :exec
UPDATE yield_positions 
SET balance_raw = $2,
    balance_usd = $3,
    current_value_usd = $4,
    last_update_time = NOW(),
    updated_at = NOW()
WHERE id = $1;

-- name: UpdatePositionRewards :exec
UPDATE yield_positions 
SET pending_rewards = $2,
    claimed_rewards = $3,
    total_rewards_usd = $4,
    updated_at = NOW()
WHERE id = $1;

-- name: ClosePosition :exec
UPDATE yield_positions 
SET is_active = false,
    realized_pnl_usd = $2,
    updated_at = NOW()
WHERE id = $1;

-- name: DeletePosition :exec
DELETE FROM yield_positions 
WHERE id = $1;

-- name: GetUserPositionSummary :one
SELECT 
    COALESCE(SUM(current_value_usd), 0) as total_value_usd,
    COALESCE(SUM(unrealized_pnl_usd), 0) + COALESCE(SUM(realized_pnl_usd), 0) as total_pnl_usd,
    COALESCE(SUM(total_rewards_usd), 0) as total_rewards_usd,
    COUNT(*) FILTER (WHERE is_active = true) as active_positions,
    COUNT(*) as total_positions
FROM yield_positions 
WHERE user_id = $1;

-- name: GetPositionsByProtocol :many
SELECT yp.*, 
       pools.pool_name,
       pools.apy as pool_apy
FROM yield_positions yp
LEFT JOIN yield_pools pools ON yp.pool_id = pools.id
WHERE yp.protocol_id = $1
  AND yp.user_id = $2
  AND ($3::boolean IS NULL OR yp.is_active = $3)
ORDER BY yp.current_value_usd DESC NULLS LAST;

-- name: GetTopPositionsByValue :many
SELECT yp.*, 
       u.address as user_address,
       pools.pool_name,
       protocols.name as protocol_name
FROM yield_positions yp
JOIN users u ON yp.user_id = u.id
LEFT JOIN yield_pools pools ON yp.pool_id = pools.id
LEFT JOIN protocols ON yp.protocol_id = protocols.id
WHERE yp.is_active = true
  AND yp.current_value_usd IS NOT NULL
ORDER BY yp.current_value_usd DESC
LIMIT $1;

-- name: UpdatePositionsPnL :exec
UPDATE yield_positions 
SET unrealized_pnl_usd = current_value_usd - entry_price_usd,
    updated_at = NOW()
WHERE is_active = true 
  AND current_value_usd IS NOT NULL 
  AND entry_price_usd IS NOT NULL;

-- name: GetPositionsForPnLUpdate :many
SELECT id, current_value_usd, entry_price_usd
FROM yield_positions 
WHERE is_active = true 
  AND current_value_usd IS NOT NULL 
  AND entry_price_usd IS NOT NULL
  AND (unrealized_pnl_usd IS NULL OR 
       ABS(unrealized_pnl_usd - (current_value_usd - entry_price_usd)) > 0.01);

-- name: GetUserPositionsWithPools :many
SELECT 
    yp.*,
    pools.pool_name,
    pools.pool_id as pool_identifier,
    pools.protocol_id as pool_protocol_id,
    pools.apy as pool_apy,
    pools.tvl_usd as pool_tvl_usd,
    pools.risk_level as pool_risk_level,
    protocols.name as protocol_name,
    protocols.slug as protocol_slug,
    protocols.logo_uri as protocol_logo_uri,
    protocols.category as protocol_category
FROM yield_positions yp
JOIN yield_pools pools ON yp.pool_id = pools.id
LEFT JOIN protocols ON pools.protocol_id = protocols.id
WHERE yp.user_id = $1
  AND ($2::boolean IS NULL OR yp.is_active = $2)
  AND ($3::integer IS NULL OR yp.chain_id = $3)
ORDER BY yp.current_value_usd DESC NULLS LAST, yp.entry_time DESC;