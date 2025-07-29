-- Create feature_flags table
CREATE TABLE IF NOT EXISTS feature_flags (
    name TEXT PRIMARY KEY,
    value JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create index on created_at for sorting
CREATE INDEX idx_feature_flags_created_at ON feature_flags(created_at);

-- Create trigger for updated_at
CREATE TRIGGER update_feature_flags_updated_at BEFORE UPDATE
    ON feature_flags FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();