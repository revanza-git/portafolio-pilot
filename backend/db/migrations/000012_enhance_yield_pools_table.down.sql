-- Remove added columns from yield_pools table
ALTER TABLE yield_pools 
DROP COLUMN IF EXISTS protocol_id,
DROP COLUMN IF EXISTS chain_id,
DROP COLUMN IF EXISTS pool_address,
DROP COLUMN IF EXISTS token_addresses,
DROP COLUMN IF EXISTS fees_apr,
DROP COLUMN IF EXISTS risk_level,
DROP COLUMN IF EXISTS min_deposit_usd,
DROP COLUMN IF EXISTS max_deposit_usd,
DROP COLUMN IF EXISTS is_active,
DROP COLUMN IF EXISTS metadata;