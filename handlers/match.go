package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/p4u/padelfriends/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MatchHandler struct {
	GroupService *services.GroupService
	MatchService *services.MatchService
}

// POST /api/group/{name}/matches?password=SECRET
// Payload: { "player_ids": ["playerID1","playerID2","playerID3","playerID4"] }
func (h *MatchHandler) CreateMatch(w http.ResponseWriter, r *http.Request) {
	groupName := chi.URLParam(r, "name")

	if !checkGroupPassword(w, r, h.GroupService, groupName) {
		return
	}

	var payload struct {
		PlayerIDs []string `json:"player_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if len(payload.PlayerIDs) != 4 {
		http.Error(w, "Exactly 4 player IDs required", http.StatusBadRequest)
		return
	}

	var pids []primitive.ObjectID
	for _, pid := range payload.PlayerIDs {
		objID, err := parseObjectID(pid)
		if err != nil {
			http.Error(w, "Invalid player ID: "+pid, http.StatusBadRequest)
			return
		}
		pids = append(pids, objID)
	}

	match, err := h.MatchService.CreateMatch(r.Context(), groupName, pids)
	if err != nil {
		http.Error(w, "Error creating match: "+err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusCreated, match)
}

// POST /api/group/{name}/matches/{match_id}/results?password=SECRET
// Payload: { "score_team1": X, "score_team2": Y }
func (h *MatchHandler) SubmitResults(w http.ResponseWriter, r *http.Request) {
	groupName := chi.URLParam(r, "name")
	matchIDStr := chi.URLParam(r, "match_id")

	if !checkGroupPassword(w, r, h.GroupService, groupName) {
		return
	}

	matchID, err := parseObjectID(matchIDStr)
	if err != nil {
		http.Error(w, "Invalid match ID", http.StatusBadRequest)
		return
	}

	var payload struct {
		ScoreTeam1 int `json:"score_team1"`
		ScoreTeam2 int `json:"score_team2"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.MatchService.SubmitResults(r.Context(), matchID, payload.ScoreTeam1, payload.ScoreTeam2); err != nil {
		http.Error(w, "Error submitting results: "+err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

// GET /api/group/{name}/matches?password=SECRET
func (h *MatchHandler) ListMatches(w http.ResponseWriter, r *http.Request) {
	groupName := chi.URLParam(r, "name")
	if !checkGroupPassword(w, r, h.GroupService, groupName) {
		return
	}

	matches, err := h.MatchService.ListMatches(r.Context(), groupName)
	if err != nil {
		http.Error(w, "Error listing matches: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, matches)
}
