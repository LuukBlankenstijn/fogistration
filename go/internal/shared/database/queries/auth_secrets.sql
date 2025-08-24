-- name: GetAuthSecret :one
SELECT *
FROM auth_secrets
WHERE user_id = $1
LIMIT 1;

-- name: CreateAuthSecret :one
INSERT INTO auth_secrets (user_id, password_hash, salt)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateAuthSecret :one
UPDATE auth_secrets
SET
  password_hash = $2,
  salt          = $3,
  updated_at    = now()
WHERE user_id = $1
RETURNING *;

-- name: UpsertAuthSecret :one
INSERT INTO auth_secrets (user_id, password_hash, salt)
VALUES ($1, $2, $3)
ON CONFLICT (user_id) DO UPDATE
SET
  password_hash = EXCLUDED.password_hash,
  salt          = EXCLUDED.salt,
  updated_at    = now()
RETURNING *;

-- name: DeleteAuthSecret :exec
DELETE FROM auth_secrets
WHERE user_id = $1;
