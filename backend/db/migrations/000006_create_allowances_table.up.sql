-- Create token allowances table
CREATE TABLE IF NOT EXISTS token_allowances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_id UUID NOT NULL REFERENCES wallets(id) ON DELETE CASCADE,
    token_id UUID NOT NULL REFERENCES tokens(id) ON DELETE CASCADE,
    spender_address VARCHAR(42) NOT NULL,
    spender_name VARCHAR(255),
    allowance DECIMAL(78, 0) NOT NULL,
    allowance_usd DECIMAL(30, 10),
    transaction_hash VARCHAR(66),
    block_number BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(wallet_id, token_id, spender_address)
);

-- Create indexes
CREATE INDEX idx_allowances_wallet_id ON token_allowances(wallet_id);
CREATE INDEX idx_allowances_token_id ON token_allowances(token_id);
CREATE INDEX idx_allowances_spender ON token_allowances(spender_address);
CREATE INDEX idx_allowances_updated_at ON token_allowances(updated_at);

-- Create trigger for updated_at
CREATE TRIGGER update_allowances_updated_at BEFORE UPDATE
    ON token_allowances FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();