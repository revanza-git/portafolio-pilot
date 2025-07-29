-- Create PnL lot type enum
CREATE TYPE pnl_lot_type AS ENUM ('buy', 'sell');

-- Create pnl_lots table for FIFO/LIFO calculations
CREATE TABLE IF NOT EXISTS pnl_lots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_id UUID NOT NULL REFERENCES wallets(id) ON DELETE CASCADE,
    token_id UUID NOT NULL REFERENCES tokens(id) ON DELETE CASCADE,
    transaction_hash VARCHAR(66) NOT NULL REFERENCES transactions(hash) ON DELETE CASCADE,
    chain_id INTEGER NOT NULL,
    type pnl_lot_type NOT NULL,
    quantity DECIMAL(78, 18) NOT NULL,
    price_usd DECIMAL(30, 10) NOT NULL,
    remaining_quantity DECIMAL(78, 18) NOT NULL,
    block_number BIGINT NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create indexes for performance
CREATE INDEX idx_pnl_lots_wallet_id ON pnl_lots(wallet_id);
CREATE INDEX idx_pnl_lots_token_id ON pnl_lots(token_id);
CREATE INDEX idx_pnl_lots_wallet_token ON pnl_lots(wallet_id, token_id);
CREATE INDEX idx_pnl_lots_timestamp ON pnl_lots(timestamp ASC);
CREATE INDEX idx_pnl_lots_block_number ON pnl_lots(block_number ASC);
CREATE INDEX idx_pnl_lots_type ON pnl_lots(type);
CREATE INDEX idx_pnl_lots_remaining_quantity ON pnl_lots(remaining_quantity) WHERE remaining_quantity > 0;

-- Create composite index for FIFO/LIFO queries
CREATE INDEX idx_pnl_lots_fifo_lifo ON pnl_lots(wallet_id, token_id, timestamp ASC) WHERE remaining_quantity > 0;

-- Create trigger for updated_at
CREATE TRIGGER update_pnl_lots_updated_at BEFORE UPDATE
    ON pnl_lots FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();