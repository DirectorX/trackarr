package main

import (
	"github.com/l3uddz/trackarr/autodl"
	"github.com/l3uddz/trackarr/autodl/parser"
	"github.com/l3uddz/trackarr/config"
	"github.com/l3uddz/trackarr/database"
	"github.com/l3uddz/trackarr/ircclient"
	"github.com/l3uddz/trackarr/logger"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
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

/* Misc */

func waitForSignal() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	signal.Notify(sigs, syscall.SIGTERM)
	<-sigs
}

/* Main */
func main() {
	log.Info("Initialized core")

	// validate we have some configured trackers
	trackersCount := len(config.Config.Trackers)
	if trackersCount < 1 {
		log.Fatalf("You must configure at-least one tracker")
	}

	// load trackers
	ircClients := make([]*ircclient.IRCClient, 0)

	log.Infof("Initializing trackers...")
	for trackerName, tracker := range config.Config.Trackers {
		// load parser
		log.Debugf("Initializing parser: %s", trackerName)
		p, err := parser.Init(trackerName, flagTrackerPath)
		if err != nil {
			log.WithError(err).Fatalf("Failed initializing parser for tracker: %s", trackerName)
		}
		log.Debugf("Initialized parser: %s", trackerName)

		// TODO: validate tracker settings are configured (authkey / passkey / torrent_pass etc...)

		// load irc client
		log.Debugf("Initializing irc client: %s", trackerName)
		c, err := ircclient.Init(p, &tracker)
		if err != nil {
			log.WithError(err).Fatalf("Failed initializing irc client for tracker: %s", trackerName)
		}
		log.Debugf("Initialized irc client: %s", trackerName)

		// start client
		if err := c.Start(); err != nil {
			log.WithError(err).Errorf("Failed starting irc client for tracker: %s", trackerName)
		} else {
			// add client to slice
			ircClients = append(ircClients, c)
		}
	}

	// wait for shutdown signal
	waitForSignal()

	// graceful shutdown
	log.Info("Shutting down")
	for _, ircClient := range ircClients {
		ircClient.Stop()
	}

	log.Info("Finished")
}
