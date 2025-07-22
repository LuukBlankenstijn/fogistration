package server

import (
	"context"
	"fmt"

	"github.com/LuukBlankenstijn/fogistration/internal/grpc/pubsub"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	pb "github.com/LuukBlankenstijn/fogistration/internal/shared/pb"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type registrationService struct {
	queries *database.Queries
	pubsub  *pubsub.Manager
}

func (r *registrationService) register(ctx context.Context, client database.Client) error {
	sendTeam := func(team database.Team) {
		msg := &pb.ServerMessage{
			Message: &pb.ServerMessage_SetTeam{
				SetTeam: &pb.SetTeam{
					Name:        team.Name,
					ImageUrl:    "",
					DisplayName: database.StringValueFromPgText(team.DisplayName),
				},
			},
		}
		r.pubsub.Publish(client.Ip, msg)
	}

	team, err := r.queries.GetTeamByIp(ctx, pgtype.Text{String: client.Ip, Valid: true})

	// err not nill, already registered
	if err == nil {
		sendTeam(team)
		return nil
	}

	if err != pgx.ErrNoRows {
		return fmt.Errorf("error when getting team from database: %w", err)
	}

	contest, err := r.queries.GetNextOrActiveContest(ctx)
	if err != nil {
		return fmt.Errorf("error when getting active contest form database")
	}
	team, err = r.queries.ClaimTeam(ctx, database.ClaimTeamParams{
		Ip:        database.PgTextFromString(client.Ip),
		ContestID: contest.ID,
	})
	if err == pgx.ErrNoRows {
		// no team available, do nothing
		logging.Info("no team found for client %s", client.Ip)
		return nil
	}
	if err != nil {
		return fmt.Errorf("error when claiming team: %w", err)
	}

	// team claimed, send and return
	sendTeam(team)
	return nil
}
