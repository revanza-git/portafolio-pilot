package services

import (
	"context"
	"fmt"
	"time"

	"github.com/defi-dashboard/backend/internal/models"
	"github.com/defi-dashboard/backend/internal/repos"
	"github.com/google/uuid"
)

type AlertService interface {
	CreateAlert(ctx context.Context, userID uuid.UUID, req *models.CreateAlertRequest) (*models.Alert, error)
	GetAlert(ctx context.Context, alertID uuid.UUID, userID uuid.UUID) (*models.Alert, error)
	GetUserAlerts(ctx context.Context, userID uuid.UUID, status *string, limit, offset int) ([]models.Alert, error)
	UpdateAlert(ctx context.Context, alertID uuid.UUID, userID uuid.UUID, req *models.UpdateAlertRequest) (*models.Alert, error)
	DeleteAlert(ctx context.Context, alertID uuid.UUID, userID uuid.UUID) error
	GetAlertHistory(ctx context.Context, alertID *uuid.UUID, userID uuid.UUID, limit, offset int) ([]models.AlertHistory, error)
	TriggerAlert(ctx context.Context, alertID uuid.UUID, triggeredValue map[string]interface{}) error
}

type alertService struct {
	alertRepo repos.AlertRepository
	userRepo  repos.UserRepository
}

func NewAlertService(alertRepo repos.AlertRepository, userRepo repos.UserRepository) AlertService {
	return &alertService{
		alertRepo: alertRepo,
		userRepo:  userRepo,
	}
}

func (s *alertService) CreateAlert(ctx context.Context, userID uuid.UUID, req *models.CreateAlertRequest) (*models.Alert, error) {
	// Validate user exists
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Validate alert type and conditions
	if err := s.validateAlertConditions(req.Type, req.Conditions); err != nil {
		return nil, fmt.Errorf("invalid alert conditions: %w", err)
	}

	alert := &models.Alert{
		ID:           uuid.New(),
		UserID:       userID,
		Type:         req.Type,
		Status:       models.AlertStatusActive,
		Target:       req.Target,
		Conditions:   req.Conditions,
		Notification: req.Notification,
		TriggerCount: 0,
	}

	if err := s.alertRepo.Create(ctx, alert); err != nil {
		return nil, fmt.Errorf("failed to create alert: %w", err)
	}

	return alert, nil
}

func (s *alertService) GetAlert(ctx context.Context, alertID uuid.UUID, userID uuid.UUID) (*models.Alert, error) {
	alert, err := s.alertRepo.GetByID(ctx, alertID)
	if err != nil {
		return nil, fmt.Errorf("failed to get alert: %w", err)
	}

	// Verify ownership
	if alert.UserID != userID {
		return nil, fmt.Errorf("alert not found")
	}

	return alert, nil
}

func (s *alertService) GetUserAlerts(ctx context.Context, userID uuid.UUID, status *string, limit, offset int) ([]models.Alert, error) {
	// Validate limit and offset
	if limit <= 0 || limit > 100 {
		limit = 20 // default
	}
	if offset < 0 {
		offset = 0
	}

	// Validate status if provided
	if status != nil {
		validStatuses := map[string]bool{
			models.AlertStatusActive:    true,
			models.AlertStatusTriggered: true,
			models.AlertStatusExpired:   true,
			models.AlertStatusDisabled:  true,
		}
		if !validStatuses[*status] {
			return nil, fmt.Errorf("invalid status: %s", *status)
		}
	}

	alerts, err := s.alertRepo.GetByUserID(ctx, userID, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user alerts: %w", err)
	}

	return alerts, nil
}

