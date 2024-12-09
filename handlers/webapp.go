package handlers

import (
	"net/http"
	"os"
	"path"
	"strings"
)

const webappdir = "ui/dist"

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	// Clean the URL path
	cleanPath := path.Clean(strings.TrimPrefix(r.URL.Path, "/app"))
	filePath := path.Join(webappdir, cleanPath)

	// Check if the file exists
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		// If file doesn't exist, serve index.html for client-side routing
		http.ServeFile(w, r, path.Join(webappdir, "index.html"))
		return
	}

	// Serve the actual file if it exists
	http.ServeFile(w, r, filePath)
}
