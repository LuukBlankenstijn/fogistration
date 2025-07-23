package listener

import (
	"context"
	"fmt"

	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database/object"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/pb"
)

func (d *DatabaseListener) handleIpChange(ctx context.Context, obj object.DatabaseObject) error {
	change, ok := obj.(object.ChangeIp)
	if !ok {
		return fmt.Errorf("unexpected type for IpChangeType")
	}

	if change.Ip == nil && change.IpOld != nil {
		message := pb.ServerMessage{
			Message: &pb.ServerMessage_UnsetTeam{
				UnsetTeam: &pb.UnsetTeam{},
			},
		}
		d.pubsub.Publish(*change.IpOld, &message)
		logging.Info("unset team for client %s", *change.IpOld)
	}

	if change.Ip != nil {
		team, err := d.queries.GetTeamById(ctx, change.Id)
		if err != nil {
			return fmt.Errorf("failed to get team from db: %w", err)
		}
		message := pb.ServerMessage{
			Message: &pb.ServerMessage_SetTeam{
				SetTeam: &pb.SetTeam{
					Name:        team.Name,
					DisplayName: database.StringValueFromPgText(team.DisplayName),
					ImageUrl:    "",
				},
			},
		}

		d.pubsub.Publish(*change.Ip, &message)
		logging.Info("set team %s for client %s", team.Name, *change.Ip)
	}

	return nil
}
