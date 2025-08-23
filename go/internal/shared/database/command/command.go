package command

import (
	"encoding/json"
	"fmt"
)

type CommandType string

const (
	SyncDjType CommandType = "sync_dj"
)

type Command interface {
	Type() CommandType
}

type SyncDj struct{}

func (c SyncDj) Type() CommandType { return SyncDjType }

func ParseDatabaseObject(t string, payload []byte) (Command, error) {
	switch CommandType(t) {
	case SyncDjType:
		var cmd SyncDj
		err := json.Unmarshal(payload, &cmd)
		return cmd, err
	default:
		return nil, fmt.Errorf("unknown processor type: %s", t)
	}
}
