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

// GET /api/group/{id}/statistics?password=SECRET
func (h *StatsHandler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	groupIDStr := chi.URLParam(r, "id")
	groupID, err := parseObjectID(groupIDStr)
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}
	if !checkGroupPassword(w, r, h.GroupService, groupID) {
		return
	}

	stats, err := h.StatsService.ComputeStats(r.Context(), groupID)
	if err != nil {
		http.Error(w, "Error computing stats: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, stats)
}
