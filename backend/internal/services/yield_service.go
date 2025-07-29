package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/defi-dashboard/backend/internal/models"
	"github.com/defi-dashboard/backend/internal/repos"
	"github.com/defi-dashboard/backend/pkg/errors"
	"github.com/google/uuid"
)

type YieldService struct {
	poolRepo     repos.YieldPoolRepository
	positionRepo repos.YieldPositionRepository
	protocolRepo repos.ProtocolRepository
	userRepo     repos.UserRepository
}

func NewYieldService(poolRepo repos.YieldPoolRepository, positionRepo repos.YieldPositionRepository, protocolRepo repos.ProtocolRepository, userRepo repos.UserRepository) *YieldService {
	return &YieldService{
		poolRepo:     poolRepo,
		positionRepo: positionRepo,
		protocolRepo: protocolRepo,
		userRepo:     userRepo,
	}
}

// Pool Management

func (s *YieldService) GetPools(ctx context.Context, filters repos.YieldPoolFilters) ([]*models.YieldPool, int64, error) {
	pools, err := s.poolRepo.GetAll(ctx, filters)
	if err != nil {
		return nil, 0, errors.Internal("Failed to fetch yield pools")
	}

	count, err := s.poolRepo.Count(ctx, filters)
	if err != nil {
		return nil, 0, errors.Internal("Failed to count yield pools")
	}

	return pools, count, nil
}

func (s *YieldService) GetPoolByID(ctx context.Context, poolID uuid.UUID) (*models.YieldPool, error) {
	pool, err := s.poolRepo.GetByID(ctx, poolID)
	if err != nil {
		return nil, errors.NotFound("Yield pool not found")
	}

	return pool, nil
}

func (s *YieldService) GetPoolsByProtocol(ctx context.Context, protocolID uuid.UUID, activeOnly bool) ([]*models.YieldPool, error) {
	pools, err := s.poolRepo.GetByProtocol(ctx, protocolID, activeOnly)
	if err != nil {
		return nil, errors.Internal("Failed to fetch pools by protocol")
	}

	return pools, nil
}

func (s *YieldService) GetPoolsByChain(ctx context.Context, chainID int) ([]*models.YieldPool, error) {
	pools, err := s.poolRepo.GetByChain(ctx, chainID)
	if err != nil {
		return nil, errors.Internal("Failed to fetch pools by chain")
	}

	return pools, nil
}

func (s *YieldService) GetTopPoolsByTVL(ctx context.Context, limit int) ([]*models.YieldPool, error) {
	pools, err := s.poolRepo.GetTopByTVL(ctx, limit)
	if err != nil {
		return nil, errors.Internal("Failed to fetch top pools by TVL")
	}

	return pools, nil
}

func (s *YieldService) RefreshPoolData(ctx context.Context, poolID string, tvlUSD, apy, apyBase, apyReward float64) error {
	// Update pool APY and TVL - this would typically be called by the worker
	if err := s.poolRepo.UpdateAPY(ctx, poolID, apy, apyBase, apyReward); err != nil {
		return errors.Internal("Failed to update pool APY")
	}

	if err := s.poolRepo.UpdateTVL(ctx, poolID, tvlUSD); err != nil {
		return errors.Internal("Failed to update pool TVL")
	}

	return nil
}

// Position Management

func (s *YieldService) GetUserPositions(ctx context.Context, userAddress string, filters repos.PositionFilters) (*models.PositionSummary, error) {
	// Get user by address
	user, err := s.userRepo.GetByAddress(ctx, userAddress)
	if err != nil {
		return nil, errors.NotFound("User not found")
	}

	// Get user positions with pool information
	positions, err := s.positionRepo.GetUserPositionsWithPools(ctx, user.ID, filters)
	if err != nil {
		return nil, errors.Internal("Failed to fetch user positions")
	}

	// Get position summary
	summary, err := s.positionRepo.GetUserSummary(ctx, user.ID)
	if err != nil {
		return nil, errors.Internal("Failed to fetch position summary")
	}

	// Set positions in summary
	summary.Positions = make([]models.YieldPosition, len(positions))
	for i, pos := range positions {
		summary.Positions[i] = *pos
	}

	return summary, nil
}

