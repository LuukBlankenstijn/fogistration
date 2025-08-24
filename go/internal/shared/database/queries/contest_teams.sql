-- name: DeleteAllContestTeams :exec
DELETE FROM contest_teams WHERE contest_id = $1;

-- name: InsertContestTeams :copyfrom
INSERT INTO contest_teams (contest_id, team_id)
VALUES ($1, $2);
