package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/p4u/padelfriends/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// writeJSON writes the given data as JSON to the response writer.
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		// If encoding fails, log the error and write a generic error message
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// writeError writes an error message in JSON format to the response writer.
func writeError(w http.ResponseWriter, status int, message string) {
	response := map[string]string{"error": message}
	writeJSON(w, status, response)
}

// parseObjectID parses a string into a MongoDB ObjectID or returns an error.
func parseObjectID(s string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(s)
}

// getQueryParam retrieves a query parameter from the request.
func getQueryParam(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

// checkGroupPassword checks if the provided password matches the group’s password.
func checkGroupPassword(w http.ResponseWriter, r *http.Request, groupService *services.GroupService, groupID primitive.ObjectID) bool {
	password := getQueryParam(r, "password")
	if password == "" {
		http.Error(w, "Missing password query parameter", http.StatusBadRequest)
		return false
	}

	// Retrieve group by ID
	g, err := groupService.GetGroupByID(r.Context(), groupID)
	if err != nil {
		http.Error(w, "Group not found or error retrieving group", http.StatusNotFound)
		return false
	}

	valid := services.CheckPassword(password, g.PasswordHash)
	if !valid {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return false
	}
	return true
}
