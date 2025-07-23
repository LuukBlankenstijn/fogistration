package object

import (
	"encoding/json"
	"fmt"
)

type DatabaseObjectType string

const (
	IpChangeType DatabaseObjectType = "change_ip"
	SyncDjType   DatabaseObjectType = "sync_dj"
)

type DatabaseObject interface {
	Type() DatabaseObjectType
}

type ChangeIp struct {
	Id    int32   `json:"id"`
	Ip    *string `json:"ip"`
	IpOld *string `json:"ip_old"`
}

func (c ChangeIp) Type() DatabaseObjectType { return IpChangeType }

type SyncDj struct{}

func (c SyncDj) Type() DatabaseObjectType { return SyncDjType }

func ParseDatabaseObject(t string, payload []byte) (DatabaseObject, error) {
	switch DatabaseObjectType(t) {
	case IpChangeType:
		var cmd ChangeIp
		err := json.Unmarshal(payload, &cmd)
		return cmd, err
	case SyncDjType:
		var cmd SyncDj
		err := json.Unmarshal(payload, &cmd)
		return cmd, err
	default:
		return nil, fmt.Errorf("unknown processor type: %s", t)
	}
}
