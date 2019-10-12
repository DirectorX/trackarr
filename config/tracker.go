package config

type TrackerConfiguration struct {
	Nickname string
	IRC      TrackerIrcConfiguration
}

type TrackerIrcConfiguration struct {
	Host *string
	Port int
	TLS  bool
}
