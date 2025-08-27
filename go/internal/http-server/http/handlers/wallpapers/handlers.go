package wallpapers

import (
	"context"
	"errors"
	"strings"

	"github.com/LuukBlankenstijn/fogistration/internal/http-server/service"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"github.com/danielgtaylor/huma/v2"
)

func (h *Handlers) getWallpaperLayout(ctx context.Context, request *getWallpaperLayoutRequest) (*wallpaperLayoutResponse, error) {
	wallpaper, err := h.Q.GetWallpaperById(ctx, request.ID)
	if err != nil || wallpaper.Layout == nil {
		return nil, huma.Error404NotFound("wallpaper not found")
	}

	return &wallpaperLayoutResponse{*wallpaper.Layout}, nil
}

func (h *Handlers) putWallpaperLayout(ctx context.Context, request *putWallpaperLayoutRequest) (*wallpaperLayoutResponse, error) {
	wallpaper, err := h.Q.UpsertWallpaperLayout(ctx, database.UpsertWallpaperLayoutParams{
		ID:     request.ID,
		Layout: &request.Body,
	})
	if err != nil {
		return nil, huma.Error500InternalServerError("")
	}

	return &wallpaperLayoutResponse{*wallpaper.Layout}, nil
}

func (h *Handlers) deleteWallpaperLayout(ctx context.Context, request *deleteWallpaperLayoutRequest) (*struct{}, error) {
	wallpaper, err := h.Q.GetWallpaperById(ctx, request.ID)
	if err != nil || wallpaper.Layout == nil {
		return nil, huma.Error404NotFound("wallpaper not found")
	}
	if _, err := h.Q.UpsertWallpaperLayout(ctx, database.UpsertWallpaperLayoutParams{
		Layout: nil,
		ID:     request.ID,
	}); err != nil {
		return nil, huma.Error500InternalServerError("failed to update database")
	}

	return &struct{}{}, nil
}

func (h *Handlers) getWallpaperFile(ctx context.Context, request *getWallpaperFileRequest) (*wallpaperFileResponse, error) {
	wallpaper, err := h.Q.GetWallpaperById(ctx, request.ID)
	if err != nil {
		return nil, huma.Error404NotFound("wallpaper not found")
	}

	if !wallpaper.Filename.Valid {
		return nil, huma.Error404NotFound("wallpaper file not found")
	}

	imageBytes, err := h.S.Wallpaper.LoadWallpaperFile(wallpaper.Filename.String)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrWallpaperFileNotFound):
			return nil, huma.Error404NotFound(err.Error())
		case errors.Is(err, service.ErrWallpaperRead):
			return nil, huma.Error500InternalServerError(err.Error())
		default:
			return nil, huma.Error500InternalServerError("unexpected error")
		}
	}

	return &wallpaperFileResponse{
		ContentType: "image/png",
		Body:        imageBytes,
	}, nil
}

func (h *Handlers) putWallpaperFile(ctx context.Context, request *putWallpaperFileRequest) (*wallpaperFileResponse, error) {
	if !strings.HasPrefix(request.ContentType, "image/png") {
		return nil, huma.Error415UnsupportedMediaType("contentType must be image/png")
	}
	_, err := h.S.Wallpaper.SaveWallpaperFile(ctx, request.ID, request.RawBody)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidContentType):
			return nil, huma.Error415UnsupportedMediaType(err.Error())
		case errors.Is(err, service.ErrNotPNG):
			return nil, huma.Error400BadRequest(err.Error())
		case errors.Is(err, service.ErrWriteFile):
			return nil, huma.Error500InternalServerError(err.Error())
		case errors.Is(err, service.ErrDB):
			return nil, huma.Error500InternalServerError("database error")
		default:
			return nil, huma.Error500InternalServerError("unexpected error")
		}
	}

	logging.Info("%+v", request.RawBody)

	return &wallpaperFileResponse{
		ContentType: "image/png",
		Body:        request.RawBody,
	}, nil
}

func (h *Handlers) deleteWallpaperFile(ctx context.Context, request *deleteWallpaperFileRequest) (*struct{}, error) {
	err := h.S.Wallpaper.DeleteWallpaperFile(ctx, request.ID)
	if err != nil {

		switch {
		case errors.Is(err, service.ErrDeleteFile):
		case errors.Is(err, service.ErrWriteFile):
			return nil, huma.Error500InternalServerError(err.Error())
		case errors.Is(err, service.ErrWallpaperNotFound):
			return nil, huma.Error404NotFound(err.Error())
		default:
			return nil, huma.Error500InternalServerError("unexpected error")
		}
	}
	return &struct{}{}, nil
}
