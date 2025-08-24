-- name: UpsertWallpaperConfig :one
INSERT INTO wallpaper_configs (contest_id, config)
VALUES ($1, $2)
ON CONFLICT (contest_id) DO UPDATE
SET config = EXCLUDED.config
RETURNING *;

-- name: UpsertWallpaperFile :one
INSERT INTO wallpaper_configs (contest_id, filename)
VALUES ($1, $2)
ON CONFLICT (contest_id) DO UPDATE
SET filename = EXCLUDED.filename
RETURNING *;

-- name: GetWallpaperConfigByContest :one
SELECT * FROM wallpaper_configs
WHERE contest_id = $1;

-- name: DeleteWallpaperConfigByContest :exec
DELETE FROM wallpaper_configs WHERE contest_id = $1;
