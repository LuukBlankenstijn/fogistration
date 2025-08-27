package service

import (
	"context"
	"errors"

	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	// Returned when another team already uses this client's IP.
	ErrClientIPUsedByOtherTeam = errors.New("client is already associated to another team")
)

type teamService struct {
	q    *database.Queries
	pool *pgxpool.Pool
}

func newTeamService(pool *pgxpool.Pool) *teamService {
	return &teamService{
		q:    database.New(pool),
		pool: pool,
	}
}

// SetClient sets team's IP to client's IP; if client == nil, clears the team's IP.
// Preconditions: `team` exists; if client != nil, `client` exists (controller validated).
func (s *teamService) SetClient(ctx context.Context, team database.Team, client *database.Client) error {
	return withTx(ctx, s.pool, func(ctx context.Context, q *database.Queries) error {
		// Clear assignment
		if client == nil {
			if _, err := q.UpdateIp(ctx, database.UpdateIpParams{
				ExternalID: team.ExternalID,
				Ip:         database.PgTextFromString(nil),
			}); err != nil {
				logging.Error("SetTeamClient: clear team IP failed", err, "team", team.ExternalID)
				return err
			}
			return nil
		}

		// Ensure client's IP isn't used by another team (idempotent if same team)
		existing, err := q.GetTeamByIp(ctx, database.PgTextFromString(&client.Ip))
		switch {
		case err == nil:
			if existing.ExternalID != team.ExternalID {
				return ErrClientIPUsedByOtherTeam
			}
			// already assigned to this team — nothing to do
			return nil
		case errors.Is(err, pgx.ErrNoRows):
			// free to assign
		default:
			logging.Error("SetTeamClient: lookup team by IP failed", err, "ip", client.Ip)
			return err
		}

		// Assign IP to team
		if _, err := q.UpdateIp(ctx, database.UpdateIpParams{
			ExternalID: team.ExternalID,
			Ip:         database.PgTextFromString(&client.Ip),
		}); err != nil {
			// unique violation on teams(ip) -> concurrent assign
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				return ErrClientIPUsedByOtherTeam
			}
			logging.Error("SetTeamClient: set team IP failed", err, "team", team.ExternalID, "ip", client.Ip)
			return err
		}

		return nil
	})
}
