-- name: UpsertReward :one
INSERT INTO rewards (position_id, accrued_amount, last_calculated_at)
VALUES ($1, $2, NOW())
ON CONFLICT (position_id) DO UPDATE SET
    accrued_amount = EXCLUDED.accrued_amount,
    last_calculated_at = NOW()
RETURNING *;

-- name: GetRewardByPosition :one
SELECT * FROM rewards WHERE position_id = $1;

-- name: TruncateRewards :exec
TRUNCATE TABLE rewards;
