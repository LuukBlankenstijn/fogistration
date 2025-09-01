package notifications

import (
	"context"

	"github.com/LuukBlankenstijn/fogistration/internal/grpc/service"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database/dblisten"
)

type teamNotificationHandler struct {
	s *service.ServiceContainer
}

func newTeamHanlder(s *service.ServiceContainer) *teamNotificationHandler {
	return &teamNotificationHandler{s}
}

func (t *teamNotificationHandler) Handle(ctx context.Context, notification dblisten.Notification[database.Team]) {
	t.ipChange(ctx, notification)
}

func (t *teamNotificationHandler) ipChange(ctx context.Context, notification dblisten.Notification[database.Team]) bool {
	if notification.Old != nil && notification.New != nil && notification.Old.Ip == notification.New.Ip {
		return false
	}

	reloaded := false
	if notification.Old != nil && notification.Old.Ip.Valid {
		ip := notification.Old.Ip.String
		t.s.Reload.PushUpdate(ctx, ip)
		reloaded = true
	}

	if notification.New != nil && notification.New.Ip.Valid {
		ip := notification.New.Ip.String
		t.s.Reload.PushUpdate(ctx, ip)
		reloaded = true
	}
	return reloaded
}
