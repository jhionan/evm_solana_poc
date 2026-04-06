-- name: InsertChainEvent :one
INSERT INTO chain_events (chain_id, event_type, tx_hash, log_index, block_number, raw_data)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (tx_hash, log_index) DO NOTHING
RETURNING *;

-- name: GetEventsByBlock :many
SELECT * FROM chain_events WHERE chain_id = $1 AND block_number = $2;

-- name: GetEventsByTxHash :many
SELECT * FROM chain_events WHERE tx_hash = $1;

-- name: TruncateChainEvents :exec
TRUNCATE TABLE chain_events;
