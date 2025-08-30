package notifications

import (
	"context"

	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database/dblisten"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
)

func Run(ctx context.Context, dbUrl string) error {
	l, err := dblisten.NewNotify(ctx, dbUrl)
	if err != nil {
		return err
	}
	defer l.Close(ctx)

	teams, err := dblisten.RegisterTyped[database.Team](ctx, l, "teams", 32)
	if err != nil {
		logging.Error("failed to register teams", err)
	}
	clients, err := dblisten.RegisterTyped[database.Client](ctx, l, "clients", 32)
	if err != nil {
		logging.Error("failed to register clients", err)
	}
	wallpapers, err := dblisten.RegisterTyped[database.Wallpaper](ctx, l, "wallpapers", 32)
	if err != nil {
		logging.Error("failed to register wallpapers", err)
	}

	for {
		select {
		case t, ok := <-teams:
			if !ok {
				teams = nil
				continue
			}
			logging.Info("received team notification: %+v, %+v", t.New, t.Old)
		case c, ok := <-clients:
			if !ok {
				clients = nil
				continue
			}
			logging.Info("received client notification: %+v, %+v", c.New, c.Old)
		case c, ok := <-wallpapers:
			if !ok {
				clients = nil
				continue
			}
			logging.Info("received wallpaper notification: %+v, %+v", c.New, c.Old)
		}
	}
}
