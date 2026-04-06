-- +goose Up
CREATE TABLE block_cursors (chain_id TEXT PRIMARY KEY, last_block BIGINT NOT NULL, updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW());
CREATE TABLE positions (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), chain_id TEXT NOT NULL, wallet TEXT NOT NULL, amount NUMERIC NOT NULL, tier TEXT NOT NULL, staked_at TIMESTAMPTZ NOT NULL, lock_until TIMESTAMPTZ NOT NULL, status TEXT NOT NULL DEFAULT 'active', tx_hash TEXT NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW());
CREATE INDEX idx_positions_chain_wallet ON positions(chain_id, wallet);
CREATE INDEX idx_positions_status ON positions(status);
CREATE TABLE chain_events (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), chain_id TEXT NOT NULL, event_type TEXT NOT NULL, tx_hash TEXT NOT NULL, log_index INTEGER NOT NULL, block_number BIGINT NOT NULL, raw_data JSONB NOT NULL, indexed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), UNIQUE (tx_hash, log_index));
CREATE INDEX idx_chain_events_chain_block ON chain_events(chain_id, block_number);
CREATE TABLE rewards (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), position_id UUID NOT NULL REFERENCES positions(id) ON DELETE CASCADE, accrued_amount NUMERIC NOT NULL DEFAULT 0, last_calculated_at TIMESTAMPTZ NOT NULL DEFAULT NOW());
CREATE UNIQUE INDEX idx_rewards_position ON rewards(position_id);
CREATE TABLE audit_log (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), action TEXT NOT NULL, actor TEXT NOT NULL, chain_id TEXT, details JSONB, prev_hash TEXT, hash TEXT NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW());

-- +goose Down
DROP TABLE IF EXISTS audit_log; DROP TABLE IF EXISTS rewards; DROP TABLE IF EXISTS chain_events; DROP TABLE IF EXISTS positions; DROP TABLE IF EXISTS block_cursors;
