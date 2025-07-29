-- Create alert status enum
CREATE TYPE alert_status AS ENUM ('active', 'triggered', 'expired', 'disabled');

-- Create alerts table
CREATE TABLE IF NOT EXISTS alerts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    status alert_status NOT NULL DEFAULT 'active',
    target JSONB NOT NULL, -- {"type": "token|address|pool", "identifier": "...", "chainId": 1}
    conditions JSONB NOT NULL, -- Alert-specific conditions
    notification JSONB NOT NULL, -- {"email": true, "webhook": "https://..."}
    last_triggered_at TIMESTAMPTZ,
    trigger_count INT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create indexes
CREATE INDEX idx_alerts_user_id ON alerts(user_id);
CREATE INDEX idx_alerts_status ON alerts(status);
CREATE INDEX idx_alerts_type ON alerts(type);
CREATE INDEX idx_alerts_last_triggered ON alerts(last_triggered_at);
CREATE INDEX idx_alerts_target ON alerts USING GIN(target);

-- Create trigger for updated_at
CREATE TRIGGER update_alerts_updated_at BEFORE UPDATE
    ON alerts FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create alert history table for tracking triggers
CREATE TABLE IF NOT EXISTS alert_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    alert_id UUID NOT NULL REFERENCES alerts(id) ON DELETE CASCADE,
    triggered_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    conditions_snapshot JSONB,
    triggered_value JSONB, -- The actual value that triggered the alert
    notification_sent BOOLEAN DEFAULT FALSE,
    notification_error TEXT
);

-- Create index for alert history
CREATE INDEX idx_alert_history_alert_id ON alert_history(alert_id);
CREATE INDEX idx_alert_history_triggered_at ON alert_history(triggered_at DESC);