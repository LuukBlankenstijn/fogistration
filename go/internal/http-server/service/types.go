package service

import (
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
)

type TokenIssuer interface {
	IssueJWT(u database.User) (string, error)
}
