package client

import (
	"context"
	"net/http"

	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"github.com/jackc/pgx/v5"
)

type service struct {
	q *database.Queries
}

func New(q *database.Queries) Service {
	return &service{
		q: q,
	}
}

func (s *service) SetTeam(ctx context.Context, id int, teamId *string) (int, string) {
	// check if the new team does not already have a client assigned
	if teamId != nil {
		newTeam, err := s.q.GetTeamByExternalId(ctx, *teamId)
		if err != nil {
			logging.Error("could not get team from database", err)
			return http.StatusInternalServerError, ""
		}
		if newTeam.Ip.Valid {
			return http.StatusBadRequest, "new team does already have a client assigned"
		}
	}

	// get client
	client, err := s.q.GetClientById(ctx, int32(id))
	if err != nil {
		return http.StatusBadRequest, "could not find client"
	}

	// get the old team
	team, err := s.q.GetTeamByIp(ctx, database.PgTextFromString(&client.Ip))
	if err != nil && err != pgx.ErrNoRows {
		logging.Error("could not get team from database", err)
		return http.StatusInternalServerError, ""
	}

	// if the old team exists, set the ip to nil
	if err == nil {
		_, err = s.q.UpdateIp(ctx, database.UpdateIpParams{
			Ip:         database.PgTextFromString(nil),
			ExternalID: team.ExternalID,
		})
		if err != nil {
			logging.Error("failed to update ip in database", err)
		}
	}

	// if the new team exists, set ip to the ip of the client
	if teamId != nil {
		_, err = s.q.UpdateIp(ctx, database.UpdateIpParams{
			Ip:         database.PgTextFromString(&client.Ip),
			ExternalID: *teamId,
		})
		if err != nil {
			logging.Error("failed to update ip in database", err)
		}
	}

	return http.StatusOK, ""
}
