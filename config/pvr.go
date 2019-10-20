package config

type PvrConfiguration struct {
	Name    string
	Enabled bool
	URL     string
	ApiKey  string
	Ignores []string
	Accepts []string
}
