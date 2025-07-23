package processor

import (
	"context"
	"fmt"

	"github.com/LuukBlankenstijn/fogistration/internal/cmdhandler/client"
	dbObject "github.com/LuukBlankenstijn/fogistration/internal/shared/database/object"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
)

func (c *CommandHandler) handleSetIpCommand(ctx context.Context, cmd dbObject.ChangeIp) error {
	logging.Info("setting ip")
	users, err := c.client.ListUsers(ctx, &client.GetV4AppApiUserListParams{})
	if err != nil {
		return fmt.Errorf("failed to get users: %w", err)
	}

	team, err := c.queries.GetTeamById(ctx, cmd.Id)
	if err != nil {
		return fmt.Errorf("failed to get team: %w", err)
	}

	var filteredUsers []client.User
	for _, user := range users {
		if user.TeamId != nil && team.ExternalID == *user.TeamId {
			filteredUsers = append(filteredUsers, user)
		}
	}

	for _, user := range filteredUsers {
		var ip string
		if cmd.Ip != nil {
			ip = *cmd.Ip
		} else {
			ip = ""
		}
		params := client.UpdateUser{
			Ip:    &ip,
			Roles: &[]string{},
		}
		_, err = c.client.UpdateUser(ctx, *user.Id, params)
		if err != nil {
			logging.Error("failed update user ip", err)
		}
	}

	return nil
}
