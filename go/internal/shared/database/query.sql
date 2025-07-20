-- CONTESTS


-- name: ListContests :many
SELECT * FROM contests
ORDER BY id;

-- name: GetNextOrActiveContest :many
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


-- TEAMS:


-- name: GetTeamById :one
SELECT * FROM teams
WHERE id = $1;

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

-- name: UpdateIp :exec
UPDATE teams
SET ip = $2
WHERE id = $1;


-- CONTEST TEAMS


-- name: DeleteAllContestTeams :exec
DELETE FROM contest_teams WHERE contest_id = $1;

-- name: InsertContestTeams :copyfrom
INSERT INTO contest_teams (contest_id, team_id)
VALUES ($1, $2);


-- COMMANDS


-- name: EnqueueCommand :exec
INSERT INTO message_queue (command_type, payload)
VALUES ($1, $2);

-- name: DequeueCommand :one
DELETE FROM message_queue 
WHERE id = (
    SELECT id FROM message_queue 
    ORDER BY created_at ASC 
    FOR UPDATE SKIP LOCKED 
    LIMIT 1
)
RETURNING id, command_type, payload;
