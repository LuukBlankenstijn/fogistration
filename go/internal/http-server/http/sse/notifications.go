package sse

import (
	"context"

	"github.com/LuukBlankenstijn/fogistration/internal/http-server/models"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database/dblisten"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
)

func (s *SSEManager) Start(ctx context.Context, dbUrl string) error {
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
			var update SSEUpdate[models.Team]
			if t.Op == "delete" {
				update = SSEUpdate[models.Team]{
					Data:      nil,
					Id:        int(t.Old.ID),
					Operation: Delete,
				}
			} else if t.New != nil {
				update = SSEUpdate[models.Team]{
					Data:      &models.MapTeam(*t.New)[0],
					Id:        int(t.New.ID),
					Operation: Update,
				}
			} else {
				continue
			}
			s.Broadcast(update)
		case c, ok := <-clients:
			if !ok {
				clients = nil
				continue
			}
			var update SSEUpdate[models.Client]
			if c.Op == "delete" {
				update = SSEUpdate[models.Client]{
					Data:      nil,
					Id:        int(c.Old.ID),
					Operation: Delete,
				}
			} else if c.New != nil {
				update = SSEUpdate[models.Client]{
					Data:      &models.MapClient(*c.New)[0],
					Id:        int(c.New.ID),
					Operation: Update,
				}
			} else {
				continue
			}
			s.Broadcast(update)
		case w, ok := <-wallpapers:
			if !ok {
				wallpapers = nil
				continue
			}
			var update SSEUpdate[models.Wallpaper]
			if w.Op == "delete" {
				update = SSEUpdate[models.Wallpaper]{
					Data:      nil,
					Id:        int(w.Old.ID),
					Operation: Delete,
				}
			} else if w.New != nil {
				update = SSEUpdate[models.Wallpaper]{
					Data:      &models.MapWallpaper(*w.New)[0],
					Id:        int(w.New.ID),
					Operation: Update,
				}
			} else {
				continue
			}
			s.Broadcast(update)
		}
	}

}
