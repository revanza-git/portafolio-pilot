-- name: GetUserByAddress :one
SELECT * FROM users
WHERE address = $1
LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (address, nonce)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateUserNonce :one
UPDATE users
SET nonce = $2, updated_at = NOW()
WHERE address = $1
RETURNING *;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1
LIMIT 1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;