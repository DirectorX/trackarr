package torrent

import "github.com/zeebo/bencode"

type Metadata struct {
	Announce string       `bencode:"announce"`
	Comment  string       `bencode:"comment"`
	Info     InfoMetadata `bencode:"info"`
}

type InfoMetadata struct {
	Name  string             `bencode:"name"`
	Size  int64              `bencode:"length"`
	Files bencode.RawMessage `bencode:"files"`
}

type FileMetadata struct {
	Path   []string `bencode:"path"`
	Length int64    `bencode:"length"`
}

type Data struct {
	Name  string
	Size  int64
	Files []string
}
