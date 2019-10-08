package main

import (
	"github.com/l3uddz/trackarr/autodl"
	"github.com/l3uddz/trackarr/config"
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

	// Init Autodl
	if err := autodl.Init(flagTrackerPath); err != nil {
		log.WithError(err).Fatal("Failed to initialize autodl")
	}
}

/* Main */
func main() {
	log.Info("Initialized")
}
