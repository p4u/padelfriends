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

// GetGroupByName handles GET /api/group/byname/{name}?password=SECRET
// Retrieves a group by name if password is correct
func (h *GroupHandler) GetGroupByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if name == "" {
		writeError(w, http.StatusBadRequest, "Missing group name")
		return
	}

	password := getQueryParam(r, "password")
	if password == "" {
		writeError(w, http.StatusBadRequest, "Missing password")
		return
	}

	g, err := h.GroupService.GetGroupByName(r.Context(), name)
	if err != nil {
		writeError(w, http.StatusNotFound, "Group not found")
		return
	}

	if !services.CheckPassword(password, g.PasswordHash) {
		writeError(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	writeJSON(w, http.StatusOK, g)
}

// ListGroups handles GET /api/groups?password=SECRET
// Retrieves a list of all groups. This endpoint can be secured as needed.
func (h *GroupHandler) ListGroups(w http.ResponseWriter, r *http.Request) {
	password := r.Header.Get("X-Group-Password") // Securely retrieve password from headers
	if password == "" {
		writeError(w, http.StatusBadRequest, "Missing X-Group-Password header")
		return
	}

	// Optionally, verify password here if required for listing groups
	// For example, only allow listing if the password matches a master password
	// This depends on your application's authentication logic

	// For demonstration, assume a master password is "adminsecret"
	masterPassword := "adminsecret"
	if password != masterPassword {
		writeError(w, http.StatusUnauthorized, "Invalid master password")
		return
	}

	groups, err := h.GroupService.ListGroups(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error listing groups: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, groups)
}
