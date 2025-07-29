-- name: GetActiveAlerts :many
SELECT id, user_id, type, status, target, conditions, 
       notification, last_triggered_at, created_at
FROM alerts
WHERE status = 'active'
  AND (last_triggered_at IS NULL 
       OR last_triggered_at < NOW() - INTERVAL '1 hour')
ORDER BY created_at;

-- name: UpdateAlertTriggered :exec
UPDATE alerts 
SET last_triggered_at = NOW(),
    trigger_count = trigger_count + 1
WHERE id = $1;

-- name: CreateAlertHistory :exec
INSERT INTO alert_history (
    alert_id, triggered_at, conditions_snapshot, 
    triggered_value, notification_sent, notification_error
) VALUES ($1, NOW(), $2, $3, $4, $5);

-- name: GetAlert :one
SELECT * FROM alerts
WHERE id = $1
LIMIT 1;

-- name: GetUserAlerts :many
SELECT * FROM alerts
WHERE user_id = $1
  AND ($2::alert_status IS NULL OR status = $2)
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: CreateAlert :one
INSERT INTO alerts (
    user_id, type, status, target, conditions, notification
) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateAlert :one
UPDATE alerts
SET status = $2,
    conditions = $3,
    notification = $4,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteAlert :exec
DELETE FROM alerts
WHERE id = $1;