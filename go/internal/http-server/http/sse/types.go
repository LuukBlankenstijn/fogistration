package sse

import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
)

type SSEOperation string

const (
	Create SSEOperation = "create"
	Update SSEOperation = "update"
	Delete SSEOperation = "delete"
)

var SSEOperationValues = []SSEOperation{
	Create,
	Update,
	Delete,
}

func (SSEOperation) Schema(r huma.Registry) *huma.Schema {
	if r.Map()["SSEOperation"] == nil {
		schemaRef := r.Schema(reflect.TypeOf(""), true, "SSEOperation")
		schemaRef.Title = "SSEOperation"
		for _, v := range SSEOperationValues {
			schemaRef.Enum = append(schemaRef.Enum, string(v))
		}
		r.Map()["SSEOperation"] = schemaRef
	}
	return &huma.Schema{Ref: "#/components/schemas/SSEOperation"}
}

type SSEUpdate[T any] struct {
	Operation SSEOperation
	Id        int
	Data      *T
}

type GetResponse[T any] struct {
	Body T
}
