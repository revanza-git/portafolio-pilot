-- Create item type enum
CREATE TYPE watchlist_item_type AS ENUM ('token', 'pool', 'protocol');

-- Create watchlists table
CREATE TABLE IF NOT EXISTS watchlists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    item_type watchlist_item_type NOT NULL,
    item_ref_id INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    -- Unique constraint to prevent duplicate entries
    UNIQUE(user_id, item_type, item_ref_id)
);

-- Create indexes
CREATE INDEX idx_watchlists_user_id ON watchlists(user_id);
CREATE INDEX idx_watchlists_item_type ON watchlists(item_type);
CREATE INDEX idx_watchlists_item_ref_id ON watchlists(item_ref_id);
CREATE INDEX idx_watchlists_user_item ON watchlists(user_id, item_type);

-- Create trigger for updated_at
CREATE TRIGGER update_watchlists_updated_at BEFORE UPDATE
    ON watchlists FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();