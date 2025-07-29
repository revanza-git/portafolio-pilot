package repos

import (
	"context"
	"time"

	"github.com/defi-dashboard/backend/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NonceRepository interface {
	Store(ctx context.Context, address, nonce string, expiresAt time.Time) error
	ValidateAndUse(ctx context.Context, address, nonce string) (bool, error)
	CleanupExpired(ctx context.Context) error
	GetByAddressAndNonce(ctx context.Context, address, nonce string) (*models.NonceStorage, error)
}

type nonceRepository struct {
	db *pgxpool.Pool
}

func NewNonceRepository(db *pgxpool.Pool) NonceRepository {
	return &nonceRepository{db: db}
}

// Store saves a new nonce for the given address
func (r *nonceRepository) Store(ctx context.Context, address, nonce string, expiresAt time.Time) error {
	query := `
		INSERT INTO nonce_storage (address, nonce, expires_at)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.Exec(ctx, query, address, nonce, expiresAt)
	return err
}

// ValidateAndUse checks if nonce is valid and marks it as used
func (r *nonceRepository) ValidateAndUse(ctx context.Context, address, nonce string) (bool, error) {
	// Start transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)

	// Check if nonce exists and is valid
	var exists bool
	checkQuery := `
		SELECT EXISTS(
			SELECT 1 FROM nonce_storage 
			WHERE address = $1 AND nonce = $2 
			AND expires_at > NOW() AND used = FALSE
		)
	`
	err = tx.QueryRow(ctx, checkQuery, address, nonce).Scan(&exists)
	if err != nil {
		return false, err
	}

	if !exists {
		return false, nil
	}

	// Mark nonce as used
	updateQuery := `
		UPDATE nonce_storage 
		SET used = TRUE 
		WHERE address = $1 AND nonce = $2 AND used = FALSE
	`
	_, err = tx.Exec(ctx, updateQuery, address, nonce)
	if err != nil {
		return false, err
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		return false, err
	}

	return true, nil
}

// CleanupExpired removes expired and used nonces
func (r *nonceRepository) CleanupExpired(ctx context.Context) error {
	query := `
		DELETE FROM nonce_storage 
		WHERE expires_at < NOW() OR used = TRUE
	`
	_, err := r.db.Exec(ctx, query)
	return err
}

// GetByAddressAndNonce retrieves a specific nonce record
func (r *nonceRepository) GetByAddressAndNonce(ctx context.Context, address, nonce string) (*models.NonceStorage, error) {
	query := `
		SELECT id, address, nonce, expires_at, used, created_at
		FROM nonce_storage 
		WHERE address = $1 AND nonce = $2
	`
	
	var ns models.NonceStorage
	err := r.db.QueryRow(ctx, query, address, nonce).Scan(
		&ns.ID, &ns.Address, &ns.Nonce, &ns.ExpiresAt, &ns.Used, &ns.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &ns, nil
}