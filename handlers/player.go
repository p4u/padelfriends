package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/p4u/padelfriends/services"
)

type PlayerHandler struct {
	GroupService  *services.GroupService
	PlayerService *services.PlayerService
}

// POST /api/group/{name}/players?password=SECRET
// Payload: { "name": "PlayerName" }
func (h *PlayerHandler) AddPlayer(w http.ResponseWriter, r *http.Request) {
	groupName := chi.URLParam(r, "name")
	// Check password
	if !checkGroupPassword(w, r, h.GroupService, groupName) {
		return
	}

	var payload struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if payload.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	player, err := h.PlayerService.AddPlayer(r.Context(), groupName, payload.Name)
	if err != nil {
		http.Error(w, "Error adding player: "+err.Error(), http.StatusConflict)
		return
	}

	writeJSON(w, http.StatusCreated, player)
}

// GET /api/group/{name}/players?password=SECRET
func (h *PlayerHandler) ListPlayers(w http.ResponseWriter, r *http.Request) {
	groupName := chi.URLParam(r, "name")

	// Check password
	if !checkGroupPassword(w, r, h.GroupService, groupName) {
		return
	}

	players, err := h.PlayerService.ListPlayers(r.Context(), groupName)
	if err != nil {
		http.Error(w, "Error listing players: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, players)
}
