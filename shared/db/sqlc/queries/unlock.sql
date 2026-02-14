-- name: CreateUnlock :one
INSERT INTO unlock (user_id, expires_at)
VALUES ($1, $2)
RETURNING *;

-- name: GetUnlock :one
SELECT *
FROM unlock
WHERE unlock_id = $1
LIMIT 1;

-- name: ListUnlocksByUser :many
SELECT *
FROM unlock
WHERE user_id = $1
ORDER BY expires_at DESC;

-- name: DeleteUnlock :exec
DELETE FROM unlock
WHERE unlock_id = $1;
