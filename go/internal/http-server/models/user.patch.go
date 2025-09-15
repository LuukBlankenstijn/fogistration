package models

import "github.com/LuukBlankenstijn/fogistration/internal/shared/database/models"

type UserPatch struct {
	Username *string          `json:"username,omitempty"`
	Email    *string          `json:"email,omitempty"`
	Role     *models.UserRole `json:"role,omitempty"`
}

func PatchUser(user User, patch UserPatch) UserPut {
	var newUser UserPut

	if patch.Role != nil {
		newUser.Role = *patch.Role
	} else {
		newUser.Role = user.Role
	}

	if patch.Username != nil {
		newUser.Username = *patch.Username
	} else {
		newUser.Username = user.Username
	}

	if patch.Email != nil {
		newUser.Email = *patch.Email
	} else {
		newUser.Email = user.Email
	}
	return newUser
}
