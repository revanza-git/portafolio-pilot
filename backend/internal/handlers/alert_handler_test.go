package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/defi-dashboard/backend/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock AlertService for testing
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

func setupTestApp() *fiber.App {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		// Mock auth middleware - set userID in locals
		userID := uuid.New()
		c.Locals("userID", userID)
		return c.Next()
	})
	return app
}

func TestAlertHandler_CreateAlert(t *testing.T) {
	app := setupTestApp()
	mockService := new(MockAlertService)
	handler := NewAlertHandler(mockService)

	app.Post("/alerts", handler.CreateAlert)

	t.Run("Successful alert creation", func(t *testing.T) {
		price := 3000.0
		reqBody := models.CreateAlertRequest{
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

		expectedAlert := &models.Alert{
			ID:           uuid.New(),
			Type:         models.AlertTypePriceAbove,
			Status:       models.AlertStatusActive,
			Target:       reqBody.Target,
			Conditions:   reqBody.Conditions,
			Notification: reqBody.Notification,
		}

		// Setup mock
		mockService.On("CreateAlert", mock.Anything, mock.AnythingOfType("uuid.UUID"), &reqBody).Return(expectedAlert, nil)

		// Create request
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/alerts", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		// Execute request
		resp, err := app.Test(req)
		require.NoError(t, err)

		// Assertions
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response models.Alert
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, expectedAlert.Type, response.Type)
		assert.Equal(t, expectedAlert.Status, response.Status)

		// Verify mock
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid request body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/alerts", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Missing alert type", func(t *testing.T) {
		reqBody := models.CreateAlertRequest{
			// Type is missing
			Target: models.AlertTarget{
				Type:       "token",
				Identifier: "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
				ChainID:    1,
			},
		}

		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/alerts", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestAlertHandler_GetAlerts(t *testing.T) {
	app := setupTestApp()
	mockService := new(MockAlertService)
	handler := NewAlertHandler(mockService)

	app.Get("/alerts", handler.GetAlerts)

	t.Run("Get user alerts", func(t *testing.T) {
		price1 := 3000.0
		price2 := 2000.0

		expectedAlerts := []models.Alert{
			{
				ID:     uuid.New(),
				Type:   models.AlertTypePriceAbove,
				Status: models.AlertStatusActive,
				Conditions: models.AlertConditions{
					Price: &price1,
				},
			},
			{
				ID:     uuid.New(),
				Type:   models.AlertTypePriceBelow,
				Status: models.AlertStatusDisabled,
				Conditions: models.AlertConditions{
					Price: &price2,
				},
			},
		}

		// Setup mock
		mockService.On("GetUserAlerts", mock.Anything, mock.AnythingOfType("uuid.UUID"), (*string)(nil), 20, 0).Return(expectedAlerts, nil)

		// Execute request
		req := httptest.NewRequest("GET", "/alerts", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)

		// Assertions
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		
		data := response["data"].([]interface{})
		assert.Len(t, data, 2)

		meta := response["meta"].(map[string]interface{})
		assert.Equal(t, float64(20), meta["limit"])
		assert.Equal(t, float64(1), meta["page"])

		// Verify mock
		mockService.AssertExpectations(t)
	})

	t.Run("Get alerts with status filter", func(t *testing.T) {
		expectedAlerts := []models.Alert{
			{
				ID:     uuid.New(),
				Type:   models.AlertTypePriceAbove,
				Status: models.AlertStatusActive,
			},
		}

		status := models.AlertStatusActive
		mockService.On("GetUserAlerts", mock.Anything, mock.AnythingOfType("uuid.UUID"), &status, 20, 0).Return(expectedAlerts, nil)

		req := httptest.NewRequest("GET", "/alerts?status=active", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mockService.AssertExpectations(t)
	})
}

func TestAlertHandler_GetAlert(t *testing.T) {
	app := setupTestApp()
	mockService := new(MockAlertService)
	handler := NewAlertHandler(mockService)

	app.Get("/alerts/:alertId", handler.GetAlert)

	t.Run("Get specific alert", func(t *testing.T) {
		alertID := uuid.New()
		price := 3000.0

		expectedAlert := &models.Alert{
			ID:     alertID,
			Type:   models.AlertTypePriceAbove,
			Status: models.AlertStatusActive,
			Conditions: models.AlertConditions{
				Price: &price,
			},
		}

		// Setup mock
		mockService.On("GetAlert", mock.Anything, alertID, mock.AnythingOfType("uuid.UUID")).Return(expectedAlert, nil)

		// Execute request
		req := httptest.NewRequest("GET", fmt.Sprintf("/alerts/%s", alertID.String()), nil)
		resp, err := app.Test(req)
		require.NoError(t, err)

		// Assertions
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response models.Alert
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, alertID, response.ID)
		assert.Equal(t, expectedAlert.Type, response.Type)

		// Verify mock
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid alert ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/alerts/invalid-uuid", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestAlertHandler_UpdateAlert(t *testing.T) {
	app := setupTestApp()
	mockService := new(MockAlertService)
	handler := NewAlertHandler(mockService)

	app.Patch("/alerts/:alertId", handler.UpdateAlert)

	t.Run("Update alert successfully", func(t *testing.T) {
		alertID := uuid.New()
		price := 3500.0

		reqBody := models.UpdateAlertRequest{
			Conditions: &models.AlertConditions{
				Price: &price,
			},
		}

		expectedAlert := &models.Alert{
			ID:     alertID,
			Type:   models.AlertTypePriceAbove,
			Status: models.AlertStatusActive,
			Conditions: models.AlertConditions{
				Price: &price,
			},
		}

		// Setup mock
		mockService.On("UpdateAlert", mock.Anything, alertID, mock.AnythingOfType("uuid.UUID"), &reqBody).Return(expectedAlert, nil)

		// Create request
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/alerts/%s", alertID.String()), bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		// Execute request
		resp, err := app.Test(req)
		require.NoError(t, err)

		// Assertions
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response models.Alert
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, alertID, response.ID)
		assert.Equal(t, price, *response.Conditions.Price)

		// Verify mock
		mockService.AssertExpectations(t)
	})
}

func TestAlertHandler_DeleteAlert(t *testing.T) {
	app := setupTestApp()
	mockService := new(MockAlertService)
	handler := NewAlertHandler(mockService)

	app.Delete("/alerts/:alertId", handler.DeleteAlert)

	t.Run("Delete alert successfully", func(t *testing.T) {
		alertID := uuid.New()

		// Setup mock
		mockService.On("DeleteAlert", mock.Anything, alertID, mock.AnythingOfType("uuid.UUID")).Return(nil)

		// Execute request
		req := httptest.NewRequest("DELETE", fmt.Sprintf("/alerts/%s", alertID.String()), nil)
		resp, err := app.Test(req)
		require.NoError(t, err)

		// Assertions
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify mock
		mockService.AssertExpectations(t)
	})
}

func TestAlertHandler_PauseActivateAlert(t *testing.T) {
	app := setupTestApp()
	mockService := new(MockAlertService)
	handler := NewAlertHandler(mockService)

	app.Patch("/alerts/:alertId/pause", handler.PauseAlert)
	app.Patch("/alerts/:alertId/activate", handler.ActivateAlert)

	t.Run("Pause alert", func(t *testing.T) {
		alertID := uuid.New()

		expectedAlert := &models.Alert{
			ID:     alertID,
			Status: models.AlertStatusDisabled,
		}

		status := models.AlertStatusDisabled
		expectedReq := models.UpdateAlertRequest{
			Status: &status,
		}

		// Setup mock
		mockService.On("UpdateAlert", mock.Anything, alertID, mock.AnythingOfType("uuid.UUID"), &expectedReq).Return(expectedAlert, nil)

		// Execute request
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/alerts/%s/pause", alertID.String()), nil)
		resp, err := app.Test(req)
		require.NoError(t, err)

		// Assertions
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response models.Alert
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, models.AlertStatusDisabled, response.Status)

		// Verify mock
		mockService.AssertExpectations(t)
	})

	t.Run("Activate alert", func(t *testing.T) {
		alertID := uuid.New()

		expectedAlert := &models.Alert{
			ID:     alertID,
			Status: models.AlertStatusActive,
		}

		status := models.AlertStatusActive
		expectedReq := models.UpdateAlertRequest{
			Status: &status,
		}

		// Setup mock
		mockService.On("UpdateAlert", mock.Anything, alertID, mock.AnythingOfType("uuid.UUID"), &expectedReq).Return(expectedAlert, nil)

		// Execute request
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/alerts/%s/activate", alertID.String()), nil)
		resp, err := app.Test(req)
		require.NoError(t, err)

		// Assertions
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response models.Alert
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, models.AlertStatusActive, response.Status)

		// Verify mock
		mockService.AssertExpectations(t)
	})
}