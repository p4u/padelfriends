package handlers

import (
	"net/http"
	"path"
	"strings"
)

const webappdir = "ui/dist"

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	var filePath string
	if r.URL.Path == "/app" || r.URL.Path == "/app/" {
		filePath = path.Join(webappdir, "index.html")
	} else {
		filePath = path.Join(webappdir, strings.TrimPrefix(path.Clean(r.URL.Path), "/app"))
	}
	// Serve the file using http.ServeFile
	http.ServeFile(w, r, filePath)
}
