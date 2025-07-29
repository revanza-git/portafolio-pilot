package repos

import (
	"context"
	"testing"
	"time"

	"github.com/defi-dashboard/backend/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockDB for testing - would use actual test database in real implementation
type MockAlertRepo struct {
	alerts  map[uuid.UUID]*models.Alert
	history map[uuid.UUID]*models.AlertHistory
}

func NewMockAlertRepo() *MockAlertRepo {
	return &MockAlertRepo{
		alerts:  make(map[uuid.UUID]*models.Alert),
		history: make(map[uuid.UUID]*models.AlertHistory),
	}
}

func (m *MockAlertRepo) Create(ctx context.Context, alert *models.Alert) error {
	alert.CreatedAt = time.Now()
	alert.UpdatedAt = time.Now()
	m.alerts[alert.ID] = alert
	return nil
}

func (m *MockAlertRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Alert, error) {
	alert, exists := m.alerts[id]
	if !exists {
		return nil, assert.AnError
	}
	return alert, nil
}

func (m *MockAlertRepo) GetByUserID(ctx context.Context, userID uuid.UUID, status *string, limit, offset int) ([]models.Alert, error) {
	var alerts []models.Alert
	for _, alert := range m.alerts {
		if alert.UserID == userID {
			if status == nil || alert.Status == *status {
				alerts = append(alerts, *alert)
			}
		}
	}
	return alerts, nil
}

func (m *MockAlertRepo) Update(ctx context.Context, alert *models.Alert) error {
	alert.UpdatedAt = time.Now()
	m.alerts[alert.ID] = alert
	return nil
}

func (m *MockAlertRepo) Delete(ctx context.Context, id uuid.UUID) error {
	delete(m.alerts, id)
	return nil
}

func (m *MockAlertRepo) GetActiveAlerts(ctx context.Context) ([]models.Alert, error) {
	var alerts []models.Alert
	for _, alert := range m.alerts {
		if alert.Status == models.AlertStatusActive {
			alerts = append(alerts, *alert)
		}
	}
	return alerts, nil
}

func (m *MockAlertRepo) UpdateTriggered(ctx context.Context, alertID uuid.UUID) error {
	if alert, exists := m.alerts[alertID]; exists {
		now := time.Now()
		alert.LastTriggeredAt = &now
		alert.TriggerCount++
	}
	return nil
}

func (m *MockAlertRepo) CreateHistory(ctx context.Context, history *models.AlertHistory) error {
	m.history[history.ID] = history
	return nil
}

func (m *MockAlertRepo) GetHistory(ctx context.Context, alertID *uuid.UUID, limit, offset int) ([]models.AlertHistory, error) {
	var history []models.AlertHistory
	for _, h := range m.history {
		if alertID == nil || h.AlertID == *alertID {
			history = append(history, *h)
		}
	}
	return history, nil
}

func TestAlertRepository(t *testing.T) {
	ctx := context.Background()
	repo := NewMockAlertRepo()

	t.Run("Create and Get Alert", func(t *testing.T) {
		userID := uuid.New()
		price := 3000.0
		
		alert := &models.Alert{
			ID:     uuid.New(),
			UserID: userID,
			Type:   models.AlertTypePriceAbove,
			Status: models.AlertStatusActive,
			Target: models.AlertTarget{
				Type:       "token",
				Identifier: "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
				ChainID:    1,
			},
			Conditions: models.AlertConditions{
				Price: &price,
			},
			Notification: models.AlertNotification{
				Email: true,
			},
		}

		// Create alert
		err := repo.Create(ctx, alert)
		require.NoError(t, err)
		assert.NotZero(t, alert.CreatedAt)

		// Get alert by ID
		retrieved, err := repo.GetByID(ctx, alert.ID)
		require.NoError(t, err)
		assert.Equal(t, alert.ID, retrieved.ID)
		assert.Equal(t, alert.Type, retrieved.Type)
		assert.Equal(t, alert.Status, retrieved.Status)
		assert.Equal(t, *alert.Conditions.Price, *retrieved.Conditions.Price)
	})

	t.Run("Update Alert", func(t *testing.T) {
		userID := uuid.New()
		price := 2500.0
		
		alert := &models.Alert{
			ID:     uuid.New(),
			UserID: userID,
			Type:   models.AlertTypePriceBelow,
			Status: models.AlertStatusActive,
			Target: models.AlertTarget{
				Type:       "token",
				Identifier: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
				ChainID:    1,
			},
			Conditions: models.AlertConditions{
				Price: &price,
			},
			Notification: models.AlertNotification{
				Email: false,
			},
		}

		// Create alert
		err := repo.Create(ctx, alert)
		require.NoError(t, err)

		// Update alert
		alert.Status = models.AlertStatusDisabled
		newPrice := 2000.0
		alert.Conditions.Price = &newPrice
		
		err = repo.Update(ctx, alert)
		require.NoError(t, err)

		// Verify update
		retrieved, err := repo.GetByID(ctx, alert.ID)
		require.NoError(t, err)
		assert.Equal(t, models.AlertStatusDisabled, retrieved.Status)
		assert.Equal(t, newPrice, *retrieved.Conditions.Price)
	})

	t.Run("Delete Alert", func(t *testing.T) {
		userID := uuid.New()
		price := 1500.0
		
		alert := &models.Alert{
			ID:     uuid.New(),
			UserID: userID,
			Type:   models.AlertTypePriceAbove,
			Status: models.AlertStatusActive,
			Conditions: models.AlertConditions{
				Price: &price,
			},
		}

		// Create alert
		err := repo.Create(ctx, alert)
		require.NoError(t, err)

		// Delete alert
		err = repo.Delete(ctx, alert.ID)
		require.NoError(t, err)

		// Verify deletion
		_, err = repo.GetByID(ctx, alert.ID)
		assert.Error(t, err)
	})

	t.Run("Get Active Alerts", func(t *testing.T) {
		userID := uuid.New()
		price1 := 3000.0
		price2 := 2000.0

		// Create active alert
		activeAlert := &models.Alert{
			ID:     uuid.New(),
			UserID: userID,
			Type:   models.AlertTypePriceAbove,
			Status: models.AlertStatusActive,
			Conditions: models.AlertConditions{
				Price: &price1,
			},
		}

		// Create disabled alert
		disabledAlert := &models.Alert{
			ID:     uuid.New(),
			UserID: userID,
			Type:   models.AlertTypePriceBelow,
			Status: models.AlertStatusDisabled,
			Conditions: models.AlertConditions{
				Price: &price2,
			},
		}

		err := repo.Create(ctx, activeAlert)
		require.NoError(t, err)
		err = repo.Create(ctx, disabledAlert)
		require.NoError(t, err)

		// Get active alerts
		activeAlerts, err := repo.GetActiveAlerts(ctx)
		require.NoError(t, err)
		
		// Should only return the active alert
		found := false
		for _, alert := range activeAlerts {
			if alert.ID == activeAlert.ID {
				found = true
				assert.Equal(t, models.AlertStatusActive, alert.Status)
			}
			assert.NotEqual(t, disabledAlert.ID, alert.ID) // Should not include disabled alert
		}
		assert.True(t, found, "Active alert should be found")
	})

	t.Run("Alert History", func(t *testing.T) {
		alertID := uuid.New()
		price := 3000.0

		history := &models.AlertHistory{
			ID:      uuid.New(),
			AlertID: alertID,
			ConditionsSnapshot: models.AlertConditions{
				Price: &price,
			},
			TriggeredValue: map[string]interface{}{
				"currentPrice": 3100.0,
			},
			TriggeredAt:      time.Now(),
			NotificationSent: true,
		}

		// Create history
		err := repo.CreateHistory(ctx, history)
		require.NoError(t, err)

		// Get history
		historyList, err := repo.GetHistory(ctx, &alertID, 10, 0)
		require.NoError(t, err)
		assert.Len(t, historyList, 1)
		assert.Equal(t, history.AlertID, historyList[0].AlertID)
		assert.Equal(t, true, historyList[0].NotificationSent)
		assert.Equal(t, 3100.0, historyList[0].TriggeredValue["currentPrice"])
	})
}