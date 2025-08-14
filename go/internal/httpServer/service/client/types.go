package client

import "context"

type Service interface {
	SetTeam(ctx context.Context, id int, teamId *string) (int, string)
}
