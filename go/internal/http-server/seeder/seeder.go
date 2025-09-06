package seeder

import (
	"context"

	"github.com/LuukBlankenstijn/fogistration/internal/http-server/utils/auth"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database/models"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Seeder struct {
	db *pgxpool.Pool
	q  *database.Queries
}

func New(db *pgxpool.Pool, q *database.Queries) *Seeder {
	return &Seeder{db: db, q: q}
}

func (s *Seeder) SeedDefaultUser(ctx context.Context) error {
	count, err := s.q.CountUsers(ctx)
	if err != nil {
		return err
	}

	if count > 0 {
		logging.Info("Users already exist, skipping seeding.")
		return nil
	}
	logging.Info("seeding ")

	// Start transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	username := "admin"
	email := "admin@example.com"
	role := models.Admin
	password := "admin123"

	salt, hash, err := auth.NewSecret(password, auth.Default)
	if err != nil {
		return err
	}

	user, err := qtx.CreateLocalUser(ctx, database.CreateLocalUserParams{
		Username: username,
		Email:    email,
		Role:     role,
	})
	if err != nil {
		return err
	}

	if _, err := qtx.CreateAuthSecret(ctx, database.CreateAuthSecretParams{
		UserID:       user.ID,
		Salt:         salt,
		PasswordHash: hash,
	}); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	logging.Info("Seeded default user: {username: %s, password: %s}", username, password)
	return nil
}
