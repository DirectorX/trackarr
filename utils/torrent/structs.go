package torrent

type TorrentFile struct {
	Announce string          `bencode:"announce"`
	Comment  string          `bencode:"comment"`
	Info     TorrentFileInfo `bencode:"info"`
}

type TorrentFileInfo struct {
	Size int64  `bencode:"length"`
	Name string `bencode:"name"`
}
