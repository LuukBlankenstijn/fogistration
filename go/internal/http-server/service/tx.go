package service

import (
	"context"

	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// withTx runs the given function inside a transaction, automatically rolling
// back on error/panic and committing on success.
func withTx(
	ctx context.Context,
	pool *pgxpool.Pool,
	fn func(ctx context.Context, q *database.Queries) error,
) error {
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		logging.Error("WithTx: begin failed", err)
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	q := database.New(pool).WithTx(tx)

	if err := fn(ctx, q); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		logging.Error("WithTx: commit failed", err)
		return err
	}
	return nil
}
