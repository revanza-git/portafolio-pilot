package repos

import (
	"context"
	"fmt"

	"github.com/defi-dashboard/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WatchlistRepository interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.Watchlist, error)
	Create(ctx context.Context, watchlist *models.Watchlist) error
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	ExistsByUserIDAndItem(ctx context.Context, userID uuid.UUID, itemType string, itemRefID int) (bool, error)
}

type watchlistRepository struct {
	db *pgxpool.Pool
}

func NewWatchlistRepository(db *pgxpool.Pool) WatchlistRepository {
	return &watchlistRepository{db: db}
}

func (r *watchlistRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.Watchlist, error) {
	query := `
		SELECT id, user_id, item_type, item_ref_id, created_at, updated_at
		FROM watchlists
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get watchlists: %w", err)
	}
	defer rows.Close()

	var watchlists []models.Watchlist
	for rows.Next() {
		var watchlist models.Watchlist
		err := rows.Scan(
			&watchlist.ID,
			&watchlist.UserID,
			&watchlist.ItemType,
			&watchlist.ItemRefID,
			&watchlist.CreatedAt,
			&watchlist.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan watchlist: %w", err)
		}
		watchlists = append(watchlists, watchlist)
	}

	return watchlists, rows.Err()
}

func (r *watchlistRepository) Create(ctx context.Context, watchlist *models.Watchlist) error {
	query := `
		INSERT INTO watchlists (user_id, item_type, item_ref_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		watchlist.UserID,
		watchlist.ItemType,
		watchlist.ItemRefID,
	).Scan(&watchlist.ID, &watchlist.CreatedAt, &watchlist.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create watchlist item: %w", err)
	}

	return nil
}

func (r *watchlistRepository) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	query := `DELETE FROM watchlists WHERE id = $1 AND user_id = $2`
	
	result, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete watchlist item: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("watchlist item not found")
	}

	return nil
}

func (r *watchlistRepository) ExistsByUserIDAndItem(ctx context.Context, userID uuid.UUID, itemType string, itemRefID int) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM watchlists 
			WHERE user_id = $1 AND item_type = $2 AND item_ref_id = $3
		)
	`

	var exists bool
	err := r.db.QueryRow(ctx, query, userID, itemType, itemRefID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check watchlist item existence: %w", err)
	}

	return exists, nil
}