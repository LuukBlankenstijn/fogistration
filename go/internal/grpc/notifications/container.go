package notifications

import (
	"github.com/LuukBlankenstijn/fogistration/internal/grpc/service"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
)

type NotificationHandlerContainer struct {
	Team       *teamNotificationHandler
	Client     *clientNotificationHandler
	Wallpapper *wallpaperNotificationHandler
}

func newHandlerContainer(s *service.ServiceContainer, queries *database.Queries) *NotificationHandlerContainer {
	return &NotificationHandlerContainer{
		Team:       newTeamHanlder(s),
		Client:     newClientHandler(s, queries),
		Wallpapper: newWallpaperHandler(s, queries),
	}
}
