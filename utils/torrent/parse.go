package torrent

import (
	"github.com/jpillora/backoff"
	"path/filepath"
	"time"

	"github.com/l3uddz/trackarr/cache"
	"github.com/l3uddz/trackarr/utils/web"

	bencode "github.com/IncSW/go-bencode"
	"github.com/imroc/req"
	"github.com/pkg/errors"
)

var ()

/* Public */

func GetTorrentDetails(torrentUrl string, timeout int, headers req.Header) (*Data, error) {
	// retrieve torrent file
	torrentBytes, err := web.GetBodyBytes(web.GET, torrentUrl, timeout, &web.Retry{
		MaxAttempts:         5,
		ExpectedContentType: "torrent",
		Backoff: backoff.Backoff{
			Jitter: true,
			Min:    500 * time.Millisecond,
			Max:    3 * time.Second,
		}}, headers)
	if err != nil {
		return nil, errors.Wrapf(err, "failed retrieving torrent bytes from: %s", torrentUrl)
	}

	tf, err := TorrentDecode(torrentBytes)
	if err != nil {
		return nil, errors.Wrapf(err, "failed decoding torrent file: %s", torrentUrl)
	}

	// Files path and length, multipart vs single file torrent
	var files []string
	if tf.Info.Length > 0 {
		// there is only a single file
		files = append(files, tf.Info.Name)
	} else {
		// there are multiple files
		// add files to files slice and increase torrent size
		for _, f := range tf.Info.Files {
			files = append(files, f.Path)
			tf.Info.Length += f.Length
		}
	}

	// add torrent to cache
	go cache.AddItem(torrentUrl, &cache.CacheItem{
		Name: tf.Info.Name,
		Data: torrentBytes,
	})

	return &Data{
		Name:  tf.Info.Name,
		Size:  tf.Info.Length,
		Files: files,
	}, nil
}

//Decode byte-array of torrent file into TorrentMeta struct
func TorrentDecode(b []byte) (*TorrentMeta, error) {
	obj, err := bencode.Unmarshal(b)
	if err != nil {
		return nil, err
	}

	// Check if root object can be converted
	r, ok := obj.(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid torrent file metadata")
	}

	// Main torrent struct
	torrent := &TorrentMeta{}

	// Root
	// Announce
	if belem, ok := r["announce"].([]byte); ok {
		torrent.AnnounceList = append(torrent.AnnounceList, string(belem))
	}
	if belem, ok := r["announce-list"].([][]byte); ok {
		for _, elem := range belem {
			value := string(elem)
			if value != torrent.AnnounceList[0] { //Prevent duplicated with Announce_1
				torrent.AnnounceList = append(torrent.AnnounceList, value)
			}
		}
	}

	// Creation date
	if belem, ok := r["creation date"].(int64); ok {
		torrent.CreationDate = time.Unix(belem/1000, 0)
	}

	// Encoding
	if belem, ok := r["encoding"].([]byte); ok {
		torrent.Encoding = string(belem)
	}

	// Comment
	if belem, ok := r["comment"].([]byte); ok {
		torrent.Comment = string(belem)
	}

	// Created by
	if belem, ok := r["created by"].([]byte); ok {
		torrent.CreatedBy = string(belem)
	}

	// Info
	if info, ok := r["info"].(map[string]interface{}); ok {
		tInfo := &TorrentInfo{}

		// Name of root file or folder
		if belem, ok := info["name"].([]byte); ok {
			tInfo.Name = string(belem)
		}

		// Size
		if belem, ok := info["length"].(int64); ok {
			tInfo.Length = belem
		}

		// Size per piece
		if belem, ok := info["piece length"].(int64); ok {
			tInfo.PieceLength = belem
		}

		// Piece's SHA-1 hash
		// if belem, ok := info["pieces"].([]byte); ok {
		// 	tInfo.Pieces = hex.EncodeToString(belem)
		// }

		// Files
		if belem, ok := info["files"].([]interface{}); ok {
			tFiles := make([]*TorrentFile, 0, len(belem))

			for i := 0; i < len(belem); i++ {
				if file, ok := belem[i].(map[string]interface{}); ok {
					tfile := &TorrentFile{}
					if length, ok := file["length"].(int64); ok {
						tfile.Length = length
					}
					if path, ok := file["path"].([]interface{}); ok {
						sPath := make([]string, len(path))
						for x := 0; x < len(path); x++ {
							sPath[x] = string(path[x].([]byte))
						}
						tfile.Path = filepath.Join(sPath...)

						tFiles = append(tFiles, tfile)
					}
				}
			}

			tInfo.Files = tFiles
		}

		torrent.Info = tInfo
	}

	return torrent, nil
}
