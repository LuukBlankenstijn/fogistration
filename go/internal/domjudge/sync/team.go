package syncer

import (
	"fmt"

	"github.com/LuukBlankenstijn/fogistration/internal/domjudge/client"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *DomJudgeSyncer) syncTeams() error {
	nextContest, err := s.ContestRepositry.GetNextOrActive(s.ctx)
	if err != nil {
		return fmt.Errorf("could not get next contest: %w", err)
	}
	if nextContest == nil {
		return nil
	}
	apiTeams, err := s.client.ListTeams(s.ctx, &client.GetV4AppApiTeamListParams{}, nextContest.ExternalID)
	if err != nil {
		return fmt.Errorf("could not get teams from api: %w", err)
	}

	// get hashes
	hashes, err := s.ContestRepositry.GetHashes(s.ctx)
	if err != nil {
		return fmt.Errorf("could not get team hashes from db: %w", err)
	}

	// mapping Id => hash
	hashMap := make(map[int32]string)
	for _, row := range hashes {
		hashMap[row.ID] = row.Hash
	}

	teams := []database.Team{}

	// mapping externalId => ip
	ipMap, err := s.getTeamIpMap()
	if err != nil {
		return fmt.Errorf("failed to get ip map: %w", err)
	}
	for _, apiContest := range apiTeams {
		team := mapTeamToDb(apiContest)
		if ip, exists := ipMap[team.ExternalID]; exists {
			team.Ip = pgtype.Text{
				String: ip,
				Valid:  true,
			}
		}

		if existingHash, exists := hashMap[team.ID]; !exists || existingHash != team.Hash {
			err = s.TeamRepository.Upsert(s.ctx, database.UpsertTeamParams{
				ID:          team.ID,
				ExternalID:  team.ExternalID,
				Name:        team.Name,
				DisplayName: team.DisplayName,
				Ip:          team.Ip,
				Hash:        team.Hash,
			})
			if err != nil {
				logging.Error("failed to upsert contest: %w", err)
			}
		}
		teams = append(teams, team)
	}

	err = s.ContestRepositry.DeleteAllTeams(s.ctx, nextContest.ID)
	if err != nil {
		return fmt.Errorf("failed to delete all teams for contest %d: %w", nextContest.ID, err)
	}

	err = s.ContestRepositry.InsertContestTeams(s.ctx, nextContest.ID, teams)
	if err != nil {
		return fmt.Errorf("failed to add update teams for contest %d: %w", nextContest.ID, err)
	}

	return nil
}

func (s DomJudgeSyncer) getTeamIpMap() (map[string]string, error) {
	users, err := s.client.ListUsers(s.ctx, &client.GetV4AppApiUserListParams{})
	if err != nil {
		return nil, err
	}

	// map from team id to ip
	ipMap := map[string]string{}
	for _, user := range users {
		if user.TeamId == nil || user.Ip == nil {
			continue
		}
		ipMap[*user.TeamId] = *user.Ip
	}

	return ipMap, nil
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
		Name:        c.Name,
		DisplayName: displayName,
	}

	hash := computeHash(team)
	team.Hash = hash

	return team
}
