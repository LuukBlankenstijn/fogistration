package models

import (
	"time"

	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
)

type Contest struct {
	ExternalID string    `json:"id" binding:"required"`
	FormalName string    `json:"name" binding:"required"`
	StartTime  time.Time `json:"startTime" binding:"required"`
	EndTime    time.Time `json:"endTime" binding:"required"`
}

func MapContest(contests ...database.Contest) []Contest {
	newContests := []Contest{}
	for _, contest := range contests {
		c := Contest{
			ExternalID: contest.ExternalID,
			FormalName: contest.FormalName,
			StartTime:  contest.StartTime.Time,
			EndTime:    contest.EndTime.Time,
		}
		newContests = append(newContests, c)
	}
	return newContests
}
