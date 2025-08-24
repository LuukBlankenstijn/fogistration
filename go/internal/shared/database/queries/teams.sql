-- name: GetTeamById :one
SELECT * FROM teams
WHERE id = $1;

-- name: GetTeamByExternalId :one
SELECT * FROM teams
WHERE external_id = $1;

-- name: GetTeamByIp :one
SELECT * FROM teams
WHERE ip = $1;

-- name: GetAllTeams :many
SELECT * FROM teams;

-- name: GetTeamHashes :many
SELECT id, hash FROM teams
ORDER BY id;

-- name: UpsertTeam :exec
INSERT INTO teams (
    id,
    external_id,
    name,
    display_name,
    ip,
    created_at,
    updated_at,
    hash
) VALUES (
    $1, $2, $3, $4, $5, NOW(), NOW(), $6
)
ON CONFLICT (id)
DO UPDATE SET
    external_id = EXCLUDED.external_id,
    name = EXCLUDED.name,
    display_name = EXCLUDED.display_name,
    ip = EXCLUDED.ip,
    updated_at = NOW(),
    hash = EXCLUDED.hash;

-- name: UpdateIp :one
UPDATE teams
SET ip = $2
WHERE external_id= $1
RETURNING *;

-- name: ClaimTeam :one
WITH target AS (
  SELECT t.id
  FROM teams t
  JOIN contest_teams ct ON t.id = ct.team_id
  WHERE t.ip IS NULL AND ct.contest_id = $2
  FOR UPDATE SKIP LOCKED
  LIMIT 1
)
UPDATE teams t
SET ip = $1
FROM target
WHERE t.id = target.id
RETURNING t.*;
