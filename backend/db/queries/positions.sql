-- name: UpsertPosition :one
INSERT INTO positions (chain_id, wallet, amount, tier, staked_at, lock_until, status, tx_hash)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (id) DO UPDATE SET
    status = EXCLUDED.status,
    updated_at = NOW()
RETURNING *;

-- name: GetPosition :one
SELECT * FROM positions WHERE id = $1;

-- name: ListPositionsByWallet :many
SELECT * FROM positions WHERE chain_id = $1 AND wallet = $2 ORDER BY staked_at DESC;

-- name: ListPositionsByChain :many
SELECT * FROM positions WHERE chain_id = $1 ORDER BY staked_at DESC;

-- name: UpdatePositionStatus :exec
UPDATE positions SET status = $2, updated_at = NOW() WHERE id = $1;

-- name: InsertPosition :one
INSERT INTO positions (chain_id, wallet, amount, tier, staked_at, lock_until, status, tx_hash)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: TruncatePositions :exec
TRUNCATE TABLE positions CASCADE;
