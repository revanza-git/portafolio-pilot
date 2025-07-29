-- name: GetProtocol :one
SELECT * FROM protocols 
WHERE id = $1 AND is_active = true
LIMIT 1;

-- name: GetProtocolBySlug :one
SELECT * FROM protocols 
WHERE slug = $1 AND is_active = true
LIMIT 1;

-- name: GetProtocols :many
SELECT * FROM protocols
WHERE ($1::varchar IS NULL OR category = $1)
  AND ($2::boolean IS NULL OR is_active = $2)
  AND ($3::varchar IS NULL OR risk_level = $3)
ORDER BY 
  CASE WHEN $4 = 'name' THEN name END ASC,
  CASE WHEN $4 = 'tvl' THEN total_tvl_usd END DESC,
  CASE WHEN $4 = 'category' THEN category END ASC,
  created_at DESC
LIMIT $5 OFFSET $6;

-- name: CountProtocols :one
SELECT COUNT(*) FROM protocols
WHERE ($1::varchar IS NULL OR category = $1)
  AND ($2::boolean IS NULL OR is_active = $2)
  AND ($3::varchar IS NULL OR risk_level = $3);

-- name: CreateProtocol :one
INSERT INTO protocols (
    name, slug, description, website_url, logo_uri, 
    category, total_tvl_usd, chains, is_active, risk_level
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: UpdateProtocol :one
UPDATE protocols 
SET name = $2,
    description = $3,
    website_url = $4,
    logo_uri = $5,
    category = $6,
    total_tvl_usd = $7,
    chains = $8,
    is_active = $9,
    risk_level = $10,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateProtocolTVL :exec
UPDATE protocols 
SET total_tvl_usd = $2,
    updated_at = NOW()
WHERE id = $1;

-- name: DeleteProtocol :exec
UPDATE protocols 
SET is_active = false,
    updated_at = NOW()
WHERE id = $1;

-- name: GetProtocolsByChain :many
SELECT * FROM protocols
WHERE chains ? $1::text
  AND is_active = true
ORDER BY total_tvl_usd DESC;

-- name: GetProtocolsWithPoolCount :many
SELECT p.*, 
       COUNT(yp.id) as pool_count,
       COALESCE(SUM(yp.tvl_usd), 0) as total_pools_tvl
FROM protocols p
LEFT JOIN yield_pools yp ON p.id = yp.protocol_id AND yp.is_active = true
WHERE p.is_active = true
GROUP BY p.id
ORDER BY total_pools_tvl DESC
LIMIT $1 OFFSET $2;