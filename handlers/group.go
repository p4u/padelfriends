package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/p4u/padelfriends/services"
)

type GroupHandler struct {
	GroupService *services.GroupService
}

// POST /api/group
// Payload: { "name": "GroupName", "password": "secret" }
func (h *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if payload.Name == "" || payload.Password == "" {
		http.Error(w, "Missing name or password", http.StatusBadRequest)
		return
	}

	group, err := h.GroupService.CreateGroup(r.Context(), payload.Name, payload.Password)
	if err != nil {
		http.Error(w, "Failed to create group: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, group)
}

// GET /api/group/byname/{name}?password=SECRET
// Retrieves a group by name if password is correct

func (h *GroupHandler) GetGroupByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if name == "" {
		http.Error(w, "Missing group name", http.StatusBadRequest)
		return
	}

	password := getQueryParam(r, "password")
	if password == "" {
		http.Error(w, "Missing password", http.StatusBadRequest)
		return
	}

	g, err := h.GroupService.GetGroupByName(r.Context(), name)
	if err != nil {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	if !services.CheckPassword(password, g.PasswordHash) {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	writeJSON(w, http.StatusOK, g)
}

// Helper to extract the group name from URL if we define {name} as a route parameter.
func (h *GroupHandler) getGroupNameFromURL(r *http.Request) string {
	// If we define the route as /api/group/byname/{name}, we can use chi.URLParam:
	name := r.URL.Path // or use chi.URLParam(r, "name")
	// We'll rely on router configuration to extract name properly soon.
	return name
}
