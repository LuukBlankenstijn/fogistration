package models

import (
	"time"

	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
)

type Client struct {
	ID       int32     `json:"id" binding:"required"`
	Ip       string    `json:"ip" binding:"required"`
	LastSeen time.Time `json:"lastSeen" binding:"required" format:"date-time" swaggertype:"string"`
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
