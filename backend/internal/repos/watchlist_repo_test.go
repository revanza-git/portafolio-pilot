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

// MockWatchlistRepo for testing
type MockWatchlistRepo struct {
	watchlists map[uuid.UUID]*models.Watchlist
}

func NewMockWatchlistRepo() *MockWatchlistRepo {
	return &MockWatchlistRepo{
		watchlists: make(map[uuid.UUID]*models.Watchlist),
	}
}

func (m *MockWatchlistRepo) GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.Watchlist, error) {
	var watchlists []models.Watchlist
	for _, watchlist := range m.watchlists {
		if watchlist.UserID == userID {
			watchlists = append(watchlists, *watchlist)
		}
	}
	return watchlists, nil
}

func (m *MockWatchlistRepo) Create(ctx context.Context, watchlist *models.Watchlist) error {
	if watchlist.ID == uuid.Nil {
		watchlist.ID = uuid.New()
	}
	watchlist.CreatedAt = time.Now()
	watchlist.UpdatedAt = time.Now()
	m.watchlists[watchlist.ID] = watchlist
	return nil
}

func (m *MockWatchlistRepo) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	watchlist, exists := m.watchlists[id]
	if !exists || watchlist.UserID != userID {
		return assert.AnError
	}
	delete(m.watchlists, id)
	return nil
}

func (m *MockWatchlistRepo) ExistsByUserIDAndItem(ctx context.Context, userID uuid.UUID, itemType string, itemRefID int) (bool, error) {
	for _, watchlist := range m.watchlists {
		if watchlist.UserID == userID && watchlist.ItemType == itemType && watchlist.ItemRefID == itemRefID {
			return true, nil
		}
	}
	return false, nil
}

func TestWatchlistRepository_Create(t *testing.T) {
	repo := NewMockWatchlistRepo()
	ctx := context.Background()

	userID := uuid.New()
	watchlist := &models.Watchlist{
		UserID:    userID,
		ItemType:  models.WatchlistItemTypeToken,
		ItemRefID: 123,
	}

	err := repo.Create(ctx, watchlist)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, watchlist.ID)
	assert.NotZero(t, watchlist.CreatedAt)
	assert.NotZero(t, watchlist.UpdatedAt)
}

func TestWatchlistRepository_GetByUserID(t *testing.T) {
	repo := NewMockWatchlistRepo()
	ctx := context.Background()

	userID := uuid.New()
	watchlist1 := &models.Watchlist{
		ID:        uuid.New(),
		UserID:    userID,
		ItemType:  models.WatchlistItemTypeToken,
		ItemRefID: 123,
	}
	watchlist2 := &models.Watchlist{
		ID:        uuid.New(),
		UserID:    userID,
		ItemType:  models.WatchlistItemTypePool,
		ItemRefID: 456,
	}

	repo.Create(ctx, watchlist1)
	repo.Create(ctx, watchlist2)

	watchlists, err := repo.GetByUserID(ctx, userID)
	require.NoError(t, err)
	assert.Len(t, watchlists, 2)
}

func TestWatchlistRepository_Delete(t *testing.T) {
	repo := NewMockWatchlistRepo()
	ctx := context.Background()

	userID := uuid.New()
	watchlist := &models.Watchlist{
		ID:        uuid.New(),
		UserID:    userID,
		ItemType:  models.WatchlistItemTypeToken,
		ItemRefID: 123,
	}

	repo.Create(ctx, watchlist)

	err := repo.Delete(ctx, watchlist.ID, userID)
	require.NoError(t, err)

	watchlists, err := repo.GetByUserID(ctx, userID)
	require.NoError(t, err)
	assert.Len(t, watchlists, 0)
}

func TestWatchlistRepository_Delete_NotFound(t *testing.T) {
	repo := NewMockWatchlistRepo()
	ctx := context.Background()

	userID := uuid.New()
	nonExistentID := uuid.New()

	err := repo.Delete(ctx, nonExistentID, userID)
	assert.Error(t, err)
}

func TestWatchlistRepository_ExistsByUserIDAndItem(t *testing.T) {
	repo := NewMockWatchlistRepo()
	ctx := context.Background()

	userID := uuid.New()
	watchlist := &models.Watchlist{
		ID:        uuid.New(),
		UserID:    userID,
		ItemType:  models.WatchlistItemTypeToken,
		ItemRefID: 123,
	}

	repo.Create(ctx, watchlist)

	// Test existing item
	exists, err := repo.ExistsByUserIDAndItem(ctx, userID, models.WatchlistItemTypeToken, 123)
	require.NoError(t, err)
	assert.True(t, exists)

	// Test non-existing item
	exists, err = repo.ExistsByUserIDAndItem(ctx, userID, models.WatchlistItemTypeToken, 456)
	require.NoError(t, err)
	assert.False(t, exists)
}