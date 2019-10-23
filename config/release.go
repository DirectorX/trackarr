package config

import "github.com/l3uddz/trackarr/utils/torrent"

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
	Torrent     *torrent.Data
}
