package eventhandler

import (
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database/dblisten"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/pb"
)

func (e *EventHandler) HandleIpChange(change dblisten.Notification[database.Team]) {
	reload := func(ip string) {
		message := pb.ServerMessage{
			Message: &pb.ServerMessage_Reload{
				Reload: &pb.Reload{},
			},
		}
		e.pubsub.Publish(ip, &message)
	}

	if change.Old != nil && change.New != nil && change.Old.Ip == change.New.Ip {
		return
	}

	if change.Old != nil && change.Old.Ip.Valid {
		reload(change.Old.Ip.String)
	}

	if change.New != nil && change.New.Ip.Valid {
		reload(change.New.Ip.String)
	}
}
