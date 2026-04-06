-- name: GetBlockCursor :one
SELECT * FROM block_cursors WHERE chain_id = $1;

-- name: UpsertBlockCursor :exec
INSERT INTO block_cursors (chain_id, last_block, updated_at)
VALUES ($1, $2, NOW())
ON CONFLICT (chain_id) DO UPDATE SET
    last_block = EXCLUDED.last_block,
    updated_at = NOW();

-- name: ResetBlockCursor :exec
DELETE FROM block_cursors WHERE chain_id = $1;

-- name: ResetAllBlockCursors :exec
TRUNCATE TABLE block_cursors;
