package http

import (
	"encoding/json"
	"net/http"
)

// @Summary Check API health
// @Description Get API health status
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} ServerHealthResponse
// @Router /health [get]
func (*Server) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := ServerHealthResponse{
		Status: "ok",
	}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
