package users

import "github.com/LuukBlankenstijn/fogistration/internal/http-server/models"

type getUserRequest struct {
	ID int32 `path:"id" doc:"User ID"`
}

type getUserResponse struct {
	Body models.User
}

type listUsersResponse struct {
	Body []models.User
}

type putUserRequest struct {
	ID   int32 `path:"id" doc:"User ID"`
	Body models.UserPut
}

type patchUserRequest struct {
	ID   int32 `path:"id" doc:"User ID"`
	Body models.UserPatch
}
