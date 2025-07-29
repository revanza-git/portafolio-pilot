-- Create yield pools table for storing APR/APY data
CREATE TABLE IF NOT EXISTS yield_pools (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pool_id VARCHAR(255) UNIQUE NOT NULL,
    protocol VARCHAR(100) NOT NULL,
    pool_name VARCHAR(255) NOT NULL,
    chain VARCHAR(50) NOT NULL,
    symbol VARCHAR(100) NOT NULL,
    tvl_usd DECIMAL(30, 2),
    apy DECIMAL(10, 4),
    apy_base DECIMAL(10, 4),
    apy_reward DECIMAL(10, 4),
    il_7d DECIMAL(10, 4), -- Impermanent loss 7 days
    stable_coin BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create indexes
CREATE INDEX idx_yield_pools_protocol ON yield_pools(protocol);
CREATE INDEX idx_yield_pools_chain ON yield_pools(chain);
CREATE INDEX idx_yield_pools_tvl ON yield_pools(tvl_usd DESC);
CREATE INDEX idx_yield_pools_apy ON yield_pools(apy DESC);
CREATE INDEX idx_yield_pools_updated ON yield_pools(updated_at);

-- Create trigger for updated_at
CREATE TRIGGER update_yield_pools_updated_at BEFORE UPDATE
    ON yield_pools FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();