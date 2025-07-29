package repos

import (
	"context"
	"time"

	"github.com/defi-dashboard/backend/internal/models"
	"github.com/google/uuid"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	GetByAddress(ctx context.Context, address string) (*models.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, address, nonce string) (*models.User, error)
	UpdateNonce(ctx context.Context, address, nonce string) (*models.User, error)
	UpdateLastLogin(ctx context.Context, id uuid.UUID, lastLogin time.Time) error
	UpdateEmail(ctx context.Context, id uuid.UUID, email string) (*models.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// WalletRepository defines the interface for wallet data access
type WalletRepository interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Wallet, error)
	GetByAddress(ctx context.Context, address string, chainID int) (*models.Wallet, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Wallet, error)
	Create(ctx context.Context, userID uuid.UUID, address string, chainID int, label *string, isPrimary bool) (*models.Wallet, error)
	Update(ctx context.Context, id, userID uuid.UUID, label *string) (*models.Wallet, error)
	SetPrimary(ctx context.Context, userID, walletID uuid.UUID) error
	Delete(ctx context.Context, id, userID uuid.UUID) error
}

// TokenRepository defines the interface for token data access
type TokenRepository interface {
	GetByAddress(ctx context.Context, address string, chainID int) (*models.Token, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Token, error)
	GetByChainID(ctx context.Context, chainID int, limit, offset int) ([]*models.Token, error)
	Search(ctx context.Context, query string, chainID *int) ([]*models.Token, error)
	Create(ctx context.Context, token *models.Token) (*models.Token, error)
	UpdatePrice(ctx context.Context, address string, chainID int, priceUSD, priceChange24h, marketCap float64) (*models.Token, error)
}

// TransactionRepository defines the interface for transaction data access
type TransactionRepository interface {
	GetByHash(ctx context.Context, hash string) (*models.Transaction, error)
	GetUserTransactions(ctx context.Context, userID uuid.UUID, filters TransactionFilters) ([]*models.Transaction, error)
	GetWalletTransactions(ctx context.Context, address string, chainID int, txType *string, limit, offset int) ([]*models.Transaction, error)
	Create(ctx context.Context, tx *models.Transaction) (*models.Transaction, error)
	UpdateStatus(ctx context.Context, hash, status string, blockNumber, gasUsed int64, gasFeeUSD float64) (*models.Transaction, error)
	LinkToUser(ctx context.Context, userID, transactionID, walletID uuid.UUID) error
}

// TransactionFilters for querying transactions
type TransactionFilters struct {
	ChainID   *int
	Type      *string
	StartTime *string
	EndTime   *string
	Limit     int
	Offset    int
}

// ProtocolRepository defines the interface for protocol data access
type ProtocolRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.Protocol, error)
	GetBySlug(ctx context.Context, slug string) (*models.Protocol, error)
	GetAll(ctx context.Context, filters ProtocolFilters) ([]*models.Protocol, error)
	Count(ctx context.Context, filters ProtocolFilters) (int64, error)
	Create(ctx context.Context, protocol *models.Protocol) (*models.Protocol, error)
	Update(ctx context.Context, protocol *models.Protocol) (*models.Protocol, error)
	UpdateTVL(ctx context.Context, id uuid.UUID, tvlUSD float64) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByChain(ctx context.Context, chainID int) ([]*models.Protocol, error)
	GetWithPoolCount(ctx context.Context, limit, offset int) ([]*models.Protocol, error)
}

// YieldPoolRepository defines the interface for yield pool data access
type YieldPoolRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.YieldPool, error)
	GetByPoolID(ctx context.Context, poolID string) (*models.YieldPool, error)
	GetAll(ctx context.Context, filters YieldPoolFilters) ([]*models.YieldPool, error)
	Count(ctx context.Context, filters YieldPoolFilters) (int64, error)
	GetByProtocol(ctx context.Context, protocolID uuid.UUID, activeOnly bool) ([]*models.YieldPool, error)
	GetByChain(ctx context.Context, chainID int) ([]*models.YieldPool, error)
	GetTopByTVL(ctx context.Context, limit int) ([]*models.YieldPool, error)
	Upsert(ctx context.Context, pool *models.YieldPool) error
	UpdateAPY(ctx context.Context, poolID string, apy, apyBase, apyReward float64) error
	UpdateTVL(ctx context.Context, poolID string, tvlUSD float64) error
	Deactivate(ctx context.Context, poolID string) error
}

// YieldPositionRepository defines the interface for position data access
type YieldPositionRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.YieldPosition, error)
	GetByUser(ctx context.Context, userID uuid.UUID, filters PositionFilters) ([]*models.YieldPosition, error)
	GetByWallet(ctx context.Context, walletID uuid.UUID, activeOnly bool) ([]*models.YieldPosition, error)
	GetByPool(ctx context.Context, poolID uuid.UUID, activeOnly bool) ([]*models.YieldPosition, error)
	GetByProtocol(ctx context.Context, protocolID, userID uuid.UUID, activeOnly bool) ([]*models.YieldPosition, error)
	GetUserSummary(ctx context.Context, userID uuid.UUID) (*models.PositionSummary, error)
	GetUserPositionsWithPools(ctx context.Context, userID uuid.UUID, filters PositionFilters) ([]*models.YieldPosition, error)
	GetTopByValue(ctx context.Context, limit int) ([]*models.YieldPosition, error)
	Create(ctx context.Context, position *models.YieldPosition) (*models.YieldPosition, error)
	Update(ctx context.Context, position *models.YieldPosition) (*models.YieldPosition, error)
	UpdateBalance(ctx context.Context, id uuid.UUID, balanceRaw string, balanceUSD, currentValueUSD float64) error
	UpdateRewards(ctx context.Context, id uuid.UUID, pendingRewards, claimedRewards interface{}, totalRewardsUSD float64) error
	Close(ctx context.Context, id uuid.UUID, realizedPnLUSD float64) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateAllPnL(ctx context.Context) error
}

// ProtocolFilters for querying protocols
type ProtocolFilters struct {
	Category  *string
	IsActive  *bool
	RiskLevel *string
	SortBy    string
	Limit     int
	Offset    int
}

// YieldPoolFilters for querying yield pools
type YieldPoolFilters struct {
	Chain         *string
	ChainID       *int
	MinTVL        *float64
	MinAPY        *float64
	ProtocolSlug  *string
	RiskLevel     *string
	IsActive      *bool
	SortBy        string
	Limit         int
	Offset        int
}

// PositionFilters for querying positions
type PositionFilters struct {
	IsActive *bool
	ChainID  *int
	Limit    int
	Offset   int
}