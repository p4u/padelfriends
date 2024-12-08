package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/p4u/padelfriends/config"
	"github.com/p4u/padelfriends/db"
	"github.com/p4u/padelfriends/handlers"
	"github.com/p4u/padelfriends/router"
	"github.com/p4u/padelfriends/services"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to MongoDB
	mdb, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Initialize services
	groupService := services.NewGroupService(mdb.Database)
	playerService := services.NewPlayerService(mdb.Database)
	matchService := services.NewMatchService(mdb.Database)
	statsService := services.NewStatsService(mdb.Database)

	// Initialize handlers
	groupHandler := &handlers.GroupHandler{GroupService: groupService}
	playerHandler := &handlers.PlayerHandler{GroupService: groupService, PlayerService: playerService}
	matchHandler := &handlers.MatchHandler{GroupService: groupService, MatchService: matchService}
	statsHandler := &handlers.StatsHandler{GroupService: groupService, StatsService: statsService}

	// Create router
	r := router.New(groupHandler, playerHandler, matchHandler, statsHandler)

	// Start server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: r,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Starting server on port %d...", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	<-stop
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	log.Println("Server exited properly")
}
