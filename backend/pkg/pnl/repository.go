package pnl

import (
	"context"
	"fmt"
	"time"

	"github.com/defi-dashboard/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	CreateLot(ctx context.Context, lot *models.PnLLot) error
	GetLotsByWallet(ctx context.Context, walletID uuid.UUID, tokenID uuid.UUID, from, to time.Time) ([]models.PnLLot, error)
	GetLotsByWalletAndToken(ctx context.Context, walletID uuid.UUID, tokenID uuid.UUID) ([]models.PnLLot, error)
	UpdateLotRemainingQuantity(ctx context.Context, lotID uuid.UUID, remainingQuantity string) error
	GetWalletTokens(ctx context.Context, walletID uuid.UUID) ([]uuid.UUID, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) CreateLot(ctx context.Context, lot *models.PnLLot) error {
	query := `
		INSERT INTO pnl_lots (
			id, wallet_id, token_id, transaction_hash, chain_id, type,
			quantity, price_usd, remaining_quantity, block_number, timestamp
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := r.db.Exec(ctx, query,
		lot.ID,
		lot.WalletID,
		lot.TokenID,
		lot.TransactionHash,
		lot.ChainID,
		lot.Type,
		lot.Quantity,
		lot.PriceUSD,
		lot.RemainingQuantity,
		lot.BlockNumber,
		lot.Timestamp,
	)

	return err
}

func (r *repository) GetLotsByWallet(ctx context.Context, walletID uuid.UUID, tokenID uuid.UUID, from, to time.Time) ([]models.PnLLot, error) {
	query := `
		SELECT 
			id, wallet_id, token_id, transaction_hash, chain_id, type,
			quantity, price_usd, remaining_quantity, block_number, timestamp,
			created_at, updated_at
		FROM pnl_lots 
		WHERE wallet_id = $1 AND token_id = $2 
		AND timestamp >= $3 AND timestamp <= $4
		ORDER BY timestamp ASC
	`

	rows, err := r.db.Query(ctx, query, walletID, tokenID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanLots(rows)
}

func (r *repository) GetLotsByWalletAndToken(ctx context.Context, walletID uuid.UUID, tokenID uuid.UUID) ([]models.PnLLot, error) {
	query := `
		SELECT 
			id, wallet_id, token_id, transaction_hash, chain_id, type,
			quantity, price_usd, remaining_quantity, block_number, timestamp,
			created_at, updated_at
		FROM pnl_lots 
		WHERE wallet_id = $1 AND token_id = $2
		ORDER BY timestamp ASC
	`

	rows, err := r.db.Query(ctx, query, walletID, tokenID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanLots(rows)
}

func (r *repository) UpdateLotRemainingQuantity(ctx context.Context, lotID uuid.UUID, remainingQuantity string) error {
	query := `
		UPDATE pnl_lots 
		SET remaining_quantity = $1, updated_at = NOW()
		WHERE id = $2
	`

	_, err := r.db.Exec(ctx, query, remainingQuantity, lotID)
	return err
}

func (r *repository) GetWalletTokens(ctx context.Context, walletID uuid.UUID) ([]uuid.UUID, error) {
	query := `
		SELECT DISTINCT token_id 
		FROM pnl_lots 
		WHERE wallet_id = $1
	`

	rows, err := r.db.Query(ctx, query, walletID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokenIDs []uuid.UUID
	for rows.Next() {
		var tokenID uuid.UUID
		if err := rows.Scan(&tokenID); err != nil {
			return nil, err
		}
		tokenIDs = append(tokenIDs, tokenID)
	}

	return tokenIDs, rows.Err()
}

func (r *repository) scanLots(rows pgx.Rows) ([]models.PnLLot, error) {
	var lots []models.PnLLot

	for rows.Next() {
		var lot models.PnLLot
		err := rows.Scan(
			&lot.ID,
			&lot.WalletID,
			&lot.TokenID,
			&lot.TransactionHash,
			&lot.ChainID,
			&lot.Type,
			&lot.Quantity,
			&lot.PriceUSD,
			&lot.RemainingQuantity,
			&lot.BlockNumber,
			&lot.Timestamp,
			&lot.CreatedAt,
			&lot.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan pnl lot: %w", err)
		}
		lots = append(lots, lot)
	}

	return lots, rows.Err()
}