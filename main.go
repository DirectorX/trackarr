package main

import (
	"gitlab.com/cloudb0x/trackarr/autodl"
	"gitlab.com/cloudb0x/trackarr/cache"
	"gitlab.com/cloudb0x/trackarr/config"
	"gitlab.com/cloudb0x/trackarr/database"
	"gitlab.com/cloudb0x/trackarr/logger"
	"gitlab.com/cloudb0x/trackarr/pvr"
	"gitlab.com/cloudb0x/trackarr/runtime"
	"gitlab.com/cloudb0x/trackarr/tasks"
	"gitlab.com/cloudb0x/trackarr/tracker"
	"gitlab.com/cloudb0x/trackarr/version"
	"gitlab.com/cloudb0x/trackarr/web"

	stringutils "gitlab.com/cloudb0x/trackarr/utils/strings"

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
	if err := logger.Init(config.Runtime.Verbose, config.Runtime.Log); err != nil {
		log.WithError(err).Fatal("Failed to initialize logging")
	}

	log = logger.GetLogger("app")

	// Version info
	if err := version.Init(buildConfig); err != nil {
		log.WithError(err).Fatal("Failed to initialize version")
	}

	log.Infof("Using %s = %s (%s@%s)", stringutils.StringLeftJust("VERSION", " ", 10),
		buildConfig.Version, buildConfig.GitCommit, buildConfig.Timestamp)

	// Logging info
	logger.ShowUsing()

	// Init Config
	if err := config.Init(buildConfig); err != nil {
		log.WithError(err).Fatal("Failed to initialize config")
	}

	// Init Database
	if err := database.Init(); err != nil {
		log.WithError(err).Fatal("Failed to initialize database")
	}

	// Init Autodl
	if err := autodl.Init(); err != nil {
		log.WithError(err).Fatal("Failed to initialize autodl")
	}

	// Init PVR
	if err := pvr.Init(); err != nil {
		log.WithError(err).Fatal("Failed to initialize PVRs")
	}

	// Init Tracker
	if err := tracker.Init(); err != nil {
		log.WithError(err).Fatal("Failed to initialize trackers")
	}

	// Init Cache
	if err := cache.Init(); err != nil {
		log.WithError(err).Fatal("Failed initializing cache")
	}

	// Init Task Scheduler
	runtime.Tasks = tasks.New()
	if err := runtime.Tasks.Init(); err != nil {
		log.WithError(err).Fatal("Failed initializing task scheduler")
	}
}

/* Main */
func main() {
	log.Info("Initialized core")

	// Defer de-inits
	defer cache.Close()

	// Alert user when no trackers were loaded
	if len(runtime.Tracker) < 1 {
		log.Warn("No trackers were enabled/loaded...")
	}

	// Start web
	web.Listen(config.Config, config.Runtime.Verbose)

	// Start IRC clients
	initIRC()

	// Startup checks
	startupChecks()

	// Startup scheduled tasks
	runtime.Tasks.Start()

	// Wait for shutdown
	waitShutdown()
}
