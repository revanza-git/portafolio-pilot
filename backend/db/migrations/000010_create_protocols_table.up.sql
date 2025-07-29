-- Create protocols table for DeFi protocol information
CREATE TABLE IF NOT EXISTS protocols (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) UNIQUE NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    website_url VARCHAR(255),
    logo_uri VARCHAR(255),
    category VARCHAR(50), -- 'dex', 'lending', 'staking', 'yield_farming', etc.
    total_tvl_usd DECIMAL(20, 2),
    chains JSONB, -- Array of supported chain IDs
    is_active BOOLEAN DEFAULT TRUE,
    risk_level VARCHAR(20) DEFAULT 'medium', -- 'low', 'medium', 'high'
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create indexes
CREATE INDEX idx_protocols_name ON protocols(name);
CREATE INDEX idx_protocols_slug ON protocols(slug);
CREATE INDEX idx_protocols_category ON protocols(category);
CREATE INDEX idx_protocols_tvl ON protocols(total_tvl_usd DESC);
CREATE INDEX idx_protocols_active ON protocols(is_active);

-- Create trigger for updated_at
CREATE TRIGGER update_protocols_updated_at BEFORE UPDATE
    ON protocols FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Insert some default protocols
INSERT INTO protocols (name, slug, description, category, website_url, logo_uri, chains, risk_level) VALUES
('Uniswap V3', 'uniswap-v3', 'A protocol for swapping and earning fees on ERC20 tokens', 'dex', 'https://uniswap.org', 'https://app.uniswap.org/favicon.ico', '[1, 10, 137, 8453, 42161]', 'low'),
('Aave V3', 'aave-v3', 'A decentralized non-custodial liquidity market protocol', 'lending', 'https://aave.com', 'https://aave.com/favicon.ico', '[1, 10, 137, 8453, 42161]', 'low'),
('Compound V3', 'compound-v3', 'An algorithmic, autonomous interest rate protocol', 'lending', 'https://compound.finance', 'https://compound.finance/favicon.ico', '[1, 10, 137, 8453]', 'low'),
('Curve Finance', 'curve', 'A decentralized exchange optimized for stablecoins', 'dex', 'https://curve.fi', 'https://curve.fi/favicon.ico', '[1, 10, 137, 8453, 42161]', 'medium'),
('Convex Finance', 'convex', 'A platform that boosts Curve Finance yields through staking CRV', 'yield_farming', 'https://convexfinance.com', 'https://convexfinance.com/favicon.ico', '[1]', 'medium'),
('Lido', 'lido', 'Liquid staking for Ethereum 2.0', 'staking', 'https://lido.fi', 'https://lido.fi/favicon.ico', '[1, 10, 137, 8453, 42161]', 'low');