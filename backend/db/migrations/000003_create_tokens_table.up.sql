-- Create tokens table
CREATE TABLE IF NOT EXISTS tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    address VARCHAR(42) NOT NULL,
    chain_id INTEGER NOT NULL,
    symbol VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    decimals INTEGER NOT NULL,
    logo_uri TEXT,
    price_usd DECIMAL(30, 10),
    price_change_24h DECIMAL(10, 4),
    market_cap DECIMAL(30, 2),
    total_supply DECIMAL(78, 0),
    last_updated TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(address, chain_id)
);

-- Create indexes
CREATE INDEX idx_tokens_chain_id ON tokens(chain_id);
CREATE INDEX idx_tokens_symbol ON tokens(symbol);
CREATE INDEX idx_tokens_address_chain ON tokens(address, chain_id);

-- Create trigger for updated_at
CREATE TRIGGER update_tokens_updated_at BEFORE UPDATE
    ON tokens FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();