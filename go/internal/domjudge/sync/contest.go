package syncer

import (
	"fmt"

	"github.com/LuukBlankenstijn/fogistration/internal/domjudge/client"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *DomJudgeSyncer) syncContests() error {
	onlyActive := false
	apiContests, err := s.client.ListContests(s.ctx, &client.GetV4AppApiContestListParams{OnlyActive: &onlyActive})
	if err != nil {
		return fmt.Errorf("could not get contests from api: %w", err)
	}

	hashes, err := s.ContestRepositry.GetHashes(s.ctx)
	if err != nil {
		return fmt.Errorf("could not get contest hashes from db: %w", err)
	}

	hashMap := make(map[int32]string)
	for _, row := range hashes {
		hashMap[row.ID] = row.Hash
	}

	for _, apiContest := range apiContests {
		contest := mapContestToDd(apiContest)
		if existingHash, exists := hashMap[contest.ID]; !exists || existingHash != contest.Hash {
			err = s.ContestRepositry.Upsert(s.ctx, database.UpsertContestParams{
				ID:         contest.ID,
				ExternalID: contest.ExternalID,
				FormalName: contest.FormalName,
				StartTime:  contest.StartTime,
				EndTime:    contest.EndTime,
				Hash:       contest.Hash,
			})
			if err != nil {
				logging.Error("failed to upsert contest: %w", err)
			}
		}
	}

	return nil
}

func mapContestToDd(c client.Contest) database.Contest {
	var startTime pgtype.Timestamp
	if c.StartTime != nil {
		startTime = pgtype.Timestamp{
			Time:  *c.StartTime,
			Valid: true,
		}
	}

	var endTime pgtype.Timestamp
	if c.EndTime != nil {
		endTime = pgtype.Timestamp{
			Time:  *c.EndTime,
			Valid: true,
		}
	}
	contest := database.Contest{
		ID:         int32(*c.Cid),
		ExternalID: *c.Id,
		FormalName: *c.FormalName,
		StartTime:  startTime,
		EndTime:    endTime,
	}

	hash := computeHash(contest)
	contest.Hash = hash

	return contest
}
