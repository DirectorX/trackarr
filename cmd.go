package main

import (
	"path/filepath"

	"github.com/l3uddz/trackarr/config"
	"github.com/l3uddz/trackarr/utils/paths"

	"github.com/spf13/pflag"
)

// CLI flags with defaults
var (
	flagLogLevel      int
	flagConfigFolder  = paths.GetCurrentBinaryPath()
	flagConfigFile    = "config.yaml"
	flagLogFile       = "activity.log"
	flagDbFile        = "vault.db"
	flagTrackerFolder = "trackers"
	flagVersion       bool
)

func cmdInit() {
	// CLI Flags
	pflag.CountVarP(&flagLogLevel, "verbose", "v", "Verbose level")
	pflag.StringVar(&flagConfigFolder, "config-dir", flagConfigFolder, "Config folder")
	pflag.StringVarP(&flagConfigFile, "config", "c", flagConfigFile, "Config file")
	pflag.StringVarP(&flagLogFile, "log", "l", flagLogFile, "Log file")
	pflag.StringVarP(&flagDbFile, "db", "d", flagDbFile, "Database file")
	pflag.StringVarP(&flagTrackerFolder, "trackers", "t", flagTrackerFolder, "Trackers folder")
	pflag.BoolVarP(&flagVersion, "version", "V", flagVersion, "Show version")

	// Parse CLI Flags
	pflag.Parse()

	// Add config folder if file not changed
	if !pflag.CommandLine.Changed("config") {
		flagConfigFile = filepath.Join(flagConfigFolder, flagConfigFile)
	}
	if !pflag.CommandLine.Changed("log") {
		flagLogFile = filepath.Join(flagConfigFolder, flagLogFile)
	}
	if !pflag.CommandLine.Changed("db") {
		flagDbFile = filepath.Join(flagConfigFolder, flagDbFile)
	}
	if !pflag.CommandLine.Changed("trackers") {
		flagTrackerFolder = filepath.Join(flagConfigFolder, flagTrackerFolder)
	}

	// Bind flags to Viper config, `config.RuntimeConfig`
	if err := config.RuntimeViper.BindPFlags(pflag.CommandLine); err != nil {
		log.Fatal("Failed to bind cmd flags to config")
	}

	// Bind Env vars
	config.RuntimeViper.AutomaticEnv()
}
