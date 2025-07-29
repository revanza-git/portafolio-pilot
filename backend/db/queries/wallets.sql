-- name: GetWalletsByUserId :many
SELECT * FROM wallets
WHERE user_id = $1
ORDER BY is_primary DESC, created_at ASC;

-- name: GetWalletByAddress :one
SELECT * FROM wallets
WHERE address = $1 AND chain_id = $2
LIMIT 1;

-- name: CreateWallet :one
INSERT INTO wallets (user_id, address, chain_id, label, is_primary)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateWallet :one
UPDATE wallets
SET label = $3, updated_at = NOW()
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: SetPrimaryWallet :exec
UPDATE wallets
SET is_primary = CASE
    WHEN id = $2 THEN true
    ELSE false
END
WHERE user_id = $1;

-- name: DeleteWallet :exec
DELETE FROM wallets
WHERE id = $1 AND user_id = $2;

-- name: GetWalletById :one
SELECT * FROM wallets
WHERE id = $1
LIMIT 1;