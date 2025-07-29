-- name: UpsertYieldPool :exec
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
    stable_coin = $11,
    updated_at = NOW();

-- name: GetYieldPool :one
SELECT * FROM yield_pools
WHERE pool_id = $1
LIMIT 1;

-- name: GetYieldPools :many
SELECT * FROM yield_pools
WHERE ($1::varchar IS NULL OR chain = $1)
  AND ($2::decimal IS NULL OR tvl_usd >= $2)
  AND ($3::decimal IS NULL OR apy >= $3)
ORDER BY 
  CASE WHEN $4 = 'apy' THEN apy END DESC,
  CASE WHEN $4 = 'tvl' THEN tvl_usd END DESC,
  CASE WHEN $4 = 'name' THEN pool_name END ASC
LIMIT $5 OFFSET $6;

-- name: GetYieldPoolHistory :many
SELECT pool_id, tvl_usd, apy, updated_at
FROM yield_pools
WHERE pool_id = $1
  AND updated_at >= $2
ORDER BY updated_at DESC;

-- name: GetPoolAPR :one
SELECT apy FROM yield_pools 
WHERE pool_id = $1
LIMIT 1;

-- name: GetPoolTVLWithPrevious :one
SELECT tvl_usd,
       LAG(tvl_usd) OVER (ORDER BY updated_at DESC) as prev_tvl
FROM yield_pools
WHERE pool_id = $1
ORDER BY updated_at DESC
LIMIT 1;

-- name: CountYieldPools :one
SELECT COUNT(*) FROM yield_pools yp
LEFT JOIN protocols p ON yp.protocol_id = p.id
WHERE ($1::varchar IS NULL OR yp.chain = $1)
  AND ($2::integer IS NULL OR yp.chain_id = $2)
  AND ($3::decimal IS NULL OR yp.tvl_usd >= $3)
  AND ($4::decimal IS NULL OR yp.apy >= $4)
  AND ($5::varchar IS NULL OR p.slug = $5)
  AND ($6::varchar IS NULL OR yp.risk_level = $6)
  AND ($7::boolean IS NULL OR yp.is_active = $7);

-- name: GetYieldPoolByID :one
SELECT yp.*, p.name as protocol_name, p.logo_uri as protocol_logo_uri
FROM yield_pools yp
LEFT JOIN protocols p ON yp.protocol_id = p.id
WHERE yp.id = $1
LIMIT 1;

-- name: GetYieldPoolsByProtocol :many
SELECT * FROM yield_pools
WHERE protocol_id = $1
  AND ($2::boolean IS NULL OR is_active = $2)
ORDER BY tvl_usd DESC;

-- name: UpdateYieldPoolAPY :exec
UPDATE yield_pools 
SET apy = $2,
    apy_base = $3,
    apy_reward = $4,
    updated_at = NOW()
WHERE pool_id = $1;

-- name: UpdateYieldPoolTVL :exec
UPDATE yield_pools 
SET tvl_usd = $2,
    updated_at = NOW()
WHERE pool_id = $1;

-- name: DeactivateYieldPool :exec
UPDATE yield_pools 
SET is_active = false,
    updated_at = NOW()
WHERE pool_id = $1;

-- name: GetActiveYieldPoolsByChain :many
SELECT yp.*, p.name as protocol_name
FROM yield_pools yp
LEFT JOIN protocols p ON yp.protocol_id = p.id
WHERE yp.chain_id = $1
  AND yp.is_active = true
ORDER BY yp.apy DESC;

-- name: GetTopYieldPoolsByTVL :many
SELECT yp.*, p.name as protocol_name
FROM yield_pools yp
LEFT JOIN protocols p ON yp.protocol_id = p.id
WHERE yp.is_active = true
  AND yp.tvl_usd IS NOT NULL
ORDER BY yp.tvl_usd DESC
LIMIT $1;

-- name: GetEnhancedYieldPools :many
SELECT yp.*, p.name as protocol_name, p.logo_uri as protocol_logo_uri, p.category as protocol_category
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
LIMIT $9 OFFSET $10;