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
		r.Post("/group", groupHandler.CreateGroup)
		r.Get("/group/byname/{name}", groupHandler.GetGroupByName)
		r.Get("/groups", groupHandler.ListGroups) // **Added this line**

		r.Route("/group/{id}", func(r chi.Router) {
			r.Post("/players", playerHandler.AddPlayer)
			r.Get("/players", playerHandler.ListPlayers)
			r.Post("/matches", matchHandler.CreateMatch)
			r.Post("/matches/{match_id}/results", matchHandler.SubmitResults)
			r.Get("/matches", matchHandler.ListMatches)
			r.Get("/statistics", statsHandler.GetStatistics)
		})

		// Health check
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		})
	})

	return r
}