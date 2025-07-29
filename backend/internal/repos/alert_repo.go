package repos

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/defi-dashboard/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AlertRepository interface {
	Create(ctx context.Context, alert *models.Alert) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Alert, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, status *string, limit, offset int) ([]models.Alert, error)
	Update(ctx context.Context, alert *models.Alert) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetActiveAlerts(ctx context.Context) ([]models.Alert, error)
	UpdateTriggered(ctx context.Context, alertID uuid.UUID) error
	CreateHistory(ctx context.Context, history *models.AlertHistory) error
	GetHistory(ctx context.Context, alertID *uuid.UUID, limit, offset int) ([]models.AlertHistory, error)
}

type alertRepository struct {
	db *pgxpool.Pool
}

func NewAlertRepository(db *pgxpool.Pool) AlertRepository {
	return &alertRepository{db: db}
}

func (r *alertRepository) Create(ctx context.Context, alert *models.Alert) error {
	targetJSON, err := json.Marshal(alert.Target)
	if err != nil {
		return fmt.Errorf("failed to marshal target: %w", err)
	}

	conditionsJSON, err := json.Marshal(alert.Conditions)
	if err != nil {
		return fmt.Errorf("failed to marshal conditions: %w", err)
	}

	notificationJSON, err := json.Marshal(alert.Notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	query := `
		INSERT INTO alerts (
			id, user_id, type, status, target, conditions, notification
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = r.db.Exec(ctx, query,
		alert.ID,
		alert.UserID,
		alert.Type,
		alert.Status,
		targetJSON,
		conditionsJSON,
		notificationJSON,
	)

	if err != nil {
		return fmt.Errorf("failed to create alert: %w", err)
	}

	// Fetch the created alert to get timestamps
	return r.populateAlertFromDB(ctx, alert.ID, alert)
}

func (r *alertRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Alert, error) {
	alert := &models.Alert{}
	err := r.populateAlertFromDB(ctx, id, alert)
	if err != nil {
		return nil, err
	}
	return alert, nil
}

func (r *alertRepository) GetByUserID(ctx context.Context, userID uuid.UUID, status *string, limit, offset int) ([]models.Alert, error) {
	query := `
		SELECT id, user_id, type, status, target, conditions, 
			   notification, last_triggered_at, trigger_count, created_at, updated_at
		FROM alerts
		WHERE user_id = $1
		  AND ($2::alert_status IS NULL OR status = $2)
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.Query(ctx, query, userID, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get alerts: %w", err)
	}
	defer rows.Close()

	return r.scanAlerts(rows)
}

func (r *alertRepository) Update(ctx context.Context, alert *models.Alert) error {
	conditionsJSON, err := json.Marshal(alert.Conditions)
	if err != nil {
		return fmt.Errorf("failed to marshal conditions: %w", err)
	}

	notificationJSON, err := json.Marshal(alert.Notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	query := `
		UPDATE alerts
		SET status = $2,
		    conditions = $3,
		    notification = $4,
		    updated_at = NOW()
		WHERE id = $1
	`

	_, err = r.db.Exec(ctx, query,
		alert.ID,
		alert.Status,
		conditionsJSON,
		notificationJSON,
	)

	if err != nil {
		return fmt.Errorf("failed to update alert: %w", err)
	}

	// Fetch updated alert to get new timestamps
	return r.populateAlertFromDB(ctx, alert.ID, alert)
}

func (r *alertRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM alerts WHERE id = $1`
	
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete alert: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("alert not found")
	}

	return nil
}

func (r *alertRepository) GetActiveAlerts(ctx context.Context) ([]models.Alert, error) {
	query := `
		SELECT id, user_id, type, status, target, conditions, 
			   notification, last_triggered_at, trigger_count, created_at, updated_at
		FROM alerts
		WHERE status = 'active'
		  AND (last_triggered_at IS NULL 
		       OR last_triggered_at < NOW() - INTERVAL '1 hour')
		ORDER BY created_at
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get active alerts: %w", err)
	}
	defer rows.Close()

	return r.scanAlerts(rows)
}

func (r *alertRepository) UpdateTriggered(ctx context.Context, alertID uuid.UUID) error {
	query := `
		UPDATE alerts 
		SET last_triggered_at = NOW(),
		    trigger_count = trigger_count + 1
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, alertID)
	if err != nil {
		return fmt.Errorf("failed to update alert trigger: %w", err)
	}

	return nil
}

func (r *alertRepository) CreateHistory(ctx context.Context, history *models.AlertHistory) error {
	conditionsJSON, err := json.Marshal(history.ConditionsSnapshot)
	if err != nil {
		return fmt.Errorf("failed to marshal conditions snapshot: %w", err)
	}

	triggeredValueJSON, err := json.Marshal(history.TriggeredValue)
	if err != nil {
		return fmt.Errorf("failed to marshal triggered value: %w", err)
	}

	query := `
		INSERT INTO alert_history (
		    id, alert_id, triggered_at, conditions_snapshot, 
		    triggered_value, notification_sent, notification_error
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = r.db.Exec(ctx, query,
		history.ID,
		history.AlertID,
		history.TriggeredAt,
		conditionsJSON,
		triggeredValueJSON,
		history.NotificationSent,
		history.NotificationError,
	)

	if err != nil {
		return fmt.Errorf("failed to create alert history: %w", err)
	}

	return nil
}

func (r *alertRepository) GetHistory(ctx context.Context, alertID *uuid.UUID, limit, offset int) ([]models.AlertHistory, error) {
	var query string
	var args []interface{}

	if alertID != nil {
		query = `
			SELECT id, alert_id, triggered_at, conditions_snapshot,
				   triggered_value, notification_sent, notification_error
			FROM alert_history
			WHERE alert_id = $1
			ORDER BY triggered_at DESC
			LIMIT $2 OFFSET $3
		`
		args = []interface{}{*alertID, limit, offset}
	} else {
		query = `
			SELECT id, alert_id, triggered_at, conditions_snapshot,
				   triggered_value, notification_sent, notification_error
			FROM alert_history
			ORDER BY triggered_at DESC
			LIMIT $1 OFFSET $2
		`
		args = []interface{}{limit, offset}
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get alert history: %w", err)
	}
	defer rows.Close()

	var history []models.AlertHistory
	for rows.Next() {
		var h models.AlertHistory
		var conditionsJSON, triggeredValueJSON []byte

		err := rows.Scan(
			&h.ID,
			&h.AlertID,
			&h.TriggeredAt,
			&conditionsJSON,
			&triggeredValueJSON,
			&h.NotificationSent,
			&h.NotificationError,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan alert history: %w", err)
		}

		// Unmarshal JSON fields
		if err := json.Unmarshal(conditionsJSON, &h.ConditionsSnapshot); err != nil {
			return nil, fmt.Errorf("failed to unmarshal conditions snapshot: %w", err)
		}
		if err := json.Unmarshal(triggeredValueJSON, &h.TriggeredValue); err != nil {
			return nil, fmt.Errorf("failed to unmarshal triggered value: %w", err)
		}

		history = append(history, h)
	}

	return history, rows.Err()
}

// Helper methods

func (r *alertRepository) populateAlertFromDB(ctx context.Context, id uuid.UUID, alert *models.Alert) error {
	query := `
		SELECT id, user_id, type, status, target, conditions, 
			   notification, last_triggered_at, trigger_count, created_at, updated_at
		FROM alerts
		WHERE id = $1
	`

	var targetJSON, conditionsJSON, notificationJSON []byte

	row := r.db.QueryRow(ctx, query, id)
	err := row.Scan(
		&alert.ID,
		&alert.UserID,
		&alert.Type,
		&alert.Status,
		&targetJSON,
		&conditionsJSON,
		&notificationJSON,
		&alert.LastTriggeredAt,
		&alert.TriggerCount,
		&alert.CreatedAt,
		&alert.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("alert not found")
		}
		return fmt.Errorf("failed to get alert: %w", err)
	}

	// Unmarshal JSON fields
	if err := json.Unmarshal(targetJSON, &alert.Target); err != nil {
		return fmt.Errorf("failed to unmarshal target: %w", err)
	}
	if err := json.Unmarshal(conditionsJSON, &alert.Conditions); err != nil {
		return fmt.Errorf("failed to unmarshal conditions: %w", err)
	}
	if err := json.Unmarshal(notificationJSON, &alert.Notification); err != nil {
		return fmt.Errorf("failed to unmarshal notification: %w", err)
	}

	return nil
}

func (r *alertRepository) scanAlerts(rows pgx.Rows) ([]models.Alert, error) {
	var alerts []models.Alert

	for rows.Next() {
		var alert models.Alert
		var targetJSON, conditionsJSON, notificationJSON []byte

		err := rows.Scan(
			&alert.ID,
			&alert.UserID,
			&alert.Type,
			&alert.Status,
			&targetJSON,
			&conditionsJSON,
			&notificationJSON,
			&alert.LastTriggeredAt,
			&alert.TriggerCount,
			&alert.CreatedAt,
			&alert.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan alert: %w", err)
		}

		// Unmarshal JSON fields
		if err := json.Unmarshal(targetJSON, &alert.Target); err != nil {
			return nil, fmt.Errorf("failed to unmarshal target: %w", err)
		}
		if err := json.Unmarshal(conditionsJSON, &alert.Conditions); err != nil {
			return nil, fmt.Errorf("failed to unmarshal conditions: %w", err)
		}
		if err := json.Unmarshal(notificationJSON, &alert.Notification); err != nil {
			return nil, fmt.Errorf("failed to unmarshal notification: %w", err)
		}

		alerts = append(alerts, alert)
	}

	return alerts, rows.Err()
}