func (s *YieldService) GetPositionByID(ctx context.Context, positionID uuid.UUID) (*models.YieldPosition, error) {
	position, err := s.positionRepo.GetByID(ctx, positionID)
	if err != nil {
		return nil, errors.NotFound("Position not found")
	}

	return position, nil
}

func (s *YieldService) CreatePosition(ctx context.Context, userAddress string, req CreatePositionRequest) (*models.YieldPosition, error) {
	// Get user by address
	user, err := s.userRepo.GetByAddress(ctx, userAddress)
	if err != nil {
		return nil, errors.NotFound("User not found")
	}

	// Validate pool exists
	pool, err := s.poolRepo.GetByID(ctx, req.PoolID)
	if err != nil {
		return nil, errors.NotFound("Pool not found")
	}

	// Create position
	position := &models.YieldPosition{
		UserID:               user.ID,
		WalletID:             req.WalletID,
		PoolID:               req.PoolID,
		ProtocolID:           pool.ProtocolID,
		PositionID:           req.PositionID,
		PoolAddress:          req.PoolAddress,
		ChainID:              req.ChainID,
		BalanceRaw:           req.BalanceRaw,
		BalanceUSD:           req.BalanceUSD,
		EntryPriceUSD:        req.EntryPriceUSD,
		EntryBlockNumber:     req.EntryBlockNumber,
		EntryTransactionHash: req.EntryTransactionHash,
		EntryTime:            time.Now(),
		CurrentValueUSD:      req.BalanceUSD, // Initially same as balance
		IsActive:             true,
		Metadata:             req.Metadata,
	}

	createdPosition, err := s.positionRepo.Create(ctx, position)
	if err != nil {
		return nil, errors.Internal("Failed to create position")
	}

	return createdPosition, nil
}

func (s *YieldService) UpdatePosition(ctx context.Context, positionID uuid.UUID, req UpdatePositionRequest) (*models.YieldPosition, error) {
	// Get existing position
	position, err := s.positionRepo.GetByID(ctx, positionID)
	if err != nil {
		return nil, errors.NotFound("Position not found")
	}

	// Update balance if provided
	if req.BalanceRaw != nil {
		position.BalanceRaw = *req.BalanceRaw
	}
	if req.BalanceUSD != nil {
		position.BalanceUSD = req.BalanceUSD
	}
	if req.CurrentValueUSD != nil {
		position.CurrentValueUSD = req.CurrentValueUSD
	}

	// Update rewards if provided
	if req.PendingRewards != nil {
		position.PendingRewards = req.PendingRewards
	}
	if req.TotalRewardsUSD != nil {
		position.TotalRewardsUSD = req.TotalRewardsUSD
	}

	// Update metadata if provided
	if req.Metadata != nil {
		position.Metadata = req.Metadata
	}

	// Update last update time
	now := time.Now()
	position.LastUpdateTime = &now
	if req.LastUpdateBlock != nil {
		position.LastUpdateBlock = req.LastUpdateBlock
	}

	updatedPosition, err := s.positionRepo.Update(ctx, position)
	if err != nil {
		return nil, errors.Internal("Failed to update position")
	}

	return updatedPosition, nil
}

func (s *YieldService) ClosePosition(ctx context.Context, positionID uuid.UUID, realizedPnLUSD float64) error {
	if err := s.positionRepo.Close(ctx, positionID, realizedPnLUSD); err != nil {
		return errors.Internal("Failed to close position")
	}

	return nil
}

