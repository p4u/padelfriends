package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/p4u/padelfriends/handlers"
)

func New(
	groupHandler *handlers.GroupHandler,
	playerHandler *handlers.PlayerHandler,
	matchHandler *handlers.MatchHandler,
	statsHandler *handlers.StatsHandler,
) http.Handler {

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// API routes
	r.Route("/api", func(r chi.Router) {
		// Public endpoints (no auth required)
		r.Get("/groups", groupHandler.ListGroups)
		r.Get("/group/{name}", groupHandler.GetGroupByName)
		r.Get("/group/byname/{name}", groupHandler.GetGroupByName)

		r.Route("/group/{name}", func(r chi.Router) {
			// Public endpoints (no auth required)
			r.Get("/matches", matchHandler.ListMatches)
			r.Get("/players", playerHandler.ListPlayers)
			r.Get("/statistics", statsHandler.GetStatistics)
			r.Get("/export/csv", groupHandler.ExportGroupMatchesCSV)

			// Authentication endpoint
			r.Post("/authenticate", groupHandler.AuthenticateGroup)

			// Protected endpoints (auth required)
			r.Group(func(r chi.Router) {
				r.Use(requireAuth)
				r.Post("/players", playerHandler.AddPlayer)
				r.Post("/matches", matchHandler.CreateMatch)
				r.Post("/matches/batch", matchHandler.CreateMatches)
				r.Post("/matches/{match_id}/cancel", matchHandler.CancelMatch)
				r.Post("/matches/{match_id}/results", matchHandler.SubmitResults)
			})
		})

		// Health check
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		})
	})

	// Serve static files
	r.HandleFunc("/*", handlers.StaticHandler)
	r.HandleFunc("/", handlers.StaticHandler)

	return r
}

// requireAuth middleware checks for password in query params
func requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		password := r.URL.Query().Get("password")
		if password == "" {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		groupName := chi.URLParam(r, "name")
		if groupName == "" {
			http.Error(w, "Group name required", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}
