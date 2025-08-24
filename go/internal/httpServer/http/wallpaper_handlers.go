package http

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/LuukBlankenstijn/fogistration/internal/httpServer/models"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	dbModels "github.com/LuukBlankenstijn/fogistration/internal/shared/database/models"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

// @Summary Get contest wallpaper config
// @Description Get the config for a wallpaper associated to the contest
// @Tags wallpaper
// @Accept json
// @Produce json
// @Param contestId path string true "contest Id"
// @Success 200 {object} models.WallpaperConfig
// @Failure 400 {object} ErrorResponse
// @Router /wallpaper/{contestId}/config [get]
// @Id GetWallpaperConfig
func (s *Server) GetWallpaperConfig(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	contest, ok := s.getContestFromRequest(w, r)
	if !ok {
		return
	}

	config := dbModels.DefaultWallpaperConfig()
	fullConfig, err := s.queries.GetWallpaperConfigByContest(ctx, contest.ExternalID)
	if err != nil && err != pgx.ErrNoRows {
		logging.Error("failed to get wallpaper config from database", err)
	} else if err == nil && fullConfig.Config != nil {
		config = *fullConfig.Config
	}

	err = json.NewEncoder(w).Encode(config)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// @Summary Set contest wallpaper config
// @Description Update the wallpaper config associated to the contest
// @Tags wallpaper
// @Accept  json
// @Produce json
// @Param   contestId     path string true "contest Id"
// @Param   config body models.WallpaperConfig true "Wallpaper config"
// @Success 200 {object} models.WallpaperConfig
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router  /wallpaper/{contestId}/config [put]
// @Id      SetWallpaperConfig
func (s *Server) SetWallpaperConfig(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	contest, ok := s.getContestFromRequest(w, r)
	if !ok {
		return
	}

	var cfg models.WallpaperConfig
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		logging.Error("failed to decode json", err)
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	updated, err := s.queries.UpsertWallpaperConfig(ctx,
		database.UpsertWallpaperConfigParams{
			ContestID: contest.ExternalID,
			Config:    &cfg,
		})
	if err != nil {
		logging.Error("failed to persist wallpaper config", err)
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updated.Config); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// @Summary Get contest wallpaper file
// @Description Get the wallpaper PNG associated with the contest
// @Tags wallpaper
// @Produce png
// @Param contestId path string true "contest Id"
// @Success 200 {string} binary "PNG bytes"
// @Header 200 {string} Content-Disposition "attachment; filename=wallpaper.png"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /wallpaper/{contestId} [get]
// @Id GetWallpaperFile
func (s *Server) GetWallpaper(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	contest, ok := s.getContestFromRequest(w, r)
	if !ok {
		return
	}

	config, err := s.queries.GetWallpaperConfigByContest(ctx, contest.ExternalID)
	if err != nil || !config.Filename.Valid {
		http.Error(w, "no wallpaper set", http.StatusNotFound)
		return
	}

	fp := filepath.Join(s.Config.WallpaperDir, filepath.Base(config.Filename.String))
	f, err := os.Open(fp)
	if err != nil {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}
	defer f.Close()

	w.Header().Set("Content-Type", "image/png")
	http.ServeContent(w, r, filepath.Base(fp), time.Now(), f)
}

// @Summary Set contest wallpaper file
// @Description Upload a new wallpaper PNG for the contest. If no file is provided (multipart without "file" or empty body), the wallpaper is removed.
// @Tags wallpaper
// @Accept  png
// @Accept  multipart/form-data
// @Produce json
// @Param   contestId path string true "contest Id"
// @Param   file formData file false "Wallpaper PNG file (omit to clear)"
// @Success 200 {object} map[string]string "filename and url (filename empty when cleared)"
// @Failure 400 {object} ErrorResponse
// @Failure 415 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router  /wallpaper/{contestId} [put]
// @Id      SetWallpaperFile
func (s *Server) SetWallpaper(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	contest, ok := s.getContestFromRequest(w, r)
	if !ok {
		return
	}

	var (
		data   []byte
		hasNew bool // true when we're setting a new file; false means "clear"
	)

	ct := r.Header.Get("Content-Type")
	switch {
	case strings.HasPrefix(ct, "multipart/form-data"):
		if err := r.ParseMultipartForm(25 << 20); err != nil { // 25MB
			http.Error(w, "bad multipart form", http.StatusBadRequest)
			return
		}
		file, _, err := r.FormFile("file")
		if err == http.ErrMissingFile {
			// No file field -> clear
			hasNew = false
		} else if err != nil {
			http.Error(w, "file read error", http.StatusBadRequest)
			return
		} else {
			defer file.Close()
			b, err := io.ReadAll(io.LimitReader(file, 50<<20))
			if err != nil {
				http.Error(w, "failed to read file", http.StatusBadRequest)
				return
			}
			if len(b) == 0 {
				// Empty file -> clear
				hasNew = false
			} else {
				data, hasNew = b, true
			}
		}

	case strings.HasPrefix(ct, "image/png"):
		// Raw PNG body; empty body => clear
		b, err := io.ReadAll(io.LimitReader(r.Body, 50<<20))
		if err != nil {
			http.Error(w, "failed to read body", http.StatusBadRequest)
			return
		}
		if len(b) == 0 {
			hasNew = false
		} else {
			data, hasNew = b, true
		}

	default:
		// If body is empty and no supported content-type, treat as "clear".
		// Otherwise, 415.
		if r.ContentLength == 0 {
			hasNew = false
		} else {
			http.Error(w, "Content-Type must be image/png or multipart/form-data", http.StatusUnsupportedMediaType)
			return
		}
	}

	// Always ensure storage dir exists (safe if clearing too)
	if err := os.MkdirAll(s.Config.WallpaperDir, 0o755); err != nil {
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}

	var name string

	if hasNew {
		// Quick PNG signature check
		if len(data) < 8 || !bytes.Equal(data[:8], []byte{137, 80, 78, 71, 13, 10, 26, 10}) {
			http.Error(w, "not a PNG", http.StatusBadRequest)
			return
		}

		name = randHex(12) + ".png"
		path := filepath.Join(s.Config.WallpaperDir, name)

		tmp := path + ".tmp"
		if err := os.WriteFile(tmp, data, 0o644); err != nil {
			http.Error(w, "write error", http.StatusInternalServerError)
			return
		}
		if err := os.Rename(tmp, path); err != nil {
			_ = os.Remove(tmp)
			http.Error(w, "save error", http.StatusInternalServerError)
			return
		}
	}

	// Fetch previous, to possibly delete file on disk
	prev, _ := s.queries.GetWallpaperConfigByContest(ctx, contest.ExternalID)

	// Update DB: set filename or clear (NULL)
	var pgText *string
	if hasNew {
		pgText = &name
	} else {
		pgText = nil
	}
	if _, err := s.queries.UpsertWallpaperFile(ctx, database.UpsertWallpaperFileParams{
		ContestID: contest.ExternalID,
		Filename:  database.PgTextFromString(pgText),
	}); err != nil {
		// If we just saved a new file but DB failed, try to remove it
		if hasNew {
			_ = os.Remove(filepath.Join(s.Config.WallpaperDir, name))
		}
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	// If clearing: remove old file (ignore errors)
	if !hasNew && prev.Filename.Valid && prev.Filename.String != "" {
		_ = os.Remove(filepath.Join(s.Config.WallpaperDir, filepath.Base(prev.Filename.String)))
	}

	// If replacing: remove old file different from new one
	if hasNew && prev.Filename.Valid && prev.Filename.String != "" && filepath.Base(prev.Filename.String) != name {
		_ = os.Remove(filepath.Join(s.Config.WallpaperDir, filepath.Base(prev.Filename.String)))
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"filename": name, // empty when cleared
		"url":      fmt.Sprintf("/wallpaper/%s", contest.ExternalID),
	})
}

func randHex(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

func (s *Server) getContestFromRequest(w http.ResponseWriter, r *http.Request) (*database.Contest, bool) {
	ctx := r.Context()
	contestId := chi.URLParam(r, "contestId")
	if contestId == "" {
		http.Error(w, "Malformed id", http.StatusBadRequest)
		return nil, false
	}

	contest, err := s.queries.GetContestByExternalId(ctx, contestId)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "failed to find contest", http.StatusBadRequest)
			return nil, false
		}
		logging.Error(fmt.Sprintf("failed to get contest with id %s", contestId), err)
		w.WriteHeader(http.StatusInternalServerError)
		return nil, false
	}

	return &contest, true
}