func (s *YieldService) ClaimRewards(ctx context.Context, userAddress string, positionID uuid.UUID) (*ClaimResponse, error) {
	// Get user by address
	user, err := s.userRepo.GetByAddress(ctx, userAddress)
	if err != nil {
		return nil, errors.NotFound("User not found")
	}

	// Get position and verify ownership
	position, err := s.positionRepo.GetByID(ctx, positionID)
	if err != nil {
		return nil, errors.NotFound("Position not found")
	}

	if position.UserID != user.ID {
		return nil, errors.Forbidden("Position does not belong to user")
	}

	// For now, return a mock transaction hash
	// In a real implementation, this would interact with the blockchain
	txHash := s.generateMockTxHash()

	// Update position to move pending rewards to claimed rewards
	if len(position.PendingRewards) > 0 {
		// Move pending to claimed
		if position.ClaimedRewards == nil {
			position.ClaimedRewards = []models.RewardInfo{}
		}

		// Add claim timestamp to rewards
		for _, reward := range position.PendingRewards {
			now := time.Now()
			reward.ClaimedAt = &now
			position.ClaimedRewards = append(position.ClaimedRewards, reward)
		}

		// Clear pending rewards
		position.PendingRewards = []models.RewardInfo{}

		// Update position
		_, err = s.positionRepo.Update(ctx, position)
		if err != nil {
			return nil, errors.Internal("Failed to update position after claim")
		}
	}

	return &ClaimResponse{
		TransactionHash: txHash,
		ClaimedAt:       time.Now(),
		Status:          "pending", // Would be "confirmed" after blockchain confirmation
	}, nil
}

func (s *YieldService) UpdateAllPositionsPnL(ctx context.Context) error {
	// This would typically be called by a background worker
	if err := s.positionRepo.UpdateAllPnL(ctx); err != nil {
		return errors.Internal("Failed to update positions P&L")
	}

	return nil
}

// Protocol Management

func (s *YieldService) GetProtocols(ctx context.Context, filters repos.ProtocolFilters) ([]*models.Protocol, int64, error) {
	protocols, err := s.protocolRepo.GetAll(ctx, filters)
	if err != nil {
		return nil, 0, errors.Internal("Failed to fetch protocols")
	}

	count, err := s.protocolRepo.Count(ctx, filters)
	if err != nil {
		return nil, 0, errors.Internal("Failed to count protocols")
	}

	return protocols, count, nil
}

func (s *YieldService) GetProtocolBySlug(ctx context.Context, slug string) (*models.Protocol, error) {
	protocol, err := s.protocolRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, errors.NotFound("Protocol not found")
	}

	return protocol, nil
}

// Helper methods

func (s *YieldService) generateMockTxHash() string {
	// Generate a mock transaction hash for development/testing
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return "0x" + hex.EncodeToString(bytes)
}

// Request/Response types

type CreatePositionRequest struct {
	PoolID               uuid.UUID   `json:"pool_id" validate:"required"`
	WalletID             uuid.UUID   `json:"wallet_id" validate:"required"`
	PositionID           *string     `json:"position_id,omitempty"`
	PoolAddress          *string     `json:"pool_address,omitempty"`
	ChainID              int         `json:"chain_id" validate:"required"`
	BalanceRaw           string      `json:"balance_raw" validate:"required"`
	BalanceUSD           *float64    `json:"balance_usd,omitempty"`
	EntryPriceUSD        *float64    `json:"entry_price_usd,omitempty"`
	EntryBlockNumber     *int64      `json:"entry_block_number,omitempty"`
	EntryTransactionHash *string     `json:"entry_transaction_hash,omitempty"`
	Metadata             interface{} `json:"metadata,omitempty"`
}

type UpdatePositionRequest struct {
	BalanceRaw        *string                `json:"balance_raw,omitempty"`
	BalanceUSD        *float64               `json:"balance_usd,omitempty"`
	CurrentValueUSD   *float64               `json:"current_value_usd,omitempty"`
	PendingRewards    []models.RewardInfo    `json:"pending_rewards,omitempty"`
	TotalRewardsUSD   *float64               `json:"total_rewards_usd,omitempty"`
	LastUpdateBlock   *int64                 `json:"last_update_block,omitempty"`
	Metadata          interface{}            `json:"metadata,omitempty"`
}

type ClaimResponse struct {
	TransactionHash string    `json:"transaction_hash"`
	ClaimedAt       time.Time `json:"claimed_at"`
	Status          string    `json:"status"`
}