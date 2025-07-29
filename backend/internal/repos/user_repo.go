package repos

import (
	"context"
	"time"

	"github.com/defi-dashboard/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByAddress(ctx context.Context, address string) (*models.User, error) {
	query := `
		SELECT id, address, email, nonce, is_admin, last_login_at, created_at, updated_at
		FROM users 
		WHERE address = $1
	`
	
	var user models.User
	err := r.db.QueryRow(ctx, query, address).Scan(
		&user.ID, &user.Address, &user.Email, &user.Nonce, &user.IsAdmin,
		&user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, address, email, nonce, is_admin, last_login_at, created_at, updated_at
		FROM users 
		WHERE id = $1
	`
	
	var user models.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Address, &user.Email, &user.Nonce, &user.IsAdmin,
		&user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, address, nonce string) (*models.User, error) {
	query := `
		INSERT INTO users (address, nonce) 
		VALUES ($1, $2)
		RETURNING id, address, email, nonce, last_login_at, created_at, updated_at
	`
	
	var user models.User
	err := r.db.QueryRow(ctx, query, address, nonce).Scan(
		&user.ID, &user.Address, &user.Email, &user.Nonce, 
		&user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) UpdateNonce(ctx context.Context, address, nonce string) (*models.User, error) {
	query := `
		UPDATE users 
		SET nonce = $2, updated_at = NOW()
		WHERE address = $1
		RETURNING id, address, email, nonce, last_login_at, created_at, updated_at
	`
	
	var user models.User
	err := r.db.QueryRow(ctx, query, address, nonce).Scan(
		&user.ID, &user.Address, &user.Email, &user.Nonce, 
		&user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, address, email, nonce, last_login_at, created_at, updated_at
		FROM users 
		WHERE email = $1
	`
	
	var user models.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Address, &user.Email, &user.Nonce, 
		&user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID, lastLogin time.Time) error {
	query := `
		UPDATE users 
		SET last_login_at = $2, updated_at = NOW()
		WHERE id = $1
	`
	
	_, err := r.db.Exec(ctx, query, id, lastLogin)
	return err
}

func (r *userRepository) UpdateEmail(ctx context.Context, id uuid.UUID, email string) (*models.User, error) {
	query := `
		UPDATE users 
		SET email = $2, updated_at = NOW()
		WHERE id = $1
		RETURNING id, address, email, nonce, last_login_at, created_at, updated_at
	`
	
	var user models.User
	err := r.db.QueryRow(ctx, query, id, email).Scan(
		&user.ID, &user.Address, &user.Email, &user.Nonce, 
		&user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}