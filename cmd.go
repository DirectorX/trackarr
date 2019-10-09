package main

import (
	"path/filepath"

	"github.com/l3uddz/trackarr/utils/paths"
	"github.com/spf13/pflag"
)

// CLI flags with defaults
var (
	flagLogLevel    int
	flagConfigPath  = filepath.Join(paths.GetCurrentBinaryPath(), "config.json")
	flagLogPath     = filepath.Join(paths.GetCurrentBinaryPath(), "activity.log")
	flagDbPath      = filepath.Join(paths.GetCurrentBinaryPath(), "vault.db")
	flagTrackerPath = filepath.Join(paths.GetCurrentBinaryPath(), "trackers")
)

func cmdInit() {
	// CLI Flags
	pflag.CountVarP(&flagLogLevel, "verbose", "v", "Verbose level")
	pflag.StringVarP(&flagConfigPath, "config", "c", flagConfigPath, "Config path")
	pflag.StringVarP(&flagLogPath, "log", "l", flagLogPath, "Log path")
	pflag.StringVarP(&flagLogPath, "db", "d", flagLogPath, "Database path")
	pflag.StringVarP(&flagTrackerPath, "track", "t", flagTrackerPath, "Trackers path")

	// Parse CLI Flags
	pflag.Parse()
}
