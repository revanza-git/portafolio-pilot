-- Drop transaction tables and types
DROP TABLE IF EXISTS user_transactions;
DROP TRIGGER IF EXISTS update_transactions_updated_at ON transactions;
DROP TABLE IF EXISTS transactions;
DROP TYPE IF EXISTS transaction_status;
DROP TYPE IF EXISTS transaction_type;