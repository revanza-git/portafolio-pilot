package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID          uuid.UUID  `json:"id"`
	Address     string     `json:"address"`
	Email       *string    `json:"email,omitempty"`
	Nonce       string     `json:"-"`
	IsAdmin     bool       `json:"is_admin"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// NonceStorage represents a nonce for SIWE authentication
type NonceStorage struct {
	ID        uuid.UUID `json:"id"`
	Address   string    `json:"address"`
	Nonce     string    `json:"nonce"`
	ExpiresAt time.Time `json:"expires_at"`
	Used      bool      `json:"used"`
	CreatedAt time.Time `json:"created_at"`
}

// Wallet represents a user's wallet
type Wallet struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Address   string    `json:"address"`
	ChainID   int       `json:"chain_id"`
	Label     *string   `json:"label,omitempty"`
	IsPrimary bool      `json:"is_primary"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Token represents a cryptocurrency token
type Token struct {
	ID            uuid.UUID `json:"id"`
	Address       string    `json:"address"`
	ChainID       int       `json:"chain_id"`
	Symbol        string    `json:"symbol"`
	Name          string    `json:"name"`
	Decimals      int       `json:"decimals"`
	LogoURI       *string   `json:"logo_uri,omitempty"`
	PriceUSD      *float64  `json:"price_usd,omitempty"`
	PriceChange24h *float64  `json:"price_change_24h,omitempty"`
	MarketCap     *float64  `json:"market_cap,omitempty"`
	TotalSupply   *string   `json:"total_supply,omitempty"`
	LastUpdated   *time.Time `json:"last_updated,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Balance represents a token balance for a wallet
type Balance struct {
	ID          uuid.UUID `json:"id"`
	WalletID    uuid.UUID `json:"wallet_id"`
	TokenID     uuid.UUID `json:"token_id"`
	Token       *Token    `json:"token,omitempty"`
	Balance     string    `json:"balance"`
	BalanceUSD  *float64  `json:"balance_usd,omitempty"`
	BlockNumber *int64    `json:"block_number,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Transaction represents a blockchain transaction
type Transaction struct {
	ID          uuid.UUID              `json:"id"`
	Hash        string                 `json:"hash"`
	ChainID     int                    `json:"chain_id"`
	FromAddress string                 `json:"from_address"`
	ToAddress   *string                `json:"to_address,omitempty"`
	Value       *string                `json:"value,omitempty"`
	GasUsed     *int64                 `json:"gas_used,omitempty"`
	GasPrice    *string                `json:"gas_price,omitempty"`
	GasFeeUSD   *float64               `json:"gas_fee_usd,omitempty"`
	BlockNumber *int64                 `json:"block_number,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
	Status      string                 `json:"status"`
	Type        string                 `json:"type"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// TokenAllowance represents a token approval/allowance
type TokenAllowance struct {
	ID              uuid.UUID  `json:"id"`
	WalletID        uuid.UUID  `json:"wallet_id"`
	TokenID         uuid.UUID  `json:"token_id"`
	Token           *Token     `json:"token,omitempty"`
	SpenderAddress  string     `json:"spender_address"`
	SpenderName     *string    `json:"spender_name,omitempty"`
	Allowance       string     `json:"allowance"`
	AllowanceUSD    *float64   `json:"allowance_usd,omitempty"`
	TransactionHash *string    `json:"transaction_hash,omitempty"`
	BlockNumber     *int64     `json:"block_number,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// Protocol represents a DeFi protocol
type Protocol struct {
	ID          uuid.UUID      `json:"id"`
	Name        string         `json:"name"`
	Slug        string         `json:"slug"`
	Description *string        `json:"description,omitempty"`
	WebsiteURL  *string        `json:"website_url,omitempty"`
	LogoURI     *string        `json:"logo_uri,omitempty"`
	Category    *string        `json:"category,omitempty"` // 'dex', 'lending', 'staking', 'yield_farming', etc.
	TotalTVLUSD *float64       `json:"total_tvl_usd,omitempty"`
	Chains      []int          `json:"chains,omitempty"`     // Supported chain IDs
	IsActive    bool           `json:"is_active"`
	RiskLevel   string         `json:"risk_level"`           // 'low', 'medium', 'high'
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// YieldPool represents a yield farming pool with enhanced information
type YieldPool struct {
	ID             uuid.UUID      `json:"id"`
	PoolID         string         `json:"pool_id"`
	ProtocolID     *uuid.UUID     `json:"protocol_id,omitempty"`
	Protocol       *Protocol      `json:"protocol,omitempty"`
	PoolName       string         `json:"pool_name"`
	ChainID        *int           `json:"chain_id,omitempty"`
	Chain          string         `json:"chain"`
	PoolAddress    *string        `json:"pool_address,omitempty"`
	Symbol         string         `json:"symbol"`
	TokenAddresses []string       `json:"token_addresses,omitempty"`
	
	// Financial metrics
	TVLUSD       *float64 `json:"tvl_usd,omitempty"`
	APY          *float64 `json:"apy,omitempty"`
	APYBase      *float64 `json:"apy_base,omitempty"`
	APYReward    *float64 `json:"apy_reward,omitempty"`
	FeesAPR      *float64 `json:"fees_apr,omitempty"`
	IL7D         *float64 `json:"il_7d,omitempty"` // Impermanent loss 7 days
	
	// Risk and limits
	RiskLevel      string   `json:"risk_level"`
	MinDepositUSD  *float64 `json:"min_deposit_usd,omitempty"`
	MaxDepositUSD  *float64 `json:"max_deposit_usd,omitempty"`
	
	// Status
	IsActive     bool        `json:"is_active"`
	StableCoin   bool        `json:"stable_coin"`
	
	// Additional data
	Metadata     interface{} `json:"metadata,omitempty"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

// TokenBalance represents a token balance in a position
type TokenBalance struct {
	TokenID    uuid.UUID `json:"token_id"`
	Token      *Token    `json:"token,omitempty"`
	Balance    string    `json:"balance"`
	BalanceUSD *float64  `json:"balance_usd,omitempty"`
}

// RewardInfo represents reward information
type RewardInfo struct {
	TokenID   uuid.UUID `json:"token_id"`
	Token     *Token    `json:"token,omitempty"`
	Amount    string    `json:"amount"`
	AmountUSD *float64  `json:"amount_usd,omitempty"`
	ClaimedAt *time.Time `json:"claimed_at,omitempty"`
}

// YieldPosition represents a user's position in a yield pool
type YieldPosition struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	WalletID   uuid.UUID `json:"wallet_id"`
	PoolID     uuid.UUID `json:"pool_id"`
	ProtocolID *uuid.UUID `json:"protocol_id,omitempty"`
	
	// Relations
	User     *User      `json:"user,omitempty"`
	Wallet   *Wallet    `json:"wallet,omitempty"`
	Pool     *YieldPool `json:"pool,omitempty"`
	Protocol *Protocol  `json:"protocol,omitempty"`
	
	// Position details
	PositionID            *string `json:"position_id,omitempty"` // External position ID
	PoolAddress           *string `json:"pool_address,omitempty"`
	ChainID               int     `json:"chain_id"`
	
	// Balance information
	BalanceRaw            string          `json:"balance_raw"`
	BalanceUSD            *float64        `json:"balance_usd,omitempty"`
	BalanceTokens         []TokenBalance  `json:"balance_tokens,omitempty"`
	
	// Entry information
	EntryPriceUSD         *float64  `json:"entry_price_usd,omitempty"`
	EntryBlockNumber      *int64    `json:"entry_block_number,omitempty"`
	EntryTransactionHash  *string   `json:"entry_transaction_hash,omitempty"`
	EntryTime             time.Time `json:"entry_time"`
	
	// Current status
	IsActive              bool       `json:"is_active"`
	LastUpdateBlock       *int64     `json:"last_update_block,omitempty"`
	LastUpdateTime        *time.Time `json:"last_update_time,omitempty"`
	
	// Rewards information
	PendingRewards        []RewardInfo `json:"pending_rewards,omitempty"`
	ClaimedRewards        []RewardInfo `json:"claimed_rewards,omitempty"`
	TotalRewardsUSD       *float64     `json:"total_rewards_usd,omitempty"`
	
	// P&L calculation
	CurrentValueUSD       *float64 `json:"current_value_usd,omitempty"`
	UnrealizedPnLUSD      *float64 `json:"unrealized_pnl_usd,omitempty"`
	RealizedPnLUSD        *float64 `json:"realized_pnl_usd,omitempty"`
	TotalFeesPaidUSD      *float64 `json:"total_fees_paid_usd,omitempty"`
	
	// Calculated fields
	PnLPercentage         *float64 `json:"pnl_percentage,omitempty"`
	
	// Additional data
	Metadata              interface{} `json:"metadata,omitempty"`
	CreatedAt             time.Time   `json:"created_at"`
	UpdatedAt             time.Time   `json:"updated_at"`
}

// PositionSummary represents aggregated position information for a user
type PositionSummary struct {
	TotalValueUSD       float64         `json:"total_value_usd"`
	TotalPnLUSD         float64         `json:"total_pnl_usd"`
	TotalPnLPercentage  float64         `json:"total_pnl_percentage"`
	TotalRewardsUSD     float64         `json:"total_rewards_usd"`
	ActivePositions     int             `json:"active_positions"`
	Positions           []YieldPosition `json:"positions"`
}

// PnLLot represents a buy/sell lot for FIFO/LIFO PnL calculations
type PnLLot struct {
	ID                uuid.UUID `json:"id"`
	WalletID          uuid.UUID `json:"wallet_id"`
	TokenID           uuid.UUID `json:"token_id"`
	TransactionHash   string    `json:"transaction_hash"`
	ChainID           int       `json:"chain_id"`
	Type              string    `json:"type"` // 'buy' or 'sell'
	Quantity          string    `json:"quantity"`
	PriceUSD          string    `json:"price_usd"`
	RemainingQuantity string    `json:"remaining_quantity"`
	BlockNumber       int64     `json:"block_number"`
	Timestamp         time.Time `json:"timestamp"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// PnLCalculation represents the result of a PnL calculation
type PnLCalculation struct {
	WalletAddress     string               `json:"wallet_address"`
	TokenAddress      string               `json:"token_address"`
	TokenSymbol       string               `json:"token_symbol"`
	Method            string               `json:"method"` // 'fifo' or 'lifo'
	RealizedPnLUSD    string               `json:"realized_pnl_usd"`
	UnrealizedPnLUSD  string               `json:"unrealized_pnl_usd"`
	TotalPnLUSD       string               `json:"total_pnl_usd"`
	TotalCostBasisUSD string               `json:"total_cost_basis_usd"`
	CurrentValueUSD   string               `json:"current_value_usd"`
	CurrentQuantity   string               `json:"current_quantity"`
	Lots              []PnLLot             `json:"lots"`
	CalculatedAt      time.Time            `json:"calculated_at"`
}

// PnLExportData represents data structure for CSV export
type PnLExportData struct {
	WalletAddress     string    `csv:"wallet_address"`
	TokenSymbol       string    `csv:"token_symbol"`
	TokenAddress      string    `csv:"token_address"`
	TransactionHash   string    `csv:"transaction_hash"`
	Type              string    `csv:"type"`
	Quantity          string    `csv:"quantity"`
	PriceUSD          string    `csv:"price_usd"`
	RemainingQuantity string    `csv:"remaining_quantity"`
	RealizedPnLUSD    string    `csv:"realized_pnl_usd"`
	Timestamp         time.Time `csv:"timestamp"`
	BlockNumber       int64     `csv:"block_number"`
}

// Alert represents an alert configuration
type Alert struct {
	ID                uuid.UUID       `json:"id"`
	UserID            uuid.UUID       `json:"user_id"`
	Type              string          `json:"type"`
	Status            string          `json:"status"`
	Target            AlertTarget     `json:"target"`
	Conditions        AlertConditions `json:"conditions"`
	Notification      AlertNotification `json:"notification"`
	LastTriggeredAt   *time.Time      `json:"last_triggered_at,omitempty"`
	TriggerCount      int             `json:"trigger_count"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

// AlertTarget represents the target entity for an alert
type AlertTarget struct {
	Type       string `json:"type"`        // token, address, pool
	Identifier string `json:"identifier"`  // token address, wallet address, pool ID
	ChainID    int    `json:"chainId"`
}

// AlertConditions represents the conditions that trigger an alert
type AlertConditions struct {
	// Price alerts
	Price         *float64 `json:"price,omitempty"`
	
	// Transfer alerts  
	Threshold     *string  `json:"threshold,omitempty"` // Wei amount
	
	// Liquidity alerts
	ChangePercent *float64 `json:"changePercent,omitempty"`
	
	// APR alerts
	MinAPR        *float64 `json:"minAPR,omitempty"`
	MaxAPR        *float64 `json:"maxAPR,omitempty"`
}

// AlertNotification represents notification preferences
type AlertNotification struct {
	Email   bool   `json:"email"`
	Webhook string `json:"webhook,omitempty"`
}

// AlertHistory represents a triggered alert event  
type AlertHistory struct {
	ID                  uuid.UUID               `json:"id"`
	AlertID             uuid.UUID               `json:"alert_id"`
	TriggeredAt         time.Time               `json:"triggered_at"`
	ConditionsSnapshot  AlertConditions         `json:"conditions_snapshot"`
	TriggeredValue      map[string]interface{}  `json:"triggered_value"`
	NotificationSent    bool                    `json:"notification_sent"`
	NotificationError   *string                 `json:"notification_error,omitempty"`
}

// Alert type constants
const (
	AlertTypePriceAbove      = "price_above"
	AlertTypePriceBelow      = "price_below" 
	AlertTypeLargeTransfer   = "large_transfer"
	AlertTypeApproval        = "approval"
	AlertTypeLiquidityChange = "liquidity_change"
	AlertTypeAPRChange       = "apr_change"
)

// Alert status constants
const (
	AlertStatusActive    = "active"
	AlertStatusTriggered = "triggered"
	AlertStatusExpired   = "expired"
	AlertStatusDisabled  = "disabled"
)

// CreateAlertRequest represents the request to create an alert
type CreateAlertRequest struct {
	Type         string            `json:"type" validate:"required,oneof=price_above price_below large_transfer approval liquidity_change apr_change"`
	Target       AlertTarget       `json:"target" validate:"required"`
	Conditions   AlertConditions   `json:"conditions" validate:"required"`
	Notification AlertNotification `json:"notification" validate:"required"`
}

// UpdateAlertRequest represents the request to update an alert
type UpdateAlertRequest struct {
	Status       *string           `json:"status,omitempty" validate:"omitempty,oneof=active disabled"`
	Conditions   *AlertConditions  `json:"conditions,omitempty"`
	Notification *AlertNotification `json:"notification,omitempty"`
}

// Watchlist represents a user's watchlist item
type Watchlist struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	ItemType   string    `json:"item_type"`
	ItemRefID  int       `json:"item_ref_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Watchlist item type constants
const (
	WatchlistItemTypeToken    = "token"
	WatchlistItemTypePool     = "pool"
	WatchlistItemTypeProtocol = "protocol"
)

// CreateWatchlistRequest represents the request to create a watchlist item
type CreateWatchlistRequest struct {
	ItemType  string `json:"item_type" validate:"required,oneof=token pool protocol"`
	ItemRefID int    `json:"item_ref_id" validate:"required,min=1"`
}

// FeatureFlag represents a feature flag configuration
type FeatureFlag struct {
	Name      string                 `json:"name"`
	Value     map[string]interface{} `json:"value"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// SystemBanner represents a system-wide banner notification
type SystemBanner struct {
	ID        uuid.UUID `json:"id"`
	Title     *string   `json:"title,omitempty"`
	Message   string    `json:"message"`
	Level     string    `json:"level"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Banner level constants
const (
	BannerLevelInfo    = "info"
	BannerLevelWarning = "warning"
	BannerLevelError   = "error"
	BannerLevelSuccess = "success"
)

// CreateFeatureFlagRequest represents the request to create/update a feature flag
type CreateFeatureFlagRequest struct {
	Name  string                 `json:"name" validate:"required"`
	Value map[string]interface{} `json:"value" validate:"required"`
}

// CreateSystemBannerRequest represents the request to create a system banner
type CreateSystemBannerRequest struct {
	Title   *string `json:"title,omitempty"`
	Message string  `json:"message" validate:"required"`
	Level   string  `json:"level" validate:"required,oneof=info warning error success"`
	Active  bool    `json:"active"`
}

// UpdateSystemBannerRequest represents the request to update a system banner
type UpdateSystemBannerRequest struct {
	Title   *string `json:"title,omitempty"`
	Message *string `json:"message,omitempty"`
	Level   *string `json:"level,omitempty" validate:"omitempty,oneof=info warning error success"`
	Active  *bool   `json:"active,omitempty"`
}