-- name: DeleteAllContestTeams :exec
DELETE FROM contest_teams WHERE contest_id = $1;

-- name: InsertContestTeams :copyfrom
INSERT INTO contest_teams (contest_id, team_id)
VALUES ($1, $2);

-- name: GetContestForTeam :one
SELECT * FROM contest_teams
WHERE team_id = $1
LIMIT 1;

-- name: GetIpsForContest :many
SELECT t.ip
FROM contest_teams ct
JOIN teams t ON ct.team_id = t.id
WHERE t.ip IS NOT NULL AND ct.contest_id = $1;

