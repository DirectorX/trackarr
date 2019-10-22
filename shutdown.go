package main

import (
	"context"
	"github.com/l3uddz/trackarr/tasks"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/l3uddz/trackarr/database"
	"github.com/l3uddz/trackarr/runtime"
)

func waitShutdown() {
	/* wait for shutdown signal */
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Warn("Shutting down...")

	shutdown()
}

func shutdown() {
	// Graceful shutdown IRC clients
	for _, ircClient := range runtime.Irc {
		ircClient.Stop()
		log.Infof("Stopped IRC client for %s", ircClient.Tracker.Name)
	}

	// Stop Web server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := runtime.Web.Shutdown(ctx); err != nil {
		log.WithError(err).Fatalf("Failed shutting down")
	}
	log.Info("Stopped web server")

	// Stop scheduled tasks
	tasks.Stop()

	// Close DB
	if err := database.DB.Close(); err != nil {
		log.WithError(err).Errorf("Failed closing database connection...")
	}
	log.Info("Stopped database")

	log.Info("Finished")
}
