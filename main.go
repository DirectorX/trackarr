package main

import (
	"github.com/l3uddz/trackarr/autodl"
	"github.com/l3uddz/trackarr/autodl/ircclient"
	"github.com/l3uddz/trackarr/autodl/parser"
	"github.com/l3uddz/trackarr/config"
	"github.com/l3uddz/trackarr/database"
	"github.com/l3uddz/trackarr/logger"
	"github.com/sirupsen/logrus"
)

var (
	// Logging
	log *logrus.Entry
)

/* Init */
func init() {
	// Setup cmd flags
	cmdInit()

	// Init Logging
	if err := logger.Init(flagLogLevel, flagLogPath); err != nil {
		log.WithError(err).Fatal("Failed to initialize logging")
	}

	log = logger.GetLogger("app")

	// Init Config
	if err := config.Init(flagConfigPath); err != nil {
		log.WithError(err).Fatal("Failed to initialize config")
	}

	// Init Database
	if err := database.Init(flagDbPath); err != nil {
		log.WithError(err).Fatal("Failed to initialize database")
	}

	// Test parser
	p, err := parser.Init("IPTorrents", flagTrackerPath)
	if err != nil {
		log.Fatal("Failed initializing tracker")
	} else {
		log.Info("Initialized tracker")
		log.Info(p.Tracker)
	}

	// Test irc
	client, err := ircclient.Init(p)
	if err != nil {
		log.Fatal("Failed initializing tracker irc client")
	} else {
		log.Info("Initialized tracker irc client")
	}
	client.Start()
	logrus.Exit(0)

	// Init Autodl
	if err := autodl.Init(flagTrackerPath); err != nil {
		log.WithError(err).Fatal("Failed to initialize autodl")
	}
}

/* Main */
func main() {
	log.Info("Initialized")

}
