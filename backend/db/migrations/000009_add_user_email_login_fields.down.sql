-- Drop cleanup function
DROP FUNCTION IF EXISTS cleanup_expired_nonces();

-- Drop nonce_storage table
DROP TABLE IF EXISTS nonce_storage;

-- Remove email and last_login_at columns from users table
ALTER TABLE users 
DROP COLUMN IF EXISTS email,
DROP COLUMN IF EXISTS last_login_at;