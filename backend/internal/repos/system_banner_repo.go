package repos

import (
	"context"
	"fmt"

	"github.com/defi-dashboard/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SystemBannerRepository interface {
	GetAll(ctx context.Context, activeOnly bool) ([]models.SystemBanner, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.SystemBanner, error)
	Create(ctx context.Context, banner *models.SystemBanner) error
	Update(ctx context.Context, banner *models.SystemBanner) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type systemBannerRepository struct {
	db *pgxpool.Pool
}

func NewSystemBannerRepository(db *pgxpool.Pool) SystemBannerRepository {
	return &systemBannerRepository{db: db}
}

func (r *systemBannerRepository) GetAll(ctx context.Context, activeOnly bool) ([]models.SystemBanner, error) {
	query := `
		SELECT id, title, message, level, active, created_at, updated_at
		FROM system_banners
	`
	args := []interface{}{}

	if activeOnly {
		query += " WHERE active = $1"
		args = append(args, true)
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get system banners: %w", err)
	}
	defer rows.Close()

	var banners []models.SystemBanner
	for rows.Next() {
		var banner models.SystemBanner
		err := rows.Scan(
			&banner.ID,
			&banner.Title,
			&banner.Message,
			&banner.Level,
			&banner.Active,
			&banner.CreatedAt,
			&banner.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan system banner: %w", err)
		}
		banners = append(banners, banner)
	}

	return banners, rows.Err()
}

func (r *systemBannerRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.SystemBanner, error) {
	query := `
		SELECT id, title, message, level, active, created_at, updated_at
		FROM system_banners
		WHERE id = $1
	`

	var banner models.SystemBanner
	err := r.db.QueryRow(ctx, query, id).Scan(
		&banner.ID,
		&banner.Title,
		&banner.Message,
		&banner.Level,
		&banner.Active,
		&banner.CreatedAt,
		&banner.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("system banner not found")
		}
		return nil, fmt.Errorf("failed to get system banner: %w", err)
	}

	return &banner, nil
}

func (r *systemBannerRepository) Create(ctx context.Context, banner *models.SystemBanner) error {
	query := `
		INSERT INTO system_banners (title, message, level, active)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		banner.Title,
		banner.Message,
		banner.Level,
		banner.Active,
	).Scan(&banner.ID, &banner.CreatedAt, &banner.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create system banner: %w", err)
	}

	return nil
}

func (r *systemBannerRepository) Update(ctx context.Context, banner *models.SystemBanner) error {
	query := `
		UPDATE system_banners
		SET title = $2, message = $3, level = $4, active = $5, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`

	err := r.db.QueryRow(ctx, query,
		banner.ID,
		banner.Title,
		banner.Message,
		banner.Level,
		banner.Active,
	).Scan(&banner.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("system banner not found")
		}
		return fmt.Errorf("failed to update system banner: %w", err)
	}

	return nil
}

func (r *systemBannerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM system_banners WHERE id = $1`
	
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete system banner: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("system banner not found")
	}

	return nil
}