package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitlab.com/cloudb0x/trackarr/database"
	"gitlab.com/cloudb0x/trackarr/runtime"
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
		log.WithError(err).Fatal("Failed shutting down")
	}
	log.Info("Stopped web server")

	// Stop scheduled tasks
	runtime.Tasks.Stop()

	// Close DB
	if err := database.DB.Close(); err != nil {
		log.WithError(err).Error("Failed closing database connection...")
	}
	log.Info("Stopped database")

	// Stop logs processor
	if err := runtime.Loghook.Stop(); err != nil {
		log.WithError(err).Fatal("Failed shutting down loghook")
	}
	log.Info("Stopped loghook")

	log.Info("Finished")
}
