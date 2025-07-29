-- Drop allowances table
DROP TRIGGER IF EXISTS update_allowances_updated_at ON token_allowances;
DROP TABLE IF EXISTS token_allowances;