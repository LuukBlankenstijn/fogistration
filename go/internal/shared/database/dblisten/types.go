package dblisten

import "errors"

type Mixed struct {
	Table string
	Op    string // INSERT | UPDATE | DELETE
	New   any
	Old   any
}

type Notification[T any] struct {
	Table string
	Op    string
	New   *T
	Old   *T
}

var ErrNotStruct = errors.New("typ must be a struct")
