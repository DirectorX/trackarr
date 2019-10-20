package main

import (
	"github.com/l3uddz/trackarr/autodl"
	"github.com/l3uddz/trackarr/autodl/parser"
	"github.com/l3uddz/trackarr/cache"
	"github.com/l3uddz/trackarr/config"
	"github.com/l3uddz/trackarr/database"
	"github.com/l3uddz/trackarr/ircclient"
	"github.com/l3uddz/trackarr/logger"
	"github.com/l3uddz/trackarr/release"
	"github.com/l3uddz/trackarr/web"
	"github.com/sirupsen/logrus"
)

var (
	// Build vars
	buildVersion   string
	buildTimestamp string
	buildGitCommit string

	// Logging
	log *logrus.Entry
)

/* Parse */
func init() {
	// Set build vars
	buildConfig := &config.BuildVars{
		Version:   buildVersion,
		GitCommit: buildGitCommit,
		Timestamp: buildTimestamp,
	}

	// Setup cmd flags
	cmdInit()

	// Init Logging
	if err := logger.Init(flagLogLevel, flagLogPath); err != nil {
		log.WithError(err).Fatal("Failed to initialize logging")
	}

	log = logger.GetLogger("app")

	// Init Config
	if err := config.Init(buildConfig, flagConfigPath); err != nil {
		log.WithError(err).Fatal("Failed to initialize config")
	}

	// Print and exit if version flag is set
	config.PrintVersion()
	if flagVersion {
		log.Logger.Exit(0)
	}

	// Init Database
	if err := database.Init(flagDbPath); err != nil {
		log.WithError(err).Fatal("Failed to initialize database")
	}

	// Init Autodl
	if err := autodl.Init(flagTrackerPath); err != nil {
		log.WithError(err).Fatal("Failed to initialize autodl")
	}

	// Init Release
	trackerCfg, err := config.GetAnyConfiguredTracker(&config.Config.Trackers)
	if err != nil {
		log.Fatal("You must configure at-least one tracker...")
	}

	if err := release.Init(&config.Config.Pvr, trackerCfg); err != nil {
		log.WithError(err).Fatal("Failed initializing release")
	}

	// Init Cache
	if err := cache.Init(); err != nil {
		log.WithError(err).Fatal("Failed initializing cache")
	}
}

/* Main */
func main() {
	log.Info("Initialized core")

	// defer de-inits
	defer cache.Close()

	// validate we have at-least one active tracker
	oneActive := false
	for _, tracker := range config.Config.Trackers {
		if tracker.Enabled {
			oneActive = true
			break
		}
	}

	if !oneActive {
		log.Fatalf("At-least one tracker must be enabled...")
	}

	// load trackers
	ircClients := make([]*ircclient.IRCClient, 0)
	connectedClients := 0

	log.Infof("Initializing trackers...")
	for trackerName, tracker := range config.Config.Trackers {
		// skip disabled trackers
		if !tracker.Enabled {
			log.Debugf("Skipping disabled tracker: %s", trackerName)
			continue
		}

		// parse tracker
		log.Debugf("Parsing tracker: %s", trackerName)
		t, err := parser.Parse(trackerName, flagTrackerPath)
		if err != nil {
			log.WithError(err).Fatalf("Failed parsing tracker: %s", trackerName)
			continue
		}
		log.Debugf("Parsed tracker: %s", trackerName)

		// validate required config settings were set for this tracker
		settingsFilled := true
		for _, trackerSetting := range t.Settings {
			if _, ok := tracker.Config[trackerSetting]; !ok {
				log.Warnf("Skipping tracker %s, missing config setting: %q", trackerName, trackerSetting)
				settingsFilled = false
				break
			}
		}

		if !settingsFilled {
			// there were missing config settings that were required by this tracker
			continue
		}

		// load irc client
		log.Debugf("Initializing irc client: %s", trackerName)
		c, err := ircclient.Init(t, tracker)
		if err != nil {
			log.WithError(err).Fatalf("Failed initializing irc client for tracker: %s", trackerName)
			continue
		}
		log.Debugf("Initialized irc client: %s", trackerName)

		// start client
		if err := c.Start(); err != nil {
			log.WithError(err).Errorf("Failed starting irc client for tracker: %s", trackerName)
			continue
		} else {
			// add client to slice
			ircClients = append(ircClients, c)
			connectedClients++
		}
	}

	// were there connected clients?
	if connectedClients < 1 {
		log.Fatal("Failed to establish a connection to any of the enabled trackers...")
	} else {
		log.Infof("Connected to %d trackers!", connectedClients)
	}

	// start web
	web.Listen(config.Config, flagLogLevel)

	// graceful shutdown non-web components
	for _, ircClient := range ircClients {
		ircClient.Stop()
	}

	if err := database.DB.Close(); err != nil {
		log.WithError(err).Errorf("Failed closing database connection...")
	}

	log.Info("Finished")
}
