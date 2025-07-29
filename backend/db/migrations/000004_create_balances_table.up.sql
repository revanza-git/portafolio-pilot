-- Create balances table
CREATE TABLE IF NOT EXISTS balances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_id UUID NOT NULL REFERENCES wallets(id) ON DELETE CASCADE,
    token_id UUID NOT NULL REFERENCES tokens(id) ON DELETE CASCADE,
    balance DECIMAL(78, 0) NOT NULL,
    balance_usd DECIMAL(30, 10),
    block_number BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(wallet_id, token_id)
);

-- Create indexes
CREATE INDEX idx_balances_wallet_id ON balances(wallet_id);
CREATE INDEX idx_balances_token_id ON balances(token_id);
CREATE INDEX idx_balances_updated_at ON balances(updated_at);

-- Create trigger for updated_at
CREATE TRIGGER update_balances_updated_at BEFORE UPDATE
    ON balances FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create balance history table for tracking portfolio over time
CREATE TABLE IF NOT EXISTS balance_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_id UUID NOT NULL REFERENCES wallets(id) ON DELETE CASCADE,
    token_id UUID NOT NULL REFERENCES tokens(id) ON DELETE CASCADE,
    balance DECIMAL(78, 0) NOT NULL,
    balance_usd DECIMAL(30, 10),
    block_number BIGINT,
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create indexes for history
CREATE INDEX idx_balance_history_wallet_id ON balance_history(wallet_id);
CREATE INDEX idx_balance_history_recorded_at ON balance_history(recorded_at);
CREATE INDEX idx_balance_history_wallet_recorded ON balance_history(wallet_id, recorded_at DESC);