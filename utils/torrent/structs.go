package torrent

import (
	"time"
)

//Contains all the info about the torrent file
type TorrentFile struct {
	Length int64
	Path   string
}

//Contains all the meta-info data from the original torrent file
type TorrentInfo struct {
	Name        string
	Length      int64
	PieceLength int64
	// Pieces      string
	Files []*TorrentFile
}

//Contains all the meta-file data from the original torrent file
type TorrentMeta struct {
	AnnounceList []string
	CreationDate time.Time
	Encoding     string
	Comment      string
	CreatedBy    string
	Info         *TorrentInfo
}

type Data struct {
	Name  string
	Size  int64
	Files []string
}
