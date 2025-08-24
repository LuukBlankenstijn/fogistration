package models

import "github.com/LuukBlankenstijn/fogistration/internal/shared/database"

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type User struct {
	ID       int64  `json:"id" binding:"required"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Role     Role   `json:"role" binding:"required"`
}

func MapUser(user database.User) User {
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
