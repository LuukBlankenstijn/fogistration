package models

import "github.com/LuukBlankenstijn/fogistration/internal/shared/database/models"

type UserPut struct {
	Username string          `json:"username"`
	Email    string          `json:"email"`
	Role     models.UserRole `json:"role"`
}
