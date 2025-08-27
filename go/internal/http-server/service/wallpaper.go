package service

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/LuukBlankenstijn/fogistration/internal/shared/config"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrWallpaperFileNotFound = errors.New("wallpaper file not found")
	ErrWallpaperNotFound     = errors.New("wallpaper not found")
	ErrWallpaperRead         = errors.New("failed to read wallpaper file")
	ErrInvalidContentType    = errors.New("content-type must be image/png or multipart/form-data")
	ErrNotPNG                = errors.New("file is not a PNG")
	ErrWriteFile             = errors.New("failed to write file")
	ErrDeleteFile            = errors.New("failed to delete file")
	ErrDB                    = errors.New("failed to update database")
)

type wallpaperService struct {
	c    *config.HttpConfig
	q    *database.Queries
	pool *pgxpool.Pool
}

func newWallpaperService(c *config.HttpConfig, pool *pgxpool.Pool) *wallpaperService {
	q := database.New(pool)
	return &wallpaperService{
		c,
		q,
		pool,
	}
}

func (s *wallpaperService) LoadWallpaperFile(filename string) ([]byte, error) {
	base := filepath.Base(filename)
	fp := filepath.Join(s.c.WallpaperDir, base)

	f, err := os.Open(fp)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrWallpaperFileNotFound
		}
		return nil, ErrWallpaperRead
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, ErrWallpaperRead
	}
	return b, nil
}

func (s *wallpaperService) SaveWallpaperFile(ctx context.Context, id int32, data []byte) (string, error) {
	hasNew := len(data) > 0
	var name string

	// 1) Validate & (if new) write file to disk first.
	if hasNew {
		if len(data) < 8 || !bytes.Equal(data[:8], []byte{137, 80, 78, 71, 13, 10, 26, 10}) {
			return "", ErrNotPNG
		}
		if err := os.MkdirAll(s.c.WallpaperDir, 0o755); err != nil {
			return "", fmt.Errorf("%w: %v", ErrWriteFile, err)
		}

		name = s.randHex(12) + ".png"
		path := filepath.Join(s.c.WallpaperDir, name)

		tmp := path + ".tmp"
		if err := os.WriteFile(tmp, data, 0o644); err != nil {
			return "", fmt.Errorf("%w: %v", ErrWriteFile, err)
		}
		if err := os.Rename(tmp, path); err != nil {
			_ = os.Remove(tmp)
			return "", fmt.Errorf("%w: %v", ErrWriteFile, err)
		}
	}

	// 2) Transaction: fetch previous row and upsert new filename (or clear).
	var oldFile string
	var oldValid bool

	err := withTx(ctx, s.pool, func(ctx context.Context, q *database.Queries) error {
		prev, err := q.GetWallpaperById(ctx, id)
		if err != nil {
			// DB error (or not found) — surface as DB failure
			return fmt.Errorf("%w: %v", ErrDB, err)
		}

		if prev.Filename.Valid && prev.Filename.String != "" {
			oldFile = filepath.Base(prev.Filename.String)
			oldValid = true
		}

		var pgText *string
		if hasNew {
			pgText = &name
		}

		if _, err := q.UpsertWallpaperFilename(ctx, database.UpsertWallpaperFilenameParams{
			ID:       prev.ID,
			Filename: database.PgTextFromString(pgText),
		}); err != nil {
			return fmt.Errorf("%w: %v", ErrDB, err)
		}
		return nil
	})
	if err != nil {
		// Rollback already done by withTx. If we created a new file, remove it.
		if hasNew {
			_ = os.Remove(filepath.Join(s.c.WallpaperDir, name))
		}
		return "", err
	}

	// 3) Post-commit cleanup: remove old file if cleared/replaced.
	if oldValid {
		if !hasNew || (hasNew && oldFile != name) {
			_ = os.Remove(filepath.Join(s.c.WallpaperDir, oldFile))
		}
	}

	return name, nil
}

func (s *wallpaperService) DeleteWallpaperFile(ctx context.Context, id int32) error {
	wallpaper, err := s.q.GetWallpaperById(ctx, id)
	if err != nil {
		return ErrWallpaperNotFound
	}
	if !wallpaper.Filename.Valid {
		return nil
	}
	_, err = s.q.UpsertWallpaperFilename(ctx, database.UpsertWallpaperFilenameParams{
		Filename: database.PgTextFromString(nil),
		ID:       wallpaper.ID,
	})
	if err != nil {
		return ErrDB
	}
	if err := os.Remove(filepath.Join(s.c.WallpaperDir, wallpaper.Filename.String)); err != nil {
		// try to fix, best effort
		_, _ = s.q.UpsertWallpaperFilename(ctx, database.UpsertWallpaperFilenameParams{
			Filename: wallpaper.Filename,
			ID:       wallpaper.ID,
		})
		return ErrDeleteFile
	}
	return nil
}

func (s *wallpaperService) randHex(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}
