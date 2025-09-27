package processor

import (
	"context"
	"fmt"

	"github.com/LuukBlankenstijn/fogistration/internal/cmdhandler/client"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database/dblisten"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
)

func (c *Worker) handleSetIpCommand(ctx context.Context, change dblisten.Notification[database.Team]) error {
	logging.Info("setting ip")
	users, err := c.client.ListUsers(ctx, &client.GetV4AppApiUserListParams{})
	if err != nil {
		return fmt.Errorf("failed to get users: %w", err)
	}

	if change.New == nil || change.Old == nil {
		logging.Info("operation is creation or deletion, not setting ip")
		return nil
	}

	team := change.New

	var filteredUsers []client.User
	for _, user := range users {
		if user.TeamId != nil && team.ExternalID == *user.TeamId {
			filteredUsers = append(filteredUsers, user)
		}
	}

	for _, user := range filteredUsers {
		params := client.UpdateUser{
			Ip:    &team.Ip.String,
			Roles: &[]string{},
		}
		_, err = c.client.UpdateUser(ctx, *user.Id, params)
		if err != nil {
			logging.Error("failed update ip", err, "teamname", team.Name)
		}
	}

	return nil
}
