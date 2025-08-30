package common

type Mixed struct {
	Table string
	Op    string
	Old   any
	New   any
}

type Notification[T any] struct {
	Table string
	Op    string
	Old   *T
	New   *T
}