func (s *alertService) UpdateAlert(ctx context.Context, alertID uuid.UUID, userID uuid.UUID, req *models.UpdateAlertRequest) (*models.Alert, error) {
	// Get existing alert and verify ownership
	alert, err := s.GetAlert(ctx, alertID, userID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Status != nil {
		alert.Status = *req.Status
	}
	if req.Conditions != nil {
		// Validate new conditions
		if err := s.validateAlertConditions(alert.Type, *req.Conditions); err != nil {
			return nil, fmt.Errorf("invalid alert conditions: %w", err)
		}
		alert.Conditions = *req.Conditions
	}
	if req.Notification != nil {
		alert.Notification = *req.Notification
	}

	if err := s.alertRepo.Update(ctx, alert); err != nil {
		return nil, fmt.Errorf("failed to update alert: %w", err)
	}

	return alert, nil
}

func (s *alertService) DeleteAlert(ctx context.Context, alertID uuid.UUID, userID uuid.UUID) error {
	// Verify ownership first
	_, err := s.GetAlert(ctx, alertID, userID)
	if err != nil {
		return err
	}

	if err := s.alertRepo.Delete(ctx, alertID); err != nil {
		return fmt.Errorf("failed to delete alert: %w", err)
	}

	return nil
}

func (s *alertService) GetAlertHistory(ctx context.Context, alertID *uuid.UUID, userID uuid.UUID, limit, offset int) ([]models.AlertHistory, error) {
	// Validate limit and offset
	if limit <= 0 || limit > 100 {
		limit = 20 // default
	}
	if offset < 0 {
		offset = 0
	}

	// If alertID is provided, verify ownership
	if alertID != nil {
		_, err := s.GetAlert(ctx, *alertID, userID)
		if err != nil {
			return nil, err
		}
	}

	history, err := s.alertRepo.GetHistory(ctx, alertID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get alert history: %w", err)
	}

	// Filter by user ownership if no specific alertID provided
	if alertID == nil {
		var filteredHistory []models.AlertHistory
		for _, h := range history {
			// Get alert to verify ownership
			alert, err := s.alertRepo.GetByID(ctx, h.AlertID)
			if err != nil {
				continue // Skip if can't access alert
			}
			if alert.UserID == userID {
				filteredHistory = append(filteredHistory, h)
			}
		}
		history = filteredHistory
	}

	return history, nil
}

func (s *alertService) TriggerAlert(ctx context.Context, alertID uuid.UUID, triggeredValue map[string]interface{}) error {
	// Get alert
	alert, err := s.alertRepo.GetByID(ctx, alertID)
	if err != nil {
		return fmt.Errorf("failed to get alert: %w", err)
	}

	// Update alert as triggered
	if err := s.alertRepo.UpdateTriggered(ctx, alertID); err != nil {
		return fmt.Errorf("failed to update alert: %w", err)
	}

	// Create history record
	history := &models.AlertHistory{
		ID:                 uuid.New(),
		AlertID:            alertID,
		TriggeredAt:        time.Now(),
		ConditionsSnapshot: alert.Conditions,
		TriggeredValue:     triggeredValue,
		NotificationSent:   false, // TODO: Implement notification sending
	}

	if err := s.alertRepo.CreateHistory(ctx, history); err != nil {
		return fmt.Errorf("failed to create alert history: %w", err)
	}

	// TODO: Send notifications (email, webhook)
	// This would be implemented based on alert.Notification preferences

	return nil
}

// validateAlertConditions validates that the conditions are appropriate for the alert type
func (s *alertService) validateAlertConditions(alertType string, conditions models.AlertConditions) error {
	switch alertType {
	case models.AlertTypePriceAbove, models.AlertTypePriceBelow:
		if conditions.Price == nil || *conditions.Price <= 0 {
			return fmt.Errorf("price must be specified and greater than 0 for price alerts")
		}
	case models.AlertTypeLargeTransfer:
		if conditions.Threshold == nil || *conditions.Threshold == "" {
			return fmt.Errorf("threshold must be specified for transfer alerts")
		}
	case models.AlertTypeLiquidityChange:
		if conditions.ChangePercent == nil || *conditions.ChangePercent <= 0 {
			return fmt.Errorf("changePercent must be specified and greater than 0 for liquidity alerts")
		}
	case models.AlertTypeAPRChange:
		if conditions.MinAPR == nil && conditions.MaxAPR == nil {
			return fmt.Errorf("either minAPR or maxAPR must be specified for APR alerts")
		}
		if conditions.MinAPR != nil && *conditions.MinAPR < 0 {
			return fmt.Errorf("minAPR must be non-negative")
		}
		if conditions.MaxAPR != nil && *conditions.MaxAPR < 0 {
			return fmt.Errorf("maxAPR must be non-negative")
		}
	case models.AlertTypeApproval:
		// No specific conditions required for approval alerts
	default:
		return fmt.Errorf("unknown alert type: %s", alertType)
	}

	return nil
}