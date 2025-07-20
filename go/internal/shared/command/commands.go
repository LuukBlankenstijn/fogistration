package command

import (
	"encoding/json"
	"fmt"
)

type CommandType string

const (
	SetIpCmd  CommandType = "set_ip"
	SyncDjCmd CommandType = "sync_dj"
)

type Command interface {
	Type() CommandType
}

type SetIpCommand struct {
	Id int32  `json:"id"`
	Ip string `json:"ip"`
}

func (c SetIpCommand) Type() CommandType { return SetIpCmd }

type SyncDjCommand struct{}

func (c SyncDjCommand) Type() CommandType { return SyncDjCmd }

func ParseCommand(cmdType string, payload []byte) (Command, error) {
	switch CommandType(cmdType) {
	case SetIpCmd:
		var cmd SetIpCommand
		err := json.Unmarshal(payload, &cmd)
		return cmd, err
	case SyncDjCmd:
		var cmd SyncDjCommand
		err := json.Unmarshal(payload, &cmd)
		return cmd, err
	default:
		return nil, fmt.Errorf("unknown processor type: %s", cmdType)
	}
}
