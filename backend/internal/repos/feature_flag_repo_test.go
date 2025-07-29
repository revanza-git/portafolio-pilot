package repos

import (
	"context"
	"testing"
	"time"

	"github.com/defi-dashboard/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockFeatureFlagRepo for testing
type MockFeatureFlagRepo struct {
	flags map[string]*models.FeatureFlag
}

func NewMockFeatureFlagRepo() *MockFeatureFlagRepo {
	return &MockFeatureFlagRepo{
		flags: make(map[string]*models.FeatureFlag),
	}
}

func (m *MockFeatureFlagRepo) GetAll(ctx context.Context) ([]models.FeatureFlag, error) {
	var flags []models.FeatureFlag
	for _, flag := range m.flags {
		flags = append(flags, *flag)
	}
	return flags, nil
}

func (m *MockFeatureFlagRepo) GetByName(ctx context.Context, name string) (*models.FeatureFlag, error) {
	flag, exists := m.flags[name]
	if !exists {
		return nil, assert.AnError
	}
	return flag, nil
}

func (m *MockFeatureFlagRepo) Upsert(ctx context.Context, flag *models.FeatureFlag) error {
	if flag.CreatedAt.IsZero() {
		flag.CreatedAt = time.Now()
	}
	flag.UpdatedAt = time.Now()
	m.flags[flag.Name] = flag
	return nil
}

func (m *MockFeatureFlagRepo) Delete(ctx context.Context, name string) error {
	_, exists := m.flags[name]
	if !exists {
		return assert.AnError
	}
	delete(m.flags, name)
	return nil
}

func TestFeatureFlagRepository_Upsert(t *testing.T) {
	repo := NewMockFeatureFlagRepo()
	ctx := context.Background()

	flag := &models.FeatureFlag{
		Name: "test-flag",
		Value: map[string]interface{}{
			"enabled": true,
			"rollout": 50,
		},
	}

	err := repo.Upsert(ctx, flag)
	require.NoError(t, err)
	assert.NotZero(t, flag.CreatedAt)
	assert.NotZero(t, flag.UpdatedAt)
}

func TestFeatureFlagRepository_GetAll(t *testing.T) {
	repo := NewMockFeatureFlagRepo()
	ctx := context.Background()

	flag1 := &models.FeatureFlag{
		Name: "flag1",
		Value: map[string]interface{}{
			"enabled": true,
		},
	}
	flag2 := &models.FeatureFlag{
		Name: "flag2",
		Value: map[string]interface{}{
			"enabled": false,
		},
	}

	repo.Upsert(ctx, flag1)
	repo.Upsert(ctx, flag2)

	flags, err := repo.GetAll(ctx)
	require.NoError(t, err)
	assert.Len(t, flags, 2)
}

func TestFeatureFlagRepository_GetByName(t *testing.T) {
	repo := NewMockFeatureFlagRepo()
	ctx := context.Background()

	flag := &models.FeatureFlag{
		Name: "test-flag",
		Value: map[string]interface{}{
			"enabled": true,
		},
	}

	repo.Upsert(ctx, flag)

	retrievedFlag, err := repo.GetByName(ctx, "test-flag")
	require.NoError(t, err)
	assert.Equal(t, "test-flag", retrievedFlag.Name)
	assert.Equal(t, true, retrievedFlag.Value["enabled"])
}

func TestFeatureFlagRepository_GetByName_NotFound(t *testing.T) {
	repo := NewMockFeatureFlagRepo()
	ctx := context.Background()

	_, err := repo.GetByName(ctx, "non-existent")
	assert.Error(t, err)
}

func TestFeatureFlagRepository_Delete(t *testing.T) {
	repo := NewMockFeatureFlagRepo()
	ctx := context.Background()

	flag := &models.FeatureFlag{
		Name: "test-flag",
		Value: map[string]interface{}{
			"enabled": true,
		},
	}

	repo.Upsert(ctx, flag)

	err := repo.Delete(ctx, "test-flag")
	require.NoError(t, err)

	_, err = repo.GetByName(ctx, "test-flag")
	assert.Error(t, err)
}

func TestFeatureFlagRepository_Delete_NotFound(t *testing.T) {
	repo := NewMockFeatureFlagRepo()
	ctx := context.Background()

	err := repo.Delete(ctx, "non-existent")
	assert.Error(t, err)
}

func TestFeatureFlagRepository_UpsertUpdate(t *testing.T) {
	repo := NewMockFeatureFlagRepo()
	ctx := context.Background()

	// Create initial flag
	flag := &models.FeatureFlag{
		Name: "test-flag",
		Value: map[string]interface{}{
			"enabled": false,
		},
	}
	repo.Upsert(ctx, flag)
	createdAt := flag.CreatedAt

	// Update the flag
	flag.Value["enabled"] = true
	err := repo.Upsert(ctx, flag)
	require.NoError(t, err)

	// Should have same created time but different updated time
	assert.Equal(t, createdAt, flag.CreatedAt)
	assert.True(t, flag.UpdatedAt.After(createdAt))

	// Verify the update
	retrievedFlag, err := repo.GetByName(ctx, "test-flag")
	require.NoError(t, err)
	assert.Equal(t, true, retrievedFlag.Value["enabled"])
}