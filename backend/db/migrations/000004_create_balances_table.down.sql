-- Drop balance tables
DROP TRIGGER IF EXISTS update_balances_updated_at ON balances;
DROP TABLE IF EXISTS balance_history;
DROP TABLE IF EXISTS balances;