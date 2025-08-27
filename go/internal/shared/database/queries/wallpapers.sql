-- name: UpsertWallpaperLayout :one
INSERT INTO wallpapers (id, layout)
VALUES ($1, $2)
ON CONFLICT (id) DO UPDATE
SET layout = EXCLUDED.layout
RETURNING *;

-- name: UpsertWallpaperFilename :one
INSERT INTO wallpapers (id, filename)
VALUES ($1, $2)
ON CONFLICT (id) DO UPDATE
SET filename = EXCLUDED.filename
RETURNING *;

-- name: GetWallpaperById :one
SELECT * FROM wallpapers
WHERE id = $1;

-- name: DeleteWallpaperById :exec
DELETE FROM wallpapers WHERE id = $1;
