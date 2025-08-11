package models

import "github.com/LuukBlankenstijn/fogistration/internal/shared/database"

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     Role   `json:"role"`
}

func MapUser(user database.AppUser) User {
	role := Role(user.Role)
	if role != RoleUser && role != RoleAdmin {
		role = RoleUser
	}
	return User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     role,
	}
}
