-- name: GetContestByExternalId :one
SELECT * FROM contests
WHERE external_id = $1;

-- name: ListContests :many
SELECT * FROM contests
ORDER BY id;

-- name: GetNextOrActiveContest :one
SELECT * FROM contests
WHERE end_time > NOW()
ORDER BY start_time ASC
LIMIT 1;

-- name: GetContestHashes :many
SELECT id, hash FROM contests
ORDER BY id;

-- name: UpsertContest :exec
INSERT INTO contests (
    id,
    external_id,
    formal_name,
    start_time,
    end_time,
    created_at,
    updated_at,
    hash
) VALUES (
    $1, $2, $3, $4, $5, NOW(), NOW(), $6
)
ON CONFLICT (id) 
DO UPDATE SET
    external_id = EXCLUDED.external_id,
    formal_name = EXCLUDED.formal_name,
    start_time = EXCLUDED.start_time,
    end_time = EXCLUDED.end_time,
    updated_at = NOW(),
    hash = EXCLUDED.hash;

-- name: GetContestByIp :one
SELECT c.*
FROM contests c
JOIN contest_teams ct ON c.id = ct.contest_id
JOIN teams t on ct.team_id = t.id
WHERE t.ip = $1;
