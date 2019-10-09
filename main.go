package main

import (
	"github.com/l3uddz/trackarr/autodl"
	"github.com/l3uddz/trackarr/config"
	"github.com/l3uddz/trackarr/database"
	models "github.com/l3uddz/trackarr/database/models"
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

	// Init Autodl
	if err := autodl.Init(flagTrackerPath); err != nil {
		log.WithError(err).Fatal("Failed to initialize autodl")
	}
}

/* Main */
func main() {
	log.Info("Initialized")

	if test, err := models.NewOrExistingTracker(database.DB, "testing2"); err != nil {
		log.Fatal("Failed finding an existing tracker...")
	} else {
		test.Version = "3"
		database.DB.Save(test)
	}
}