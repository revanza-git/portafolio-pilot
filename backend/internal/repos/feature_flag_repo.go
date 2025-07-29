package repos

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/defi-dashboard/backend/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FeatureFlagRepository interface {
	GetAll(ctx context.Context) ([]models.FeatureFlag, error)
	GetByName(ctx context.Context, name string) (*models.FeatureFlag, error)
	Upsert(ctx context.Context, flag *models.FeatureFlag) error
	Delete(ctx context.Context, name string) error
}

type featureFlagRepository struct {
	db *pgxpool.Pool
}

func NewFeatureFlagRepository(db *pgxpool.Pool) FeatureFlagRepository {
	return &featureFlagRepository{db: db}
}

func (r *featureFlagRepository) GetAll(ctx context.Context) ([]models.FeatureFlag, error) {
	query := `
		SELECT name, value, created_at, updated_at
		FROM feature_flags
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get feature flags: %w", err)
	}
	defer rows.Close()

	var flags []models.FeatureFlag
	for rows.Next() {
		var flag models.FeatureFlag
		var valueJSON []byte

		err := rows.Scan(
			&flag.Name,
			&valueJSON,
			&flag.CreatedAt,
			&flag.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan feature flag: %w", err)
		}

		// Unmarshal JSON value
		if err := json.Unmarshal(valueJSON, &flag.Value); err != nil {
			return nil, fmt.Errorf("failed to unmarshal flag value: %w", err)
		}

		flags = append(flags, flag)
	}

	return flags, rows.Err()
}

func (r *featureFlagRepository) GetByName(ctx context.Context, name string) (*models.FeatureFlag, error) {
	query := `
		SELECT name, value, created_at, updated_at
		FROM feature_flags
		WHERE name = $1
	`

	var flag models.FeatureFlag
	var valueJSON []byte

	err := r.db.QueryRow(ctx, query, name).Scan(
		&flag.Name,
		&valueJSON,
		&flag.CreatedAt,
		&flag.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get feature flag: %w", err)
	}

	// Unmarshal JSON value
	if err := json.Unmarshal(valueJSON, &flag.Value); err != nil {
		return nil, fmt.Errorf("failed to unmarshal flag value: %w", err)
	}

	return &flag, nil
}

func (r *featureFlagRepository) Upsert(ctx context.Context, flag *models.FeatureFlag) error {
	valueJSON, err := json.Marshal(flag.Value)
	if err != nil {
		return fmt.Errorf("failed to marshal flag value: %w", err)
	}

	query := `
		INSERT INTO feature_flags (name, value)
		VALUES ($1, $2)
		ON CONFLICT (name) 
		DO UPDATE SET 
			value = EXCLUDED.value,
			updated_at = NOW()
		RETURNING created_at, updated_at
	`

	err = r.db.QueryRow(ctx, query, flag.Name, valueJSON).Scan(&flag.CreatedAt, &flag.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to upsert feature flag: %w", err)
	}

	return nil
}

func (r *featureFlagRepository) Delete(ctx context.Context, name string) error {
	query := `DELETE FROM feature_flags WHERE name = $1`
	
	result, err := r.db.Exec(ctx, query, name)
	if err != nil {
		return fmt.Errorf("failed to delete feature flag: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("feature flag not found")
	}

	return nil
}