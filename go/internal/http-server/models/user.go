package models

import (
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database/models"
)

type User struct {
	ID       int32           `json:"id"`
	Username string          `json:"username"`
	Email    string          `json:"email"`
	Role     models.UserRole `json:"role"`
}

func MapUsers(users ...database.User) []User {
	var newUsers []User
	for _, user := range users {
		role := models.UserRole(user.Role)
		if role != models.User && role != models.Admin && role != models.Guest {
			role = models.User
		}

		newUsers = append(newUsers, User{
			ID:       int32(user.ID),
			Username: user.Username,
			Email:    user.Email,
			Role:     role,
		})
	}
	return newUsers
}
