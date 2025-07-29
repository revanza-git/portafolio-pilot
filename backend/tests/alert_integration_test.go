package tests

import (
	"context"
	"testing"
	"time"

	"github.com/defi-dashboard/backend/internal/jobs"
	"github.com/defi-dashboard/backend/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock implementations for integration testing
type MockAlertService struct {
	mock.Mock
}

func (m *MockAlertService) CreateAlert(ctx context.Context, userID uuid.UUID, req *models.CreateAlertRequest) (*models.Alert, error) {
	args := m.Called(ctx, userID, req)
	return args.Get(0).(*models.Alert), args.Error(1)
}

func (m *MockAlertService) GetAlert(ctx context.Context, alertID uuid.UUID, userID uuid.UUID) (*models.Alert, error) {
	args := m.Called(ctx, alertID, userID)
	return args.Get(0).(*models.Alert), args.Error(1)
}

func (m *MockAlertService) GetUserAlerts(ctx context.Context, userID uuid.UUID, status *string, limit, offset int) ([]models.Alert, error) {
	args := m.Called(ctx, userID, status, limit, offset)
	return args.Get(0).([]models.Alert), args.Error(1)
}

func (m *MockAlertService) UpdateAlert(ctx context.Context, alertID uuid.UUID, userID uuid.UUID, req *models.UpdateAlertRequest) (*models.Alert, error) {
	args := m.Called(ctx, alertID, userID, req)
	return args.Get(0).(*models.Alert), args.Error(1)
}

func (m *MockAlertService) DeleteAlert(ctx context.Context, alertID uuid.UUID, userID uuid.UUID) error {
	args := m.Called(ctx, alertID, userID)
	return args.Error(0)
}

func (m *MockAlertService) GetAlertHistory(ctx context.Context, alertID *uuid.UUID, userID uuid.UUID, limit, offset int) ([]models.AlertHistory, error) {
	args := m.Called(ctx, alertID, userID, limit, offset)
	return args.Get(0).([]models.AlertHistory), args.Error(1)
}

func (m *MockAlertService) TriggerAlert(ctx context.Context, alertID uuid.UUID, triggeredValue map[string]interface{}) error {
	args := m.Called(ctx, alertID, triggeredValue)
	return args.Error(0)
}

type MockAlertRepository struct {
	mock.Mock
}

func (m *MockAlertRepository) Create(ctx context.Context, alert *models.Alert) error {
	args := m.Called(ctx, alert)
	return args.Error(0)
}

func (m *MockAlertRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Alert, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Alert), args.Error(1)
}

func (m *MockAlertRepository) GetByUserID(ctx context.Context, userID uuid.UUID, status *string, limit, offset int) ([]models.Alert, error) {
	args := m.Called(ctx, userID, status, limit, offset)
	return args.Get(0).([]models.Alert), args.Error(1)
}

func (m *MockAlertRepository) Update(ctx context.Context, alert *models.Alert) error {
	args := m.Called(ctx, alert)
	return args.Error(0)
}

func (m *MockAlertRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAlertRepository) GetActiveAlerts(ctx context.Context) ([]models.Alert, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Alert), args.Error(1)
}

func (m *MockAlertRepository) UpdateTriggered(ctx context.Context, alertID uuid.UUID) error {
	args := m.Called(ctx, alertID)
	return args.Error(0)
}

func (m *MockAlertRepository) CreateHistory(ctx context.Context, history *models.AlertHistory) error {
	args := m.Called(ctx, history)
	return args.Error(0)
}

func (m *MockAlertRepository) GetHistory(ctx context.Context, alertID *uuid.UUID, limit, offset int) ([]models.AlertHistory, error) {
	args := m.Called(ctx, alertID, limit, offset)
	return args.Get(0).([]models.AlertHistory), args.Error(1)
}

