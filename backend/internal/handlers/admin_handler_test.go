package handlers

import (
	"bytes"
	"context"
	"encoding/json"
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

// Mock repositories for testing
type MockFeatureFlagRepository struct {
	mock.Mock
}

func (m *MockFeatureFlagRepository) GetAll(ctx context.Context) ([]models.FeatureFlag, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.FeatureFlag), args.Error(1)
}

func (m *MockFeatureFlagRepository) GetByName(ctx context.Context, name string) (*models.FeatureFlag, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(*models.FeatureFlag), args.Error(1)
}

func (m *MockFeatureFlagRepository) Upsert(ctx context.Context, flag *models.FeatureFlag) error {
	args := m.Called(ctx, flag)
	return args.Error(0)
}

func (m *MockFeatureFlagRepository) Delete(ctx context.Context, name string) error {
	args := m.Called(ctx, name)
	return args.Error(0)
}

type MockSystemBannerRepository struct {
	mock.Mock
}

func (m *MockSystemBannerRepository) GetAll(ctx context.Context, activeOnly bool) ([]models.SystemBanner, error) {
	args := m.Called(ctx, activeOnly)
	return args.Get(0).([]models.SystemBanner), args.Error(1)
}

func (m *MockSystemBannerRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.SystemBanner, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.SystemBanner), args.Error(1)
}

func (m *MockSystemBannerRepository) Create(ctx context.Context, banner *models.SystemBanner) error {
	args := m.Called(ctx, banner)
	return args.Error(0)
}

func (m *MockSystemBannerRepository) Update(ctx context.Context, banner *models.SystemBanner) error {
	args := m.Called(ctx, banner)
	return args.Error(0)
}

func (m *MockSystemBannerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func createTestAdminHandler() (*AdminHandler, *MockUserRepository, *MockFeatureFlagRepository, *MockSystemBannerRepository) {
	mockUserRepo := new(MockUserRepository)
	mockFlagRepo := new(MockFeatureFlagRepository)
	mockBannerRepo := new(MockSystemBannerRepository)
	handler := NewAdminHandler(mockUserRepo, mockFlagRepo, mockBannerRepo)
	return handler, mockUserRepo, mockFlagRepo, mockBannerRepo
}

func createTestAdminApp(handler *AdminHandler) *fiber.App {
	app := fiber.New()
	
	admin := app.Group("/admin")
	admin.Get("/feature-flags", handler.GetFeatureFlags)
	admin.Post("/feature-flags", handler.CreateFeatureFlag)
	admin.Get("/banners", handler.GetSystemBanners)
	admin.Post("/banners", handler.CreateSystemBanner)
	
	return app
}

func TestGetFeatureFlags_Success(t *testing.T) {
	handler, _, mockFlagRepo, _ := createTestAdminHandler()
	app := createTestAdminApp(handler)

	flags := []models.FeatureFlag{
		{
			Name: "test-flag",
			Value: map[string]interface{}{
				"enabled": true,
			},
		},
	}

	mockFlagRepo.On("GetAll", mock.Anything).Return(flags, nil)

	req := httptest.NewRequest("GET", "/admin/feature-flags", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result []models.FeatureFlag
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "test-flag", result[0].Name)
}

func TestCreateFeatureFlag_Success(t *testing.T) {
	handler, _, mockFlagRepo, _ := createTestAdminHandler()
	app := createTestAdminApp(handler)

	request := models.CreateFeatureFlagRequest{
		Name: "new-flag",
		Value: map[string]interface{}{
			"enabled": true,
		},
	}

	mockFlagRepo.On("Upsert", mock.Anything, mock.AnythingOfType("*models.FeatureFlag")).Return(nil)

	reqBody, _ := json.Marshal(request)
	req := httptest.NewRequest("POST", "/admin/feature-flags", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var result models.FeatureFlag
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	assert.Equal(t, "new-flag", result.Name)
	assert.Equal(t, true, result.Value["enabled"])
}

func TestCreateFeatureFlag_InvalidRequest(t *testing.T) {
	handler, _, _, _ := createTestAdminHandler()
	app := createTestAdminApp(handler)

	request := models.CreateFeatureFlagRequest{
		Name: "", // Empty name should fail
		Value: map[string]interface{}{
			"enabled": true,
		},
	}

	reqBody, _ := json.Marshal(request)
	req := httptest.NewRequest("POST", "/admin/feature-flags", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestGetSystemBanners_Success(t *testing.T) {
	handler, _, _, mockBannerRepo := createTestAdminHandler()
	app := createTestAdminApp(handler)

	banners := []models.SystemBanner{
		{
			ID:      uuid.New(),
			Message: "Test banner",
			Level:   models.BannerLevelInfo,
			Active:  true,
		},
	}

	mockBannerRepo.On("GetAll", mock.Anything, false).Return(banners, nil)

	req := httptest.NewRequest("GET", "/admin/banners", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result []models.SystemBanner
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Test banner", result[0].Message)
}

func TestGetSystemBanners_ActiveOnly(t *testing.T) {
	handler, _, _, mockBannerRepo := createTestAdminHandler()
	app := createTestAdminApp(handler)

	banners := []models.SystemBanner{
		{
			ID:      uuid.New(),
			Message: "Active banner",
			Level:   models.BannerLevelInfo,
			Active:  true,
		},
	}

	mockBannerRepo.On("GetAll", mock.Anything, true).Return(banners, nil)

	req := httptest.NewRequest("GET", "/admin/banners?active=true", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestCreateSystemBanner_Success(t *testing.T) {
	handler, _, _, mockBannerRepo := createTestAdminHandler()
	app := createTestAdminApp(handler)

	request := models.CreateSystemBannerRequest{
		Message: "New banner",
		Level:   models.BannerLevelWarning,
		Active:  true,
	}

	mockBannerRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.SystemBanner")).Return(nil)

	reqBody, _ := json.Marshal(request)
	req := httptest.NewRequest("POST", "/admin/banners", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var result models.SystemBanner
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	assert.Equal(t, "New banner", result.Message)
	assert.Equal(t, models.BannerLevelWarning, result.Level)
}

func TestCreateSystemBanner_InvalidLevel(t *testing.T) {
	handler, _, _, _ := createTestAdminHandler()
	app := createTestAdminApp(handler)

	request := models.CreateSystemBannerRequest{
		Message: "New banner",
		Level:   "invalid", // Invalid level
		Active:  true,
	}

	reqBody, _ := json.Marshal(request)
	req := httptest.NewRequest("POST", "/admin/banners", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestCreateSystemBanner_EmptyMessage(t *testing.T) {
	handler, _, _, _ := createTestAdminHandler()
	app := createTestAdminApp(handler)

	request := models.CreateSystemBannerRequest{
		Message: "", // Empty message should fail
		Level:   models.BannerLevelInfo,
		Active:  true,
	}

	reqBody, _ := json.Marshal(request)
	req := httptest.NewRequest("POST", "/admin/banners", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}