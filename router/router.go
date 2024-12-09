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
		r.Get("/group/{name}", groupHandler.GetGroupByName)
		r.Get("/group/byname/{name}", groupHandler.GetGroupByName)
		r.Get("/groups", groupHandler.ListGroups)

		r.Route("/group/{name}", func(r chi.Router) {
			r.Post("/players", playerHandler.AddPlayer)
			r.Get("/players", playerHandler.ListPlayers)

			// Match routes
			r.Post("/matches", matchHandler.CreateMatch)
			r.Post("/matches/batch", matchHandler.CreateMatches)
			r.Post("/matches/{match_id}/cancel", matchHandler.CancelMatch)
			r.Post("/matches/{match_id}/results", matchHandler.SubmitResults)
			r.Get("/matches", matchHandler.ListMatches)

			r.Get("/statistics", statsHandler.GetStatistics)
		})

		// Health check
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		})
	})
	r.HandleFunc("/*", handlers.StaticHandler)
	r.HandleFunc("/", handlers.StaticHandler)

	return r
}
