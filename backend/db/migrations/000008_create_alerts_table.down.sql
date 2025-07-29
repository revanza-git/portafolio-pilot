-- Drop alerts tables and types
DROP TABLE IF EXISTS alert_history;
DROP TRIGGER IF EXISTS update_alerts_updated_at ON alerts;
DROP TABLE IF EXISTS alerts;
DROP TYPE IF EXISTS alert_status;