package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/LuukBlankenstijn/fogistration/internal/httpServer/models"
	"github.com/go-chi/chi/v5"
)

type SetTeamRequest struct {
	TeamId *string `json:"teamId"`
}

// @Summary List clients
// @Description Get all clients
// @Tags client
// @Accept json
// @Produce json
// @Success 200 {array} models.Client
// @Router /client [get]
// @Id GetAllClients
func (s *Server) ListClients(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	clients, err := s.queries.GetAllClients(ctx)
	if err != nil {
		http.Error(w, "could not get clients", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(models.MapClient(clients...))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// @Summary Set team
// @Description Sets the team a client is associated to
// @Tags client
// @Accept json
// @Success 200
// @Param clientId path int true "User ID"
// @Param request body SetTeamRequest true "team id"
// @Router /client/{clientId}/team [put]
// @Id SetTeam
func (s *Server) SetTeam(w http.ResponseWriter, r *http.Request) {
	var req SetTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	clientId := chi.URLParam(r, "clientId")
	id, err := strconv.Atoi(clientId)
	if err != nil {
		http.Error(w, "Malformed client id", http.StatusBadRequest)
		return
	}

	statusCode, errorMessage := s.Client.SetTeam(ctx, id, req.TeamId)
	if statusCode != http.StatusOK {
		http.Error(w, errorMessage, statusCode)
	}

}
