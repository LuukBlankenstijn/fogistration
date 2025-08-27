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
	ErrTeamAlreadyAssigned = errors.New("team already assigned to a client")
)

type clientService struct {
	q    *database.Queries
	pool *pgxpool.Pool
}

func newClientService(pool *pgxpool.Pool) *clientService {
	q := database.New(pool)
	return &clientService{q: q, pool: pool}
}

// AssignClientToTeam performs the reassignment inside a transaction.
func (s *clientService) AssignClientToTeam(ctx context.Context, client database.Client, team *database.Team) error {
	return withTx(ctx, s.pool, func(ctx context.Context, q *database.Queries) error {
		// 1) Clear old team (if any) that currently has this client's IP
		old, err := q.GetTeamByIp(ctx, database.PgTextFromString(&client.Ip))
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			logging.Error("AssignClientToTeam: failed to get old team by IP", err, "ip", client.Ip)
			return err
		}
		if err == nil {
			if _, err := q.UpdateIp(ctx, database.UpdateIpParams{
				Ip:         database.PgTextFromString(nil),
				ExternalID: old.ExternalID,
			}); err != nil {
				logging.Error("AssignClientToTeam: failed to clear old team IP", err, "team", old.ExternalID)
				return err
			}
		}

		// 2) If target team provided, ensure free & set its IP to client's IP
		if team != nil {
			// quick in-memory check to fail fast (still race-safe thanks to unique index handling below)
			if team.Ip.Valid {
				return ErrTeamAlreadyAssigned
			}

			if _, err := q.UpdateIp(ctx, database.UpdateIpParams{
				Ip:         database.PgTextFromString(&client.Ip),
				ExternalID: team.ExternalID,
			}); err != nil {
				// If a unique constraint on teams(ip) fires, map to domain error
				var pgErr *pgconn.PgError
				if errors.As(err, &pgErr) && pgErr.Code == "23505" {
					return ErrTeamAlreadyAssigned
				}
				logging.Error("AssignClientToTeam: failed to set new team IP", err, "team", team.ExternalID, "ip", client.Ip)
				return err
			}
		}

		return nil
	})
}
