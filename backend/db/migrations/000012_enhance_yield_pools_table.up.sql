-- Enhance yield_pools table with protocol relationship and additional fields
ALTER TABLE yield_pools 
ADD COLUMN IF NOT EXISTS protocol_id UUID REFERENCES protocols(id) ON DELETE SET NULL,
ADD COLUMN IF NOT EXISTS chain_id INTEGER,
ADD COLUMN IF NOT EXISTS pool_address VARCHAR(42),
ADD COLUMN IF NOT EXISTS token_addresses JSONB, -- Array of token contract addresses
ADD COLUMN IF NOT EXISTS fees_apr DECIMAL(10, 4), -- Fee-based APR
ADD COLUMN IF NOT EXISTS risk_level VARCHAR(20) DEFAULT 'medium',
ADD COLUMN IF NOT EXISTS min_deposit_usd DECIMAL(20, 2),
ADD COLUMN IF NOT EXISTS max_deposit_usd DECIMAL(20, 2),
ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT TRUE,
ADD COLUMN IF NOT EXISTS metadata JSONB; -- Additional pool-specific data

-- Create additional indexes
CREATE INDEX IF NOT EXISTS idx_yield_pools_protocol_id ON yield_pools(protocol_id);
CREATE INDEX IF NOT EXISTS idx_yield_pools_chain_id ON yield_pools(chain_id);
CREATE INDEX IF NOT EXISTS idx_yield_pools_active ON yield_pools(is_active);
CREATE INDEX IF NOT EXISTS idx_yield_pools_risk_level ON yield_pools(risk_level);

-- Update existing pools with protocol references (if protocols exist)
UPDATE yield_pools 
SET protocol_id = p.id,
    chain_id = CASE 
        WHEN chain = 'ethereum' THEN 1
        WHEN chain = 'polygon' THEN 137
        WHEN chain = 'arbitrum' THEN 42161
        WHEN chain = 'optimism' THEN 10
        WHEN chain = 'base' THEN 8453
        ELSE 1
    END
FROM protocols p 
WHERE LOWER(yield_pools.protocol) = LOWER(p.name) 
   OR LOWER(yield_pools.protocol) LIKE LOWER(p.name) || '%';

-- Add some sample token addresses for existing pools
UPDATE yield_pools 
SET token_addresses = '[
    "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
    "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"
]'::jsonb
WHERE pool_id LIKE '%weth-usdc%' OR symbol LIKE '%WETH%' OR symbol LIKE '%USDC%';

UPDATE yield_pools 
SET token_addresses = '[
    "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"
]'::jsonb
WHERE pool_id LIKE '%usdc%' AND pool_id NOT LIKE '%weth%';

-- Set reasonable defaults for existing pools
UPDATE yield_pools 
SET risk_level = CASE 
        WHEN protocol ILIKE '%uniswap%' OR protocol ILIKE '%aave%' THEN 'low'
        WHEN protocol ILIKE '%curve%' OR protocol ILIKE '%compound%' THEN 'low'
        WHEN protocol ILIKE '%convex%' THEN 'medium'
        ELSE 'medium'
    END,
    min_deposit_usd = 1.00,
    is_active = TRUE
WHERE risk_level IS NULL;