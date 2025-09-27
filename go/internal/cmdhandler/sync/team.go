package syncer

import (
	"fmt"

	"github.com/LuukBlankenstijn/fogistration/internal/cmdhandler/client"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *DomJudgeSyncer) syncTeams(queries *database.Queries) error {
	contestRepo := repository.NewContestRepository(queries)

	nextContest, err := queries.GetNextOrActiveContest(s.ctx)
	if err == pgx.ErrNoRows {
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not get next contest: %w", err)
	}

	apiTeams, err := s.client.ListTeams(s.ctx, &client.GetV4AppApiTeamListParams{}, nextContest.ExternalID)
	if err != nil {
		return fmt.Errorf("could not get teams from api: %w", err)
	}

	// hashes from DB
	hashes, err := queries.GetTeamHashes(s.ctx)
	if err != nil {
		return fmt.Errorf("could not get team hashes from db: %w", err)
	}
	hashMap := make(map[int32]string, len(hashes))
	for _, row := range hashes {
		hashMap[row.ID] = row.Hash
	}

	teams := make([]database.Team, 0, len(apiTeams))
	for _, apiContest := range apiTeams {
		team := mapTeamToDb(apiContest) // hash no longer includes IP

		// Upsert WITHOUT touching ip
		if existingHash, exists := hashMap[team.ID]; !exists || existingHash != team.Hash {
			if err := queries.UpsertTeam(s.ctx, database.UpsertTeamParams{
				ID:          team.ID,
				ExternalID:  team.ExternalID,
				Name:        team.Name,
				DisplayName: team.DisplayName,
				Hash:        team.Hash,
			}); err != nil {
				logging.Error("failed to upsert team", err)
			}
		}
		teams = append(teams, team)
	}

	if err := queries.DeleteAllContestTeams(s.ctx, nextContest.ID); err != nil {
		return fmt.Errorf("failed to delete all teams for contest %d: %w", nextContest.ID, err)
	}
	if err := contestRepo.InsertContestTeams(s.ctx, nextContest.ID, teams); err != nil {
		return fmt.Errorf("failed to add/update teams for contest %d: %w", nextContest.ID, err)
	}
	return nil
}

func mapTeamToDb(c client.Team) database.Team {
	var displayName pgtype.Text
	if c.DisplayName != nil {
		displayName = pgtype.Text{
			String: *c.DisplayName,
			Valid:  true,
		}
	}
	team := database.Team{
		ID:          int32(*c.Teamid),
		ExternalID:  *c.Id,
		Name:        *c.Name,
		DisplayName: displayName,
	}

	hash := computeHash(team)
	team.Hash = hash

	return team
}
