package notifications

import (
	"context"

	"github.com/LuukBlankenstijn/fogistration/internal/grpc/service"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database/dblisten"
)

type wallpaperNotificationHandler struct {
	s       *service.ServiceContainer
	queries *database.Queries
}

func newWallpaperHandler(s *service.ServiceContainer, queries *database.Queries) *wallpaperNotificationHandler {
	return &wallpaperNotificationHandler{s, queries}
}

func (w *wallpaperNotificationHandler) Handle(ctx context.Context, notification dblisten.Notification[database.Wallpaper]) {
	teamIps := []string{}
	if notification.Old != nil && notification.New == nil {
		result, err := w.queries.GetIpsForContest(ctx, notification.Old.ID)
		if err != nil {
			return
		}
		for _, ip := range result {
			if ip.Valid {
				teamIps = append(teamIps, ip.String)
			}
		}
	} else if notification.New != nil {
		result, err := w.queries.GetIpsForContest(ctx, notification.New.ID)
		if err != nil {
			return
		}
		for _, ip := range result {
			if ip.Valid {
				teamIps = append(teamIps, ip.String)
			}
		}
	}
	for _, ip := range teamIps {
		go w.s.Reload.PushUpdate(ctx, ip)
	}
}
