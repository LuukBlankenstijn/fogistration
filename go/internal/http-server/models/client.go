package models

import (
	"time"

	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
)

type Client struct {
	ID       int32     `json:"id"`
	Ip       string    `json:"ip"`
	LastSeen time.Time `json:"lastSeen"`
}

func MapClient(clients ...database.Client) []Client {
	newClients := []Client{}
	for _, client := range clients {
		newClients = append(newClients, Client{
			ID:       client.ID,
			Ip:       client.Ip,
			LastSeen: client.LastSeen.Time,
		})
	}
	return newClients
}
