-- Drop position-related triggers and functions
DROP TRIGGER IF EXISTS calculate_position_pnl_trigger ON yield_positions;
DROP FUNCTION IF EXISTS calculate_position_pnl();

-- Drop yield_positions table
DROP TABLE IF EXISTS yield_positions;