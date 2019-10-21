package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/l3uddz/trackarr/logger"
	stringutils "github.com/l3uddz/trackarr/utils/strings"

	"github.com/lithammer/shortuuid"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// BuildVars build details
type BuildVars struct {
	// Version
	Version string
	// Timestamp
	Timestamp string
	// Git commit
	GitCommit string
}

type Configuration struct {
	Server   ServerConfig
	Pvr      []PvrConfig
	Trackers map[string]TrackerConfig
}

/* Vars */

var (
	// Config exports the config object
	Config *Configuration
	// Runtime flags/config
	Runtime      *RuntimeConfig
	RuntimeViper = viper.New()
	// ldflags (set by makefile or goreleaser)
	Build *BuildVars
	// State
	Pvr = make(map[string]*PvrInstance)

	// Internal
	log          = logger.GetLogger("cfg")
	newOptionLen = 0
)

/* Public */

func (cfg Configuration) ToJsonString() (string, error) {
	c := viper.AllSettings()
	bs, err := json.MarshalIndent(c, "", "  ")
	return string(bs), err
}

func Init(build *BuildVars) error {
	// Set build vars
	Build = build

	// Unmarshal runtime/cmd config
	if err := RuntimeViper.Unmarshal(&Runtime); err != nil {
		return errors.Wrap(err, "failed to unmarshal runtime config")
	}

	// Info
	log.Infof("Using %s = %q", stringutils.StringLeftJust("CONFIG", " ", 10), Runtime.Config)

	/* Initialize Configuration */
	viper.SetConfigType("yaml")
	viper.SetConfigFile(Runtime.Config)

	// read matching env vars
	viper.AutomaticEnv()

	// Load config
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok || os.IsNotExist(err) {
			// set the default config to be written
			if err := setConfigDefaults(false); err != nil {
				log.WithError(err).Error("Failed to add config defaults")
				return errors.Wrap(err, "failed adding config defaults")
			}

			// write default config
			if err := viper.WriteConfig(); err != nil {
				log.WithError(err).Fatalf("Failed dumping default configuration to %q", Runtime.Config)
			}

			log.Infof("Dumped default configuration to %q. Please edit before running again!", viper.ConfigFileUsed())
			log.Logger.Exit(0)
		}

		log.WithError(err).Error("Configuration read error")
		return errors.Wrap(err, "failed reading config")
	}

	// Set defaults (checking whether new options were added)
	if err := setConfigDefaults(true); err != nil {
		log.WithError(err).Error("Failed to add new config defaults")
		return errors.Wrap(err, "failed adding new config defaults")
	}

	// Unmarshal into Config struct
	if err := viper.Unmarshal(&Config); err != nil {
		log.WithError(err).Error("Configuration decode error")
		return errors.Wrap(err, "failed decoding config")
	}

	return nil
}

func PrintVersion() {
	log.Infof("Trackarr version %s (%s@%s)", Build.Version, Build.GitCommit, Build.Timestamp)
}

/* Private */

func setConfigDefault(key string, value interface{}, check bool) int {
	if check {
		if viper.IsSet(key) {
			return 0
		}

		// determine padding to use for new key
		if keyLen := len(key); (keyLen + 2) > newOptionLen {
			newOptionLen = keyLen + 2
		}

		log.Warnf("New config option: %s = %v", stringutils.StringLeftJust(fmt.Sprintf("%q", key), " ", newOptionLen), value)
	}

	viper.SetDefault(key, value)

	return 1
}

func setConfigDefaults(check bool) error {
	added := 0

	// server settings
	added += setConfigDefault("server.host", "0.0.0.0", check)
	added += setConfigDefault("server.port", 7337, check)
	added += setConfigDefault("server.apikey", shortuuid.New(), check)
	added += setConfigDefault("server.publicurl", "http://trackarr.domain.com", check)

	// pvr settings
	added += setConfigDefault("pvr", []PvrConfig{
		{
			Name:    "sonarr_main",
			Enabled: false,
			URL:     "https://sonarr.domain.com",
			ApiKey:  "YOUR_API_KEY",
			Ignores: []string{
				`Category contains "Disk"`,
				`Category contains "DVD-R"`,
				`TorrentName contains "Disc"`,
			},
			Accepts: []string{
				`TrackerName == "IPT" && Category startsWith "TV/"`,
			},
			Delays: []string{
				`Category startsWith "TV/" ? 1 : 0`,
			},
		},
	}, check)

	// tracker settings
	added += setConfigDefault("trackers.iptorrents.enabled", true, check)
	added += setConfigDefault("trackers.iptorrents.bencode", false, check)
	added += setConfigDefault("trackers.iptorrents.settings.passkey", "", check)
	added += setConfigDefault("trackers.iptorrents.irc.nickname", "therugmuncher_autodl", check)
	added += setConfigDefault("trackers.iptorrents.irc.channels", []string{"#ipt.announce"}, check)
	added += setConfigDefault("trackers.iptorrents.irc.commands", []string{
		"PRIVMSG NickServ IDENTIFY YOUR_PASS",
	}, check)
	added += setConfigDefault("trackers.iptorrents.irc.verbose", false, check)

	// were new settings added?
	if check && added > 0 {
		if err := viper.WriteConfig(); err != nil {
			log.WithError(err).Error("Failed saving configuration with new options...")
			return errors.Wrap(err, "failed saving updated configuration")
		}

		log.Info("Configuration was saved with new options!")
		log.Logger.Exit(0)
	}

	return nil
}