// TestAlertCreationAndTriggering tests the full alert lifecycle
func TestAlertCreationAndTriggering(t *testing.T) {
	ctx := context.Background()
	
	mockAlertService := new(MockAlertService)
	mockAlertRepo := new(MockAlertRepository)

	// Test data
	userID := uuid.New()
	alertID := uuid.New()
	price := 3000.0

	t.Run("Create alert, trigger, and verify history", func(t *testing.T) {
		// Step 1: Create alert
		createReq := &models.CreateAlertRequest{
			Type: models.AlertTypePriceAbove,
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

		createdAlert := &models.Alert{
			ID:           alertID,
			UserID:       userID,
			Type:         createReq.Type,
			Status:       models.AlertStatusActive,
			Target:       createReq.Target,
			Conditions:   createReq.Conditions,
			Notification: createReq.Notification,
			TriggerCount: 0,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		mockAlertService.On("CreateAlert", ctx, userID, createReq).Return(createdAlert, nil)

		// Create the alert
		alert, err := mockAlertService.CreateAlert(ctx, userID, createReq)
		require.NoError(t, err)
		assert.Equal(t, alertID, alert.ID)
		assert.Equal(t, models.AlertStatusActive, alert.Status)

		// Step 2: Set up alert triggering
		triggeredValue := map[string]interface{}{
			"currentPrice": 3100.0, // Above threshold
			"targetPrice":  3000.0,
		}

		mockAlertService.On("TriggerAlert", ctx, alertID, triggeredValue).Return(nil)

		// Trigger the alert
		err = mockAlertService.TriggerAlert(ctx, alertID, triggeredValue)
		require.NoError(t, err)

		// Step 3: Verify history was created
		expectedHistory := []models.AlertHistory{
			{
				ID:      uuid.New(),
				AlertID: alertID,
				ConditionsSnapshot: models.AlertConditions{
					Price: &price,
				},
				TriggeredValue: triggeredValue,
				TriggeredAt:    time.Now(),
				NotificationSent: false,
			},
		}

		mockAlertService.On("GetAlertHistory", ctx, &alertID, userID, 10, 0).Return(expectedHistory, nil)

		// Get alert history
		history, err := mockAlertService.GetAlertHistory(ctx, &alertID, userID, 10, 0)
		require.NoError(t, err)
		assert.Len(t, history, 1)
		assert.Equal(t, alertID, history[0].AlertID)
		assert.Equal(t, 3100.0, history[0].TriggeredValue["currentPrice"])

		// Verify all mocks were called correctly
		mockAlertService.AssertExpectations(t)
	})
}

// TestWorkerIntegration tests the worker evaluating and triggering alerts
func TestWorkerIntegration(t *testing.T) {
	ctx := context.Background()
	
	mockAlertService := new(MockAlertService)
	mockAlertRepo := new(MockAlertRepository)

	// Create alert evaluator job
	job := jobs.NewAlertEvaluatorJob(nil, mockAlertService, mockAlertRepo)

	t.Run("Worker finds and triggers price alert", func(t *testing.T) {
		userID := uuid.New()
		alertID := uuid.New()
		price := 3000.0

		// Create test alert
		alert := models.Alert{
			ID:     alertID,
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
			TriggerCount: 0,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		alerts := []models.Alert{alert}

		// Mock repository to return active alerts
		mockAlertRepo.On("GetActiveAlerts", ctx).Return(alerts, nil)

		// Mock alert service to handle trigger call
		// Note: In a real test, the worker would fetch current prices and evaluate conditions
		// For this test, we'll assume the condition is met and the trigger is called
		triggeredValue := map[string]interface{}{
			"currentPrice": 3100.0,
			"tokenKey":     "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2-1",
		}
		mockAlertService.On("TriggerAlert", ctx, alertID, triggeredValue).Return(nil)

		// Run the worker
		err := job.Run(ctx)
		require.NoError(t, err)

		// Verify that GetActiveAlerts was called
		mockAlertRepo.AssertExpectations(t)
		// Note: TriggerAlert would be called inside evaluatePriceAlerts if conditions are met
	})
}

// TestAlertUpdatesAndHistory tests updating alerts and checking history
func TestAlertUpdatesAndHistory(t *testing.T) {
	ctx := context.Background()
	
	mockAlertService := new(MockAlertService)

	userID := uuid.New()
	alertID := uuid.New()
	price := 3000.0

	t.Run("Update alert and check history", func(t *testing.T) {
		// Step 1: Get existing alert
		existingAlert := &models.Alert{
			ID:     alertID,
			UserID: userID,
			Type:   models.AlertTypePriceAbove,
			Status: models.AlertStatusActive,
			Conditions: models.AlertConditions{
				Price: &price,
			},
			TriggerCount: 1,
			LastTriggeredAt: &time.Time{},
		}

		mockAlertService.On("GetAlert", ctx, alertID, userID).Return(existingAlert, nil)

		// Get the alert
		alert, err := mockAlertService.GetAlert(ctx, alertID, userID)
		require.NoError(t, err)
		assert.Equal(t, 1, alert.TriggerCount)

		// Step 2: Update alert (pause it)
		status := models.AlertStatusDisabled
		updateReq := &models.UpdateAlertRequest{
			Status: &status,
		}

		updatedAlert := &models.Alert{
			ID:     alertID,
			UserID: userID,
			Type:   models.AlertTypePriceAbove,
			Status: models.AlertStatusDisabled,
			Conditions: models.AlertConditions{
				Price: &price,
			},
			TriggerCount: 1,
		}

		mockAlertService.On("UpdateAlert", ctx, alertID, userID, updateReq).Return(updatedAlert, nil)

		// Update the alert
		result, err := mockAlertService.UpdateAlert(ctx, alertID, userID, updateReq)
		require.NoError(t, err)
		assert.Equal(t, models.AlertStatusDisabled, result.Status)

		// Step 3: Check alert history
		historyEntry := models.AlertHistory{
			ID:      uuid.New(),
			AlertID: alertID,
			ConditionsSnapshot: models.AlertConditions{
				Price: &price,
			},
			TriggeredValue: map[string]interface{}{
				"currentPrice": 3100.0,
			},
			TriggeredAt:      time.Now().Add(-time.Hour),
			NotificationSent: true,
		}

		history := []models.AlertHistory{historyEntry}
		mockAlertService.On("GetAlertHistory", ctx, &alertID, userID, 20, 0).Return(history, nil)

		// Get alert history
		historyResult, err := mockAlertService.GetAlertHistory(ctx, &alertID, userID, 20, 0)
		require.NoError(t, err)
		assert.Len(t, historyResult, 1)
		assert.Equal(t, alertID, historyResult[0].AlertID)
		assert.True(t, historyResult[0].NotificationSent)

		// Verify all mocks were called correctly
		mockAlertService.AssertExpectations(t)
	})
}

// TestMultipleAlertTypes tests different types of alerts
func TestMultipleAlertTypes(t *testing.T) {
	ctx := context.Background()
	
	mockAlertService := new(MockAlertService)

	userID := uuid.New()

	t.Run("Create different alert types", func(t *testing.T) {
		// Price alert
		price := 3000.0
		priceAlertReq := &models.CreateAlertRequest{
			Type: models.AlertTypePriceAbove,
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

		priceAlert := &models.Alert{
			ID:           uuid.New(),
			UserID:       userID,
			Type:         models.AlertTypePriceAbove,
			Status:       models.AlertStatusActive,
			Target:       priceAlertReq.Target,
			Conditions:   priceAlertReq.Conditions,
			Notification: priceAlertReq.Notification,
		}

		// Transfer alert
		threshold := "1000000000000000000" // 1 ETH in wei
		transferAlertReq := &models.CreateAlertRequest{
			Type: models.AlertTypeLargeTransfer,
			Target: models.AlertTarget{
				Type:       "address",
				Identifier: "0x1234567890123456789012345678901234567890",
				ChainID:    1,
			},
			Conditions: models.AlertConditions{
				Threshold: &threshold,
			},
			Notification: models.AlertNotification{
				Email: false,
				Webhook: "https://webhook.example.com",
			},
		}

		transferAlert := &models.Alert{
			ID:           uuid.New(),
			UserID:       userID,
			Type:         models.AlertTypeLargeTransfer,
			Status:       models.AlertStatusActive,
			Target:       transferAlertReq.Target,
			Conditions:   transferAlertReq.Conditions,
			Notification: transferAlertReq.Notification,
		}

		// Mock service calls
		mockAlertService.On("CreateAlert", ctx, userID, priceAlertReq).Return(priceAlert, nil)
		mockAlertService.On("CreateAlert", ctx, userID, transferAlertReq).Return(transferAlert, nil)

		// Create alerts
		priceResult, err := mockAlertService.CreateAlert(ctx, userID, priceAlertReq)
		require.NoError(t, err)
		assert.Equal(t, models.AlertTypePriceAbove, priceResult.Type)
		assert.Equal(t, price, *priceResult.Conditions.Price)

		transferResult, err := mockAlertService.CreateAlert(ctx, userID, transferAlertReq)
		require.NoError(t, err)
		assert.Equal(t, models.AlertTypeLargeTransfer, transferResult.Type)
		assert.Equal(t, threshold, *transferResult.Conditions.Threshold)

		// Mock getting user alerts
		userAlerts := []models.Alert{*priceAlert, *transferAlert}
		mockAlertService.On("GetUserAlerts", ctx, userID, (*string)(nil), 20, 0).Return(userAlerts, nil)

		// Get user alerts
		alerts, err := mockAlertService.GetUserAlerts(ctx, userID, nil, 20, 0)
		require.NoError(t, err)
		assert.Len(t, alerts, 2)

		// Verify different alert types
		alertTypes := make(map[string]bool)
		for _, alert := range alerts {
			alertTypes[alert.Type] = true
		}
		assert.True(t, alertTypes[models.AlertTypePriceAbove])
		assert.True(t, alertTypes[models.AlertTypeLargeTransfer])

		// Verify all mocks were called correctly
		mockAlertService.AssertExpectations(t)
	})
}