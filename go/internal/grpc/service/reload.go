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
	wallpaper := r.pushWallpaper(ctx, ip)
	contestUrl := r.pushContestUrl(ctx, ip)

	r.pubsub.Publish(ip, &pb.ServerMessage{
		Message: &pb.ServerMessage_Update{
			Update: &pb.Update{
				Wallpaper:  wallpaper,
				ContestUrl: contestUrl,
			},
		},
	})
}

func (r *reloadService) pushContestUrl(ctx context.Context, ip string) *pb.ContestUrl {
	var url string
	if contest, err := r.queries.GetContestByIp(ctx, database.PgTextFromString(&ip)); err == nil {
		url = fmt.Sprintf("http://%s/api/contests/%s", r.config.DJHost, contest.ExternalID)
	} else {
		url = ""
	}

	return &pb.ContestUrl{
		Url: url,
	}
}

func (r *reloadService) pushWallpaper(ctx context.Context, ip string) *pb.Wallpaper {
	var renderedWallpaper []byte
	returnWallpaper := func() *pb.Wallpaper {
		if len(renderedWallpaper) == 0 {
			return nil
		}
		return &pb.Wallpaper{
			Size: uint64(len(renderedWallpaper)),
			Data: renderedWallpaper,
		}
	}
	team, err := r.queries.GetTeamByIp(ctx, database.PgTextFromString(&ip))
	if err != nil {
		renderedWallpaper, err = utils.RenderNoTeamAssigned(ip)
		if err != nil {
			logging.Info(fmt.Sprintf("failed to render wallpaper for client %s: ", ip), err)
			return nil
		}
		return returnWallpaper()
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
			return nil
		}
		return returnWallpaper()
	}

	wallpaper, err := r.queries.GetWallpaperById(ctx, contest_team.ContestID)

	// No background in database, render a black background with a notfound text
	if err != nil {
		renderedWallpaper, err = utils.RenderNoWallpaperWatermark()
		if err != nil {
			logging.Info(fmt.Sprintf("failed to render wallpaper for client %s: ", ip), err)
			return nil
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
					return nil
				}
			} else {
				renderedWallpaper, err = utils.RenderLayoutOnly(*wallpaper.Layout, teamName, ip)
				if err != nil {
					logging.Info(fmt.Sprintf("failed to render wallpaper for client %s: ", ip), err)
					return nil
				}
			}
		} else {
			if wallpaper.Layout == nil {
				renderedWallpaper = baseWallpaper
			} else {
				renderedWallpaper, err = utils.RenderCompleteBackground(baseWallpaper, *wallpaper.Layout, teamName, ip)
				if err != nil {
					logging.Info(fmt.Sprintf("failed to render wallpaper for client %s: ", ip), err)
					return nil
				}
			}
		}
	}

	return returnWallpaper()
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
