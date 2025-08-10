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
	sendReload := func() {
		msg := &pb.ServerMessage{
			Message: &pb.ServerMessage_Reload{
				Reload: &pb.Reload{},
			},
		}
		r.pubsub.Publish(client.Ip, msg)
	}

	_, err := r.queries.GetTeamByIp(ctx, pgtype.Text{String: client.Ip, Valid: true})

	// err nill, already registered
	if err == nil {
		sendReload()
		return nil
	}

	if err != pgx.ErrNoRows {
		return fmt.Errorf("error when getting team from database: %w", err)
	}

	contest, err := r.queries.GetNextOrActiveContest(ctx)
	if err != nil {
		return fmt.Errorf("error when getting active contest form database")
	}
	_, err = r.queries.ClaimTeam(ctx, database.ClaimTeamParams{
		Ip:        database.PgTextFromString(&client.Ip),
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

	// don't need to send, database trigger handles that
	return nil
}
