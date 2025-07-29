-- name: GetWalletBalances :many
SELECT 
    b.*,
    t.address as token_address,
    t.symbol,
    t.name,
    t.decimals,
    t.logo_uri,
    t.price_usd as current_price,
    t.price_change_24h
FROM balances b
INNER JOIN tokens t ON t.id = b.token_id
WHERE b.wallet_id = $1
    AND b.balance > 0
ORDER BY b.balance_usd DESC NULLS LAST;

-- name: GetOrCreateBalance :one
INSERT INTO balances (wallet_id, token_id, balance, balance_usd)
VALUES ($1, $2, $3, $4)
ON CONFLICT (wallet_id, token_id) 
DO UPDATE SET
    balance = $3,
    balance_usd = $4,
    updated_at = NOW()
RETURNING *;

-- name: UpdateBalance :one
UPDATE balances
SET balance = $3,
    balance_usd = $4,
    block_number = $5,
    updated_at = NOW()
WHERE wallet_id = $1 AND token_id = $2
RETURNING *;

-- name: RecordBalanceHistory :exec
INSERT INTO balance_history (wallet_id, token_id, balance, balance_usd, block_number)
VALUES ($1, $2, $3, $4, $5);

-- name: GetPortfolioHistory :many
SELECT 
    recorded_at,
    SUM(balance_usd) as total_value
FROM balance_history
WHERE wallet_id = ANY($1::uuid[])
    AND recorded_at >= $2
    AND recorded_at <= $3
GROUP BY recorded_at
ORDER BY recorded_at ASC;

-- name: GetUserTotalBalance :one
SELECT 
    COALESCE(SUM(b.balance_usd), 0) as total_balance_usd
FROM balances b
INNER JOIN wallets w ON w.id = b.wallet_id
WHERE w.user_id = $1
    AND ($2::int IS NULL OR w.chain_id = $2);