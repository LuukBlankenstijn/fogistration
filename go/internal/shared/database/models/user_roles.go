package models

import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
)

type UserRole string

const (
	Admin UserRole = "admin"
	User  UserRole = "user"
	Guest UserRole = "guest"
)

var UserRoleValues = []UserRole{
	Admin,
	User,
	Guest,
}

func (UserRole) Schema(r huma.Registry) *huma.Schema {
	if r.Map()["UserRole"] == nil {
		schemaRef := r.Schema(reflect.TypeOf(""), true, "UserRole")
		schemaRef.Title = "UserRole"
		for _, v := range UserRoleValues {
			schemaRef.Enum = append(schemaRef.Enum, string(v))
		}
		r.Map()["UserRole"] = schemaRef
	}
	return &huma.Schema{Ref: "#/components/schemas/UserRole"}
}
