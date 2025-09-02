package notifications

import (
	"context"

	"github.com/LuukBlankenstijn/fogistration/internal/grpc/service"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database/dblisten"
)

type clientNotificationHandler struct {
	s       *service.ServiceContainer
	queries *database.Queries
}

func newClientHandler(s *service.ServiceContainer, queries *database.Queries) *clientNotificationHandler {
	return &clientNotificationHandler{s, queries}
}

func (c *clientNotificationHandler) Handle(ctx context.Context, notification dblisten.Notification[database.Client]) {
	c.sync(ctx, notification)
}

func (c *clientNotificationHandler) sync(ctx context.Context, notification dblisten.Notification[database.Client]) bool {
	if notification.New != nil && notification.New.PendingSync {
		go c.s.Reload.PushUpdate(ctx, notification.New.Ip)
		// best effort
		_ = c.queries.SetPendingSync(ctx, database.SetPendingSyncParams{
			ID:          notification.New.ID,
			PendingSync: false,
		})
		return true
	}

	return false
}
