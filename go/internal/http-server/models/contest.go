package models

import (
	"time"

	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
)

type Contest struct {
	ID         int32     `json:"id"`
	ExternalID string    `json:"externalId"`
	FormalName string    `json:"name"`
	StartTime  time.Time `json:"startTime"`
	EndTime    time.Time `json:"endTime"`
}

func MapContest(contests ...database.Contest) []Contest {
	newContests := []Contest{}
	for _, contest := range contests {
		c := Contest{
			ID:         contest.ID,
			ExternalID: contest.ExternalID,
			FormalName: contest.FormalName,
			StartTime:  contest.StartTime.Time,
			EndTime:    contest.EndTime.Time,
		}
		newContests = append(newContests, c)
	}
	return newContests
}
