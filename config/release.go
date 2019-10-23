package config

type ReleaseInfo struct {
	TrackerName string
	ReleaseTime string
	TorrentName string
	TorrentURL  string
	SizeString  string
	SizeBytes   int64
	Category    string
	Encoder     string
	Resolution  string
	Container   string
	Origin      string
	Tags        string
	Files       []string
}
