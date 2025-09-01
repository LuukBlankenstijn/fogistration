package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/LuukBlankenstijn/fogistration/internal/grpc/pubsub"
	"github.com/LuukBlankenstijn/fogistration/internal/grpc/service/utils"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/config"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/pb"
)

type reloadService struct {
	queries *database.Queries
	pubsub  *pubsub.Manager
	config  config.GrpcConfig
}

func (r *reloadService) PushUpdate(ctx context.Context, ip string) {
	go r.pushWallpaper(ctx, ip)
	go r.pushContestUrl(ctx, ip)
}

func (r *reloadService) pushContestUrl(ctx context.Context, ip string) {
	// TODO: impolement sending the autologinURL
}

func (r *reloadService) pushWallpaper(ctx context.Context, ip string) {
	var renderedWallpaper []byte
	publish := func() {
		r.pubsub.Publish(ip, &pb.ServerMessage{
			Message: &pb.ServerMessage_Wallpaper{
				Wallpaper: &pb.Wallpaper{
					Size: uint64(len(renderedWallpaper)),
					Data: renderedWallpaper,
				},
			},
		})
	}
	team, err := r.queries.GetTeamByIp(ctx, database.PgTextFromString(&ip))
	if err != nil {
		renderedWallpaper, err = utils.RenderNoTeamAssigned(ip)
		if err != nil {
			logging.Info(fmt.Sprintf("failed to render wallpaper for client %s: ", ip), err)
			return
		}
		publish()
		return
	}
	var teamName string
	if team.DisplayName.Valid {
		teamName = team.DisplayName.String
	} else {
		teamName = team.Name
	}

	contest_team, err := r.queries.GetContestForTeam(ctx, team.ID)
	if err != nil {
		renderedWallpaper, err = utils.RenderInactiveContest(ip, teamName)
		if err != nil {
			logging.Info(fmt.Sprintf("failed to render wallpaper for client %s: ", ip), err)
			return
		}
		publish()
		return
	}

	wallpaper, err := r.queries.GetWallpaperById(ctx, contest_team.ContestID)

	// No background in database, render a black background with a notfound text
	if err != nil {
		renderedWallpaper, err = utils.RenderNoWallpaperWatermark()
		if err != nil {
			logging.Info(fmt.Sprintf("failed to render wallpaper for client %s: ", ip), err)
			return
		}
	} else {
		var baseWallpaper []byte
		if wallpaper.Filename.Valid {
			baseWallpaper, err = r.loadWallpaperFile(wallpaper.Filename.String)
			if err != nil {
				logging.Info(fmt.Sprintf("failed to load wallpaper for client %s: ", ip), err)
			}
		}

		if baseWallpaper == nil {
			if wallpaper.Layout == nil {
				renderedWallpaper, err = utils.RenderNoWallpaperWatermark()
				if err != nil {
					logging.Info(fmt.Sprintf("failed to render wallpaper for client %s: ", ip), err)
					return
				}
			} else {
				renderedWallpaper, err = utils.RenderLayoutOnly(*wallpaper.Layout, teamName, ip)
				if err != nil {
					logging.Info(fmt.Sprintf("failed to render wallpaper for client %s: ", ip), err)
					return
				}
			}
		} else {
			if wallpaper.Layout == nil {
				renderedWallpaper = baseWallpaper
			} else {
				renderedWallpaper, err = utils.RenderCompleteBackground(baseWallpaper, *wallpaper.Layout, teamName, ip)
				if err != nil {
					logging.Info(fmt.Sprintf("failed to render wallpaper for client %s: ", ip), err)
					return
				}
			}
		}
	}

	publish()
}

func (r *reloadService) loadWallpaperFile(filename string) ([]byte, error) {
	base := filepath.Base(filename)
	fp := filepath.Join(r.config.WallpaperDir, base)

	f, err := os.Open(fp)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, errors.New("wallpaperfile not found")
		}
		return nil, errors.New("error while opening wallpaper file")
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, errors.New("error while reading wallpaper file")
	}
	return b, nil
}
