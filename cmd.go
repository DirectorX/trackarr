package main

import (
	"path/filepath"

	"github.com/l3uddz/trackarr/config"
	"github.com/l3uddz/trackarr/utils/paths"

	"github.com/spf13/pflag"
)

// CLI flags with defaults
var (
	flagLogLevel    int
	flagConfigPath  = filepath.Join(paths.GetCurrentBinaryPath(), "config.yaml")
	flagLogPath     = filepath.Join(paths.GetCurrentBinaryPath(), "activity.log")
	flagDbPath      = filepath.Join(paths.GetCurrentBinaryPath(), "vault.db")
	flagTrackerPath = filepath.Join(paths.GetCurrentBinaryPath(), "trackers")
	flagVersion     bool
)

func cmdInit() {
	// CLI Flags
	pflag.CountVarP(&flagLogLevel, "verbose", "v", "Verbose level")
	pflag.StringVarP(&flagConfigPath, "config", "c", flagConfigPath, "Config path")
	pflag.StringVarP(&flagLogPath, "log", "l", flagLogPath, "Log path")
	pflag.StringVarP(&flagDbPath, "db", "d", flagDbPath, "Database path")
	pflag.StringVarP(&flagTrackerPath, "trackers", "t", flagTrackerPath, "Trackers path")
	pflag.BoolVarP(&flagVersion, "version", "V", flagVersion, "Show version")

	// Parse CLI Flags
	pflag.Parse()

	// Bind flags to Viper config, `config.RuntimeConfig`
	if err := config.RuntimeViper.BindPFlags(pflag.CommandLine); err != nil {
		log.Fatal("Failed to bind cmd flags to config")
	}

	// Bind Env vars
	config.RuntimeViper.AutomaticEnv()
}
