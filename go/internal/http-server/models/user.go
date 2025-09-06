package models

import (
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database/models"
)

type User struct {
	ID       int64           `json:"id"`
	Username string          `json:"username"`
	Email    string          `json:"email"`
	Role     models.UserRole `json:"role"`
}

func MapUser(user database.User) User {
	role := models.UserRole(user.Role)
	if role != models.User && role != models.Admin && role != models.Guest {
		role = models.User
	}
	return User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     role,
	}
}
