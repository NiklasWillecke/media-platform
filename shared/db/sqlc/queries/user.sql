-- name: CreateUser :one
INSERT INTO "user" (password, email)
VALUES ($1, $2)
RETURNING *;

-- name: GetUserById :one
SELECT *
FROM "user"
WHERE user_id = $1
LIMIT 1;

-- name: GetUserByEmail :one
SELECT *
FROM "user"
WHERE email = $1
LIMIT 1;

-- name: CheckUserExists :one
SELECT EXISTS(SELECT 1 FROM "user" WHERE email = $1);

-- name: ListUsers :many
SELECT *
FROM "user"
ORDER BY created_at DESC;

-- name: UpdateUser :one
UPDATE "user"
SET password = $2,
    email = $3
WHERE user_id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM "user"
WHERE user_id = $1;
