-- name: GetTransactionByHash :one
SELECT * FROM transactions
WHERE hash = $1
LIMIT 1;

-- name: CreateTransaction :one
INSERT INTO transactions (
    hash, chain_id, from_address, to_address, value,
    gas_used, gas_price, gas_fee_usd, block_number,
    timestamp, status, type, metadata
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
)
RETURNING *;

-- name: UpdateTransactionStatus :one
UPDATE transactions
SET status = $2,
    block_number = $3,
    gas_used = $4,
    gas_fee_usd = $5,
    updated_at = NOW()
WHERE hash = $1
RETURNING *;

-- name: GetUserTransactions :many
SELECT t.* FROM transactions t
INNER JOIN user_transactions ut ON ut.transaction_id = t.id
WHERE ut.user_id = $1
    AND ($2::int IS NULL OR t.chain_id = $2)
    AND ($3::transaction_type IS NULL OR t.type = $3)
    AND ($4::timestamptz IS NULL OR t.timestamp >= $4)
    AND ($5::timestamptz IS NULL OR t.timestamp <= $5)
ORDER BY t.timestamp DESC
LIMIT $6 OFFSET $7;

-- name: GetWalletTransactions :many
SELECT * FROM transactions
WHERE (from_address = $1 OR to_address = $1)
    AND chain_id = $2
    AND ($3::transaction_type IS NULL OR type = $3)
ORDER BY timestamp DESC
LIMIT $4 OFFSET $5;

-- name: LinkTransactionToUser :exec
INSERT INTO user_transactions (user_id, transaction_id, wallet_id)
VALUES ($1, $2, $3)
ON CONFLICT DO NOTHING;