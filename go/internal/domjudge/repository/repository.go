package repository

import (
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/repository"
)

type RepositoryStore struct {
	ContestRepositry *repository.ContestRepository
	TeamRepository   *repository.TeamRepository
}

func NewRepositoryStore(queries *database.Queries) RepositoryStore {
	return RepositoryStore{
		ContestRepositry: repository.NewClientRepository(queries),
		TeamRepository:   repository.NewTeamRepository(queries),
	}
}
