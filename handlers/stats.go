package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/p4u/padelfriends/services"
)

type StatsHandler struct {
	GroupService *services.GroupService
	StatsService *services.StatsService
}

// GET /api/group/{name}/statistics
func (h *StatsHandler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	groupName := chi.URLParam(r, "name")

	stats, err := h.StatsService.ComputeStats(r.Context(), groupName)
	if err != nil {
		http.Error(w, "Error computing statistics: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, stats)
}
