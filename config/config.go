package config

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"os"

	"github.com/l3uddz/trackarr/logger"
	stringutils "github.com/l3uddz/trackarr/utils/strings"
	"github.com/spf13/viper"
)

/* Config */

type Configuration struct {
	Server   ServerConfiguration
	Trackers map[string]TrackerIrcConfiguration
}

/* Vars */

var (
	// Config exports the config object
	Config       *Configuration
	log          = logger.GetLogger("cfg")
	newOptionLen = 0
)

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

	// tracker settings
	added += setConfigDefault("trackers", map[string]TrackerConfiguration{
		"IPTorrents": {
			Nickname: "thebigmuncho",
			IRC: TrackerIrcConfiguration{
				Port: 6667,
				TLS:  true,
			},
		},
	}, check)

	// were new settings added?
	if check && added > 0 {
		if err := viper.WriteConfig(); err != nil {
			log.WithError(err).Error("Failed saving configuration with new options...")
			return errors.Wrap(err, "failed saving updated configuration")
		}

		log.Info("Configuration was saved with new options!")
		os.Exit(0)
	}

	return nil
}

/* Public */

func (cfg Configuration) ToJsonString() (string, error) {
	c := viper.AllSettings()
	bs, err := json.MarshalIndent(c, "", "  ")
	return string(bs), err
}

func Init(configFilePath string) error {
	// Info
	log.Infof("Using %s = %q", stringutils.StringLeftJust("CONFIG", " ", 10), configFilePath)

	/* Initialize Configuration */
	viper.SetConfigFile(configFilePath)

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
			if err := viper.WriteConfigAs(configFilePath); err != nil {
				log.WithError(err).Errorf("Failed dumping default configuration to %q", configFilePath)
				os.Exit(1)
			} else {
				log.Infof("Dumped default configuration to %q. Please edit before running again!", viper.ConfigFileUsed())
				os.Exit(0)
			}
		}

		log.WithError(err).Error("Configuration read error")
		return errors.Wrap(err, "failed reading config")
	}

	if err := viper.Unmarshal(&Config); err != nil {
		log.WithError(err).Error("Configuration decode error")
		return errors.Wrap(err, "failed decoding config")
	}

	// Set defaults (checking whether new options were added)
	if err := setConfigDefaults(true); err != nil {
		log.WithError(err).Error("Failed to add new config defaults")
		return errors.Wrap(err, "failed adding new config defaults")
	}

	return nil
}
