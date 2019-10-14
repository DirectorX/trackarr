package config

type TrackerConfiguration struct {
	Enabled bool
	Verbose bool
	Config  map[string]string
	IRC     TrackerIrcConfiguration
}

type TrackerIrcConfiguration struct {
	Nickname   string
	Channels   []string
	Announcers []string
	Host       *string
	Port       *string
	Commands   [][]string
}
