-- name: CreateLocalUser :one
INSERT INTO users (username, email, role)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CreateExternalUser :one
INSERT INTO users (username, email, role, external_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1
LIMIT 1;

-- name: GetExternaluserById :one
SELECT * FROM users
WHERE external_id = $1
LIMIT 1;

-- name: GetUserByUsernameCI :one
SELECT * FROM users
WHERE lower(username) = lower($1)
LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;

-- name: UpdateUserProfile :one
UPDATE users
SET
  username = COALESCE(sqlc.narg('username'), username),
  email    = COALESCE(sqlc.narg('email'), email),
  role     = COALESCE(sqlc.narg('role'), role)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: TouchLastLogin :one
UPDATE users
SET last_login_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
