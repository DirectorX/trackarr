package config

import "github.com/pkg/errors"

type TrackerConfiguration struct {
	Enabled bool
	Bencode bool
	Config  map[string]string
	IRC     TrackerIrcConfiguration
}

type TrackerIrcConfiguration struct {
	Nickname   string
	Channels   []string
	Announcers []string
	Commands   [][]string
	Host       *string
	Port       *string
	Verbose bool
}

/* Public */
func GetAnyConfiguredTracker(trackers *map[string]TrackerConfiguration) (*TrackerConfiguration, error) {
	for _, trackerConfig := range *trackers {
		return &trackerConfig, nil
	}

	return nil, errors.New("no tracker configuration found")
}
