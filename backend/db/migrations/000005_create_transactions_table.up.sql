-- Create transaction type enum
CREATE TYPE transaction_type AS ENUM ('send', 'receive', 'swap', 'approve', 'bridge', 'stake', 'unstake');
CREATE TYPE transaction_status AS ENUM ('pending', 'confirmed', 'failed');

-- Create transactions table
CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hash VARCHAR(66) UNIQUE NOT NULL,
    chain_id INTEGER NOT NULL,
    from_address VARCHAR(42) NOT NULL,
    to_address VARCHAR(42),
    value DECIMAL(78, 0),
    gas_used BIGINT,
    gas_price DECIMAL(30, 0),
    gas_fee_usd DECIMAL(30, 10),
    block_number BIGINT,
    timestamp TIMESTAMPTZ NOT NULL,
    status transaction_status NOT NULL DEFAULT 'pending',
    type transaction_type NOT NULL,
    metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create indexes
CREATE INDEX idx_transactions_hash ON transactions(hash);
CREATE INDEX idx_transactions_from_address ON transactions(from_address);
CREATE INDEX idx_transactions_to_address ON transactions(to_address);
CREATE INDEX idx_transactions_timestamp ON transactions(timestamp DESC);
CREATE INDEX idx_transactions_chain_id ON transactions(chain_id);
CREATE INDEX idx_transactions_type ON transactions(type);
CREATE INDEX idx_transactions_status ON transactions(status);

-- Create trigger for updated_at
CREATE TRIGGER update_transactions_updated_at BEFORE UPDATE
    ON transactions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create user_transactions junction table for faster lookups
CREATE TABLE IF NOT EXISTS user_transactions (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    transaction_id UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    wallet_id UUID NOT NULL REFERENCES wallets(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, transaction_id)
);

-- Create indexes for user_transactions
CREATE INDEX idx_user_transactions_user_id ON user_transactions(user_id);
CREATE INDEX idx_user_transactions_transaction_id ON user_transactions(transaction_id);