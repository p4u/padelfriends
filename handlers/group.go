package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/p4u/padelfriends/services"
)

// GroupHandler handles group-related HTTP requests.
type GroupHandler struct {
	GroupService *services.GroupService
}

// CreateGroup handles POST /api/group
// Payload: { "name": "GroupName", "password": "secret" }
func (h *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if payload.Name == "" || payload.Password == "" {
		writeError(w, http.StatusBadRequest, "Missing name or password")
		return
	}

	group, err := h.GroupService.CreateGroup(r.Context(), payload.Name, payload.Password)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create group: "+err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, group)
}

// GetGroupByName handles GET /api/group/byname/{name}
func (h *GroupHandler) GetGroupByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if name == "" {
		writeError(w, http.StatusBadRequest, "Missing group name")
		return
	}

	g, err := h.GroupService.GetGroupByName(r.Context(), name)
	if err != nil {
		writeError(w, http.StatusNotFound, "Group not found")
		return
	}

	// Check password if provided for authentication status
	password := getQueryParam(r, "password")
	if password != "" && services.CheckPassword(password, g.PasswordHash) {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"name":            g.Name,
			"created_at":      g.CreatedAt,
			"isAuthenticated": true,
		})
		return
	}

	// Return basic info for unauthenticated requests
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"name":            g.Name,
		"created_at":      g.CreatedAt,
		"isAuthenticated": false,
	})
}

// ListGroups handles GET /api/groups
func (h *GroupHandler) ListGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := h.GroupService.ListGroups(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error listing groups: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, groups)
}

// AuthenticateGroup handles POST /api/group/{name}/authenticate
func (h *GroupHandler) AuthenticateGroup(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if name == "" {
		writeError(w, http.StatusBadRequest, "Missing group name")
		return
	}

	var payload struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	g, err := h.GroupService.GetGroupByName(r.Context(), name)
	if err != nil {
		writeError(w, http.StatusNotFound, "Group not found")
		return
	}

	if !services.CheckPassword(payload.Password, g.PasswordHash) {
		writeError(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"name":            g.Name,
		"created_at":      g.CreatedAt,
		"isAuthenticated": true,
	})
}

// ExportGroupMatchesCSV handles GET /api/group/{name}/export/csv
func (h *GroupHandler) ExportGroupMatchesCSV(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if name == "" {
		writeError(w, http.StatusBadRequest, "Missing group name")
		return
	}

	csv, err := h.GroupService.ExportGroupMatchesCSV(r.Context(), name)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error exporting matches: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename="+name+"-matches.csv")
	w.Write([]byte(csv))
}
