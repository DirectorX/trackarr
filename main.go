package main

import (
	"github.com/l3uddz/trackarr/autodl"
	"github.com/l3uddz/trackarr/cache"
	"github.com/l3uddz/trackarr/config"
	"github.com/l3uddz/trackarr/database"
	"github.com/l3uddz/trackarr/tasks"
	"github.com/l3uddz/trackarr/utils/version"
	"github.com/l3uddz/trackarr/web"

	"github.com/l3uddz/trackarr/logger"
	"github.com/l3uddz/trackarr/pvr"
	"github.com/l3uddz/trackarr/runtime"
	"github.com/l3uddz/trackarr/tracker"

	stringutils "github.com/l3uddz/trackarr/utils/strings"

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
	if err := logger.Init(flagLogLevel, flagLogFile); err != nil {
		log.WithError(err).Fatal("Failed to initialize logging")
	}

	log = logger.GetLogger("app")

	// Version info
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

	// Check version
	usingLatest, latestVersion := version.IsLatestGitlabVersion(
		"https://gitlab.com/api/v4/projects/15385789/releases", "", config.Build.Version)
	if !usingLatest {
		log.Warnf("You are not using the latest version: %s", latestVersion)
	}

	// Defer de-inits
	defer cache.Close()

	// Validate we have at-least one active tracker
	if len(runtime.Tracker) < 1 {
		log.Fatalf("At-least one tracker must be enabled...")
	}

	// Start web
	web.Listen(config.Config, flagLogLevel)

	// Start IRC clients
	initIRC()

	// Startup checks
	startupChecks()

	// Startup scheduled tasks
	runtime.Tasks.Start()

	// Wait for shutdown
	waitShutdown()
}
