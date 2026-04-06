-- name: InsertAuditLog :one
INSERT INTO audit_log (action, actor, chain_id, details, prev_hash, hash)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetLatestAuditLog :one
SELECT * FROM audit_log ORDER BY created_at DESC LIMIT 1;

-- name: ListAuditLogs :many
SELECT * FROM audit_log ORDER BY created_at DESC LIMIT $1 OFFSET $2;
