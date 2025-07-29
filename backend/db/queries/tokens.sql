-- name: GetTokenByAddress :one
SELECT * FROM tokens
WHERE address = $1 AND chain_id = $2
LIMIT 1;

-- name: GetTokenById :one
SELECT * FROM tokens
WHERE id = $1
LIMIT 1;

-- name: CreateToken :one
INSERT INTO tokens (
    address, chain_id, symbol, name, decimals,
    logo_uri, price_usd, price_change_24h, market_cap,
    total_supply, last_updated
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
)
RETURNING *;

-- name: UpdateTokenPrice :one
UPDATE tokens
SET price_usd = $3,
    price_change_24h = $4,
    market_cap = $5,
    last_updated = NOW(),
    updated_at = NOW()
WHERE address = $1 AND chain_id = $2
RETURNING *;

-- name: GetTokensByChainId :many
SELECT * FROM tokens
WHERE chain_id = $1
ORDER BY market_cap DESC NULLS LAST
LIMIT $2 OFFSET $3;

-- name: SearchTokens :many
SELECT * FROM tokens
WHERE (symbol ILIKE $1 || '%' OR name ILIKE $1 || '%')
    AND ($2::int IS NULL OR chain_id = $2)
ORDER BY market_cap DESC NULLS LAST
LIMIT 20;