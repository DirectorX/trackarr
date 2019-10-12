package config

type TrackerConfiguration struct {
	IRC TrackerIrcConfiguration
}

type TrackerIrcConfiguration struct {
	Nickname string
	Host     *string
	Port     int
	TLS      bool
}
