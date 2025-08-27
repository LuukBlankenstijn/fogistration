package models

import (
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
)

type Team struct {
	ExternalID string `json:"id"`
	Name       string `json:"name"`
	Ip         string `json:"ip,omitempty"`
}

func MapTeam(teams ...database.Team) []Team {
	newteams := []Team{}
	for _, team := range teams {
		name := team.DisplayName.String
		if !team.DisplayName.Valid {
			name = team.Name
		}
		t := Team{
			ExternalID: team.ExternalID,
			Name:       name,
		}

		if team.Ip.Valid {
			t.Ip = team.Ip.String
		}
		newteams = append(newteams, t)
	}
	return newteams
}
