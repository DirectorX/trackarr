package config

type TrackerConfiguration struct {
	Enabled bool
	Verbose bool
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
}
