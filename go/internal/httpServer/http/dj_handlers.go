package http

import (
	"encoding/json"
	"net/http"

	"github.com/LuukBlankenstijn/fogistration/internal/httpServer/models"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"github.com/go-chi/chi/v5"
)

type PutClientRequest struct {
	ClientId *int32 `json:"clientId"`
}

// @Summary List teams
// @Description Get all teams
// @Tags dj
// @Accept json
// @Produce json
// @Success 200 {array} models.Team
// @Router /dj/team [get]
// @Id GetAllTeams
func (s *Server) ListTeams(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	teams, err := s.queries.GetAllTeams(ctx)
	if err != nil {
		http.Error(w, "could not get teams", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(models.MapTeam(teams...))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// @Summary List contests
// @Description Get all contests
// @Tags dj
// @Accept json
// @Produce json
// @Success 200 {array} models.Contest
// @Router /dj/contest [get]
// @Id GetAllContests
func (s *Server) ListContests(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	contests, err := s.queries.ListContests(ctx)
	if err != nil {
		http.Error(w, "could not get contests", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(models.MapContest(contests...))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// @Summary Active contest
// @Description Get the next or current active contest
// @Tags dj
// @Accept json
// @Produce json
// @Success 200 {object} models.Contest
// @Router /dj/contest/active [get]
// @Id GetActiveContest
func (s *Server) GetActiveContest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	contest, err := s.queries.GetNextOrActiveContest(ctx)
	if err != nil {
		http.Error(w, "could not get active contest", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(models.MapContest(contest)[0])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// @Summary Set team client
// @Description Set the client a team is associated to
// @Tags dj
// @Accept json
// @Produce json
// @Param teamId path string true "Team ID"
// @Param request body PutClientRequest true "Client id"
// @Success 200 {object} models.Team
// @Failure 400 {object} ErrorResponse
// @Router /dj/team/{teamId}/client [put]
// @Id SetTeamClient
func (s *Server) SetTeamClient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req PutClientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// get team
	teamId := chi.URLParam(r, "teamId")
	if teamId == "" {
		http.Error(w, "Malformed teamId id", http.StatusBadRequest)
		return
	}
	team, err := s.queries.GetTeamByExternalId(ctx, teamId)
	if err != nil {
		http.Error(w, "Team not found", http.StatusBadRequest)
		return
	}

	// if client id is null, just set the ip to rull and return
	if req.ClientId == nil {
		team, err = s.queries.UpdateIp(ctx, database.UpdateIpParams{
			ExternalID: team.ExternalID,
			Ip:         database.PgTextFromString(nil),
		})
		if err != nil {
			logging.Error("Failed to update ip in db", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(team)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// check if client exists
	client, err := s.queries.GetClientById(ctx, *req.ClientId)
	if err != nil {
		http.Error(w, "Client not found", http.StatusBadRequest)
		return
	}

	// check if client is not already used by another team
	if _, err := s.queries.GetTeamByIp(ctx, database.PgTextFromString(&client.Ip)); err == nil {
		http.Error(w, "Client is already associated to another team", http.StatusBadRequest)
		return
	}

	// set ip to that of client
	team, err = s.queries.UpdateIp(ctx, database.UpdateIpParams{
		ExternalID: team.ExternalID,
		Ip:         database.PgTextFromString(&client.Ip),
	})
	if err != nil {
		logging.Error("Failed to update ip in db", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(team)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
