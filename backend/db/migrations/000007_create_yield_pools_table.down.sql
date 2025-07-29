-- Drop yield pools table
DROP TRIGGER IF EXISTS update_yield_pools_updated_at ON yield_pools;
DROP TABLE IF EXISTS yield_pools;