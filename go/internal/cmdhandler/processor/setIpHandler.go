package processor

import (
	"context"
	"fmt"

	"github.com/LuukBlankenstijn/fogistration/internal/cmdhandler/client"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/command"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"github.com/jackc/pgx/v5/pgtype"
)

func (c *CommandHandler) handleSetIpCommand(ctx context.Context, cmd command.SetIpCommand) error {
	users, err := c.client.ListUsers(ctx, &client.GetV4AppApiUserListParams{})
	if err != nil {
		return fmt.Errorf("failed to get users: %w", err)
	}

	team, err := c.teamRepo.GetById(ctx, cmd.Id)
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
		params := client.UpdateUser{
			Ip:    &cmd.Ip,
			Roles: &[]string{},
		}
		_, err = c.client.UpdateUser(ctx, *user.Id, params)
		if err != nil {
			logging.Error("failed update user ip", err)
		}
	}

	err = c.teamRepo.UpdateIp(ctx, database.UpdateIpParams{
		ID: cmd.Id,
		Ip: pgtype.Text{
			String: cmd.Ip,
			Valid:  true,
		},
	})

	if err != nil {
		return fmt.Errorf("error updating ip in database: %w", err)
	}

	return nil
}
