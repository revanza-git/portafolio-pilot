-- Add is_admin field to users table
ALTER TABLE users 
ADD COLUMN is_admin BOOLEAN NOT NULL DEFAULT FALSE;

-- Create index on is_admin for faster admin queries
CREATE INDEX idx_users_is_admin ON users(is_admin);