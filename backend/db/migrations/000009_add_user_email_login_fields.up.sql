-- Add email and last_login_at fields to users table
ALTER TABLE users 
ADD COLUMN email VARCHAR(255) UNIQUE,
ADD COLUMN last_login_at TIMESTAMPTZ;

-- Create index on email
CREATE INDEX idx_users_email ON users(email);

-- Create nonce_storage table for better nonce management
CREATE TABLE IF NOT EXISTS nonce_storage (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    address VARCHAR(42) NOT NULL,
    nonce VARCHAR(64) NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create indexes for nonce storage
CREATE INDEX idx_nonce_storage_address ON nonce_storage(address);
CREATE INDEX idx_nonce_storage_nonce ON nonce_storage(nonce);
CREATE INDEX idx_nonce_storage_expires_at ON nonce_storage(expires_at);

-- Create cleanup function for expired nonces
CREATE OR REPLACE FUNCTION cleanup_expired_nonces()
RETURNS void AS $$
BEGIN
    DELETE FROM nonce_storage 
    WHERE expires_at < NOW() OR used = TRUE;
END;
$$ LANGUAGE plpgsql;