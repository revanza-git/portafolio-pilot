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

// Mock WatchlistRepository for testing
type MockWatchlistRepository struct {
	mock.Mock
}

func (m *MockWatchlistRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.Watchlist, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Watchlist), args.Error(1)
}

func (m *MockWatchlistRepository) Create(ctx context.Context, watchlist *models.Watchlist) error {
	args := m.Called(ctx, watchlist)
	return args.Error(0)
}

func (m *MockWatchlistRepository) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func (m *MockWatchlistRepository) ExistsByUserIDAndItem(ctx context.Context, userID uuid.UUID, itemType string, itemRefID int) (bool, error) {
	args := m.Called(ctx, userID, itemType, itemRefID)
	return args.Get(0).(bool), args.Error(1)
}

func createTestWatchlistHandler() (*WatchlistHandler, *MockWatchlistRepository) {
	mockRepo := new(MockWatchlistRepository)
	handler := NewWatchlistHandler(mockRepo)
	return handler, mockRepo
}

func createTestApp(handler *WatchlistHandler) *fiber.App {
	app := fiber.New()
	
	// Setup middleware to inject userID for testing
	app.Use(func(c *fiber.Ctx) error {
		userID := uuid.New()
		c.Locals("userID", userID)
		return c.Next()
	})
	
	watchlist := app.Group("/watchlist")
	watchlist.Get("/", handler.GetWatchlist)
	watchlist.Post("/", handler.CreateWatchlistItem)
	watchlist.Delete("/:id", handler.DeleteWatchlistItem)
	
	return app
}

func TestGetWatchlist_Success(t *testing.T) {
	handler, mockRepo := createTestWatchlistHandler()
	app := createTestApp(handler)

	userID := uuid.New()
	watchlists := []models.Watchlist{
		{
			ID:        uuid.New(),
			UserID:    userID,
			ItemType:  models.WatchlistItemTypeToken,
			ItemRefID: 123,
		},
	}

	mockRepo.On("GetByUserID", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(watchlists, nil)

	req := httptest.NewRequest("GET", "/watchlist", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result []models.Watchlist
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, models.WatchlistItemTypeToken, result[0].ItemType)
}

func TestCreateWatchlistItem_Success(t *testing.T) {
	handler, mockRepo := createTestWatchlistHandler()
	app := createTestApp(handler)

	request := models.CreateWatchlistRequest{
		ItemType:  models.WatchlistItemTypeToken,
		ItemRefID: 123,
	}

	mockRepo.On("ExistsByUserIDAndItem", mock.Anything, mock.AnythingOfType("uuid.UUID"), models.WatchlistItemTypeToken, 123).Return(false, nil)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Watchlist")).Return(nil)

	reqBody, _ := json.Marshal(request)
	req := httptest.NewRequest("POST", "/watchlist", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var result models.Watchlist
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	assert.Equal(t, models.WatchlistItemTypeToken, result.ItemType)
	assert.Equal(t, 123, result.ItemRefID)
}

func TestCreateWatchlistItem_AlreadyExists(t *testing.T) {
	handler, mockRepo := createTestWatchlistHandler()
	app := createTestApp(handler)

	request := models.CreateWatchlistRequest{
		ItemType:  models.WatchlistItemTypeToken,
		ItemRefID: 123,
	}

	mockRepo.On("ExistsByUserIDAndItem", mock.Anything, mock.AnythingOfType("uuid.UUID"), models.WatchlistItemTypeToken, 123).Return(true, nil)

	reqBody, _ := json.Marshal(request)
	req := httptest.NewRequest("POST", "/watchlist", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestCreateWatchlistItem_InvalidItemType(t *testing.T) {
	handler, _ := createTestWatchlistHandler()
	app := createTestApp(handler)

	request := models.CreateWatchlistRequest{
		ItemType:  "invalid",
		ItemRefID: 123,
	}

	reqBody, _ := json.Marshal(request)
	req := httptest.NewRequest("POST", "/watchlist", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestCreateWatchlistItem_InvalidItemRefID(t *testing.T) {
	handler, _ := createTestWatchlistHandler()
	app := createTestApp(handler)

	request := models.CreateWatchlistRequest{
		ItemType:  models.WatchlistItemTypeToken,
		ItemRefID: 0,
	}

	reqBody, _ := json.Marshal(request)
	req := httptest.NewRequest("POST", "/watchlist", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestDeleteWatchlistItem_Success(t *testing.T) {
	handler, mockRepo := createTestWatchlistHandler()
	app := createTestApp(handler)

	watchlistID := uuid.New()

	mockRepo.On("Delete", mock.Anything, watchlistID, mock.AnythingOfType("uuid.UUID")).Return(nil)

	req := httptest.NewRequest("DELETE", "/watchlist/"+watchlistID.String(), nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestDeleteWatchlistItem_InvalidID(t *testing.T) {
	handler, _ := createTestWatchlistHandler()
	app := createTestApp(handler)

	req := httptest.NewRequest("DELETE", "/watchlist/invalid-uuid", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}