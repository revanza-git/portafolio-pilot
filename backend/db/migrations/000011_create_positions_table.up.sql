-- Create positions table for tracking user yield positions
CREATE TABLE IF NOT EXISTS yield_positions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    wallet_id UUID NOT NULL REFERENCES wallets(id) ON DELETE CASCADE,
    pool_id UUID NOT NULL REFERENCES yield_pools(id) ON DELETE CASCADE,
    protocol_id UUID REFERENCES protocols(id) ON DELETE SET NULL,
    
    -- Position details
    position_id VARCHAR(255), -- External position identifier (e.g., NFT token ID for Uniswap V3)
    pool_address VARCHAR(42), -- Contract address of the pool
    chain_id INTEGER NOT NULL,
    
    -- Balance information
    balance_raw VARCHAR(78), -- Raw balance (can be very large numbers)
    balance_usd DECIMAL(20, 8), -- USD value of position
    balance_tokens JSONB, -- Array of token balances [{"token_id": "uuid", "balance": "string", "balance_usd": "decimal"}]
    
    -- Entry information
    entry_price_usd DECIMAL(20, 8), -- USD value at entry
    entry_block_number BIGINT,
    entry_transaction_hash VARCHAR(66),
    entry_time TIMESTAMPTZ NOT NULL,
    
    -- Current status
    is_active BOOLEAN DEFAULT TRUE,
    last_update_block BIGINT,
    last_update_time TIMESTAMPTZ DEFAULT NOW(),
    
    -- Rewards information
    pending_rewards JSONB, -- Array of pending rewards [{"token_id": "uuid", "amount": "string", "amount_usd": "decimal"}]
    claimed_rewards JSONB, -- Array of claimed rewards with timestamps
    total_rewards_usd DECIMAL(20, 8) DEFAULT 0,
    
    -- P&L calculation
    current_value_usd DECIMAL(20, 8),
    unrealized_pnl_usd DECIMAL(20, 8), -- current_value_usd - entry_price_usd
    realized_pnl_usd DECIMAL(20, 8) DEFAULT 0, -- From partial closes/claims
    total_fees_paid_usd DECIMAL(20, 8) DEFAULT 0,
    
    -- Metadata
    metadata JSONB, -- Additional protocol-specific data
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create indexes for efficient queries
CREATE INDEX idx_yield_positions_user_id ON yield_positions(user_id);
CREATE INDEX idx_yield_positions_wallet_id ON yield_positions(wallet_id);
CREATE INDEX idx_yield_positions_pool_id ON yield_positions(pool_id);
CREATE INDEX idx_yield_positions_protocol_id ON yield_positions(protocol_id);
CREATE INDEX idx_yield_positions_chain_id ON yield_positions(chain_id);
CREATE INDEX idx_yield_positions_active ON yield_positions(is_active);
CREATE INDEX idx_yield_positions_entry_time ON yield_positions(entry_time);
CREATE INDEX idx_yield_positions_user_active ON yield_positions(user_id, is_active);
CREATE INDEX idx_yield_positions_wallet_active ON yield_positions(wallet_id, is_active);

-- Composite index for efficient user position queries
CREATE INDEX idx_yield_positions_user_wallet_active ON yield_positions(user_id, wallet_id, is_active);

-- Create trigger for updated_at
CREATE TRIGGER update_yield_positions_updated_at BEFORE UPDATE
    ON yield_positions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create function to calculate unrealized P&L automatically
CREATE OR REPLACE FUNCTION calculate_position_pnl()
RETURNS TRIGGER AS $$
BEGIN
    -- Auto-calculate unrealized P&L when current_value_usd or entry_price_usd changes
    IF NEW.current_value_usd IS NOT NULL AND NEW.entry_price_usd IS NOT NULL THEN
        NEW.unrealized_pnl_usd = NEW.current_value_usd - NEW.entry_price_usd;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to auto-calculate P&L
CREATE TRIGGER calculate_position_pnl_trigger
    BEFORE INSERT OR UPDATE ON yield_positions
    FOR EACH ROW EXECUTE FUNCTION calculate_position_pnl();