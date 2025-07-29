package services

import (
	"context"
	"testing"
	"time"

	"github.com/defi-dashboard/backend/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock implementations for testing
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

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByAddress(ctx context.Context, address string) (*models.User, error) {
	args := m.Called(ctx, address)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func TestAlertService_CreateAlert(t *testing.T) {
	ctx := context.Background()
	
	mockAlertRepo := new(MockAlertRepository)
	mockUserRepo := new(MockUserRepository)
	service := NewAlertService(mockAlertRepo, mockUserRepo)

	userID := uuid.New()
	price := 3000.0

	user := &models.User{
		ID:      userID,
		Address: "0x1234567890123456789012345678901234567890",
	}

	req := &models.CreateAlertRequest{
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

	// Setup mocks
	mockUserRepo.On("GetByID", ctx, userID).Return(user, nil)
	mockAlertRepo.On("Create", ctx, mock.AnythingOfType("*models.Alert")).Return(nil)

	// Execute test
	alert, err := service.CreateAlert(ctx, userID, req)

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, userID, alert.UserID)
	assert.Equal(t, models.AlertTypePriceAbove, alert.Type)
	assert.Equal(t, models.AlertStatusActive, alert.Status)
	assert.Equal(t, price, *alert.Conditions.Price)
	assert.Equal(t, true, alert.Notification.Email)

	// Verify mocks
	mockUserRepo.AssertExpectations(t)
	mockAlertRepo.AssertExpectations(t)
}

func TestAlertService_CreateAlert_ValidationErrors(t *testing.T) {
	ctx := context.Background()
	
	mockAlertRepo := new(MockAlertRepository)
	mockUserRepo := new(MockUserRepository)
	service := NewAlertService(mockAlertRepo, mockUserRepo)

	userID := uuid.New()
	user := &models.User{ID: userID}

	t.Run("Invalid price condition", func(t *testing.T) {
		invalidPrice := -100.0
		req := &models.CreateAlertRequest{
			Type: models.AlertTypePriceAbove,
			Conditions: models.AlertConditions{
				Price: &invalidPrice,
			},
		}

		mockUserRepo.On("GetByID", ctx, userID).Return(user, nil)

		_, err := service.CreateAlert(ctx, userID, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "price must be specified and greater than 0")
	})

	t.Run("Missing price for price alert", func(t *testing.T) {
		req := &models.CreateAlertRequest{
			Type:       models.AlertTypePriceAbove,
			Conditions: models.AlertConditions{}, // No price specified
		}

		mockUserRepo.On("GetByID", ctx, userID).Return(user, nil)

		_, err := service.CreateAlert(ctx, userID, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "price must be specified")
	})

	t.Run("Invalid APR conditions", func(t *testing.T) {
		req := &models.CreateAlertRequest{
			Type:       models.AlertTypeAPRChange,
			Conditions: models.AlertConditions{}, // No APR conditions specified
		}

		mockUserRepo.On("GetByID", ctx, userID).Return(user, nil)

		_, err := service.CreateAlert(ctx, userID, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "either minAPR or maxAPR must be specified")
	})
}

func TestAlertService_GetAlert(t *testing.T) {
	ctx := context.Background()
	
	mockAlertRepo := new(MockAlertRepository)
	mockUserRepo := new(MockUserRepository)
	service := NewAlertService(mockAlertRepo, mockUserRepo)

	userID := uuid.New()
	alertID := uuid.New()
	price := 3000.0

	alert := &models.Alert{
		ID:     alertID,
		UserID: userID,
		Type:   models.AlertTypePriceAbove,
		Status: models.AlertStatusActive,
		Conditions: models.AlertConditions{
			Price: &price,
		},
	}

	// Setup mocks
	mockAlertRepo.On("GetByID", ctx, alertID).Return(alert, nil)

	// Execute test
	result, err := service.GetAlert(ctx, alertID, userID)

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, alertID, result.ID)
	assert.Equal(t, userID, result.UserID)

	// Verify mocks
	mockAlertRepo.AssertExpectations(t)
}

func TestAlertService_GetAlert_NotOwner(t *testing.T) {
	ctx := context.Background()
	
	mockAlertRepo := new(MockAlertRepository)
	mockUserRepo := new(MockUserRepository)
	service := NewAlertService(mockAlertRepo, mockUserRepo)

	userID := uuid.New()
	otherUserID := uuid.New()
	alertID := uuid.New()

	alert := &models.Alert{
		ID:     alertID,
		UserID: otherUserID, // Different user
		Type:   models.AlertTypePriceAbove,
		Status: models.AlertStatusActive,
	}

	// Setup mocks
	mockAlertRepo.On("GetByID", ctx, alertID).Return(alert, nil)

	// Execute test
	_, err := service.GetAlert(ctx, alertID, userID)

	// Assertions
	require.Error(t, err)
	assert.Contains(t, err.Error(), "alert not found")

	// Verify mocks
	mockAlertRepo.AssertExpectations(t)
}

func TestAlertService_UpdateAlert(t *testing.T) {
	ctx := context.Background()
	
	mockAlertRepo := new(MockAlertRepository)
	mockUserRepo := new(MockUserRepository)
	service := NewAlertService(mockAlertRepo, mockUserRepo)

	userID := uuid.New()
	alertID := uuid.New()
	price := 3000.0

	existingAlert := &models.Alert{
		ID:     alertID,
		UserID: userID,
		Type:   models.AlertTypePriceAbove,
		Status: models.AlertStatusActive,
		Conditions: models.AlertConditions{
			Price: &price,
		},
	}

	updatedStatus := models.AlertStatusDisabled
	newPrice := 3500.0
	req := &models.UpdateAlertRequest{
		Status: &updatedStatus,
		Conditions: &models.AlertConditions{
			Price: &newPrice,
		},
	}

	// Setup mocks
	mockAlertRepo.On("GetByID", ctx, alertID).Return(existingAlert, nil)
	mockAlertRepo.On("Update", ctx, mock.AnythingOfType("*models.Alert")).Return(nil)

	// Execute test
	result, err := service.UpdateAlert(ctx, alertID, userID, req)

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, models.AlertStatusDisabled, result.Status)
	assert.Equal(t, newPrice, *result.Conditions.Price)

	// Verify mocks
	mockAlertRepo.AssertExpectations(t)
}

func TestAlertService_TriggerAlert(t *testing.T) {
	ctx := context.Background()
	
	mockAlertRepo := new(MockAlertRepository)
	mockUserRepo := new(MockUserRepository)
	service := NewAlertService(mockAlertRepo, mockUserRepo)

	alertID := uuid.New()
	price := 3000.0

	alert := &models.Alert{
		ID:     alertID,
		UserID: uuid.New(),
		Type:   models.AlertTypePriceAbove,
		Status: models.AlertStatusActive,
		Conditions: models.AlertConditions{
			Price: &price,
		},
	}

	triggeredValue := map[string]interface{}{
		"currentPrice": 3100.0,
		"targetPrice":  3000.0,
	}

	// Setup mocks
	mockAlertRepo.On("GetByID", ctx, alertID).Return(alert, nil)
	mockAlertRepo.On("UpdateTriggered", ctx, alertID).Return(nil)
	mockAlertRepo.On("CreateHistory", ctx, mock.AnythingOfType("*models.AlertHistory")).Return(nil)

	// Execute test
	err := service.TriggerAlert(ctx, alertID, triggeredValue)

	// Assertions
	require.NoError(t, err)

	// Verify mocks
	mockAlertRepo.AssertExpectations(t)
}

func TestAlertService_DeleteAlert(t *testing.T) {
	ctx := context.Background()
	
	mockAlertRepo := new(MockAlertRepository)
	mockUserRepo := new(MockUserRepository)
	service := NewAlertService(mockAlertRepo, mockUserRepo)

	userID := uuid.New()
	alertID := uuid.New()

	alert := &models.Alert{
		ID:     alertID,
		UserID: userID,
		Type:   models.AlertTypePriceAbove,
		Status: models.AlertStatusActive,
	}

	// Setup mocks
	mockAlertRepo.On("GetByID", ctx, alertID).Return(alert, nil)
	mockAlertRepo.On("Delete", ctx, alertID).Return(nil)

	// Execute test
	err := service.DeleteAlert(ctx, alertID, userID)

	// Assertions
	require.NoError(t, err)

	// Verify mocks
	mockAlertRepo.AssertExpectations(t)
}

func TestAlertService_GetUserAlerts(t *testing.T) {
	ctx := context.Background()
	
	mockAlertRepo := new(MockAlertRepository)
	mockUserRepo := new(MockUserRepository)
	service := NewAlertService(mockAlertRepo, mockUserRepo)

	userID := uuid.New()
	price1 := 3000.0
	price2 := 2000.0

	alerts := []models.Alert{
		{
			ID:     uuid.New(),
			UserID: userID,
			Type:   models.AlertTypePriceAbove,
			Status: models.AlertStatusActive,
			Conditions: models.AlertConditions{
				Price: &price1,
			},
		},
		{
			ID:     uuid.New(),
			UserID: userID,
			Type:   models.AlertTypePriceBelow,
			Status: models.AlertStatusDisabled,
			Conditions: models.AlertConditions{
				Price: &price2,
			},
		},
	}

	status := models.AlertStatusActive

	// Setup mocks
	mockAlertRepo.On("GetByUserID", ctx, userID, &status, 20, 0).Return(alerts, nil)

	// Execute test
	result, err := service.GetUserAlerts(ctx, userID, &status, 20, 0)

	// Assertions
	require.NoError(t, err)
	assert.Len(t, result, 2)

	// Verify mocks
	mockAlertRepo.AssertExpectations(t)
}