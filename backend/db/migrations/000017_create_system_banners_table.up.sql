-- Create banner level enum
CREATE TYPE banner_level AS ENUM ('info', 'warning', 'error', 'success');

-- Create system_banners table
CREATE TABLE IF NOT EXISTS system_banners (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255),
    message TEXT NOT NULL,
    level banner_level NOT NULL DEFAULT 'info',
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create indexes
CREATE INDEX idx_system_banners_active ON system_banners(active);
CREATE INDEX idx_system_banners_level ON system_banners(level);
CREATE INDEX idx_system_banners_created_at ON system_banners(created_at DESC);

-- Create trigger for updated_at
CREATE TRIGGER update_system_banners_updated_at BEFORE UPDATE
    ON system_banners FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();