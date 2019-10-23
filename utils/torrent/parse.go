package torrent

import (
	"github.com/imroc/req"
	"github.com/l3uddz/trackarr/cache"
	"github.com/l3uddz/trackarr/logger"
	"github.com/l3uddz/trackarr/utils/web"
	"github.com/pkg/errors"
	"github.com/zeebo/bencode"
)

var (
	log = logger.GetLogger("torrent")
)

/* Public */

// Credits: https://github.com/j-muller/go-torrent-parser
func GetTorrentDetails(torrentUrl string, timeout int, headers req.Header) (*Data, error) {
	// retrieve torrent file
	torrentBytes, err := web.GetBodyBytes(web.GET, torrentUrl, timeout, headers)
	if err != nil {
		log.WithError(err).Errorf("Failed retrieving torrent bytes from: %s", torrentUrl)
		return nil, errors.Wrapf(err, "failed retrieving torrent bytes from: %s", torrentUrl)
	}

	// decode torrent data
	tf := &Metadata{}
	err = bencode.DecodeBytes(torrentBytes, tf)
	if err != nil {
		log.WithError(err).Errorf("Failed decoding torrent bytes from: %s", torrentUrl)
		return nil, errors.Wrapf(err, "failed decoding torrent bytes from: %s", torrentUrl)
	}

	// decode files data
	files := make([]string, 0)
	// single file context
	if tf.Info.Size > 0 {
		files = append(files, tf.Info.Name)
	} else {
		// decode files metadata
		metadataFiles := make([]*FileMetadata, 0)
		err = bencode.DecodeBytes(tf.Info.Files, &metadataFiles)
		if err != nil {
			return nil, errors.Wrapf(err, "failed decoding files torrent bytes from: %s", torrentUrl)
		}

		// add file to files slice and increase torrent size
		for _, f := range metadataFiles {
			files = append(files, f.Path...)
			tf.Info.Size += f.Length
		}
	}

	// add torrent to cache
	cache.AddItem(torrentUrl, &cache.CacheItem{
		Name: tf.Info.Name,
		Data: torrentBytes,
	})

	return &Data{
		Name:  tf.Info.Name,
		Size:  tf.Info.Size,
		Files: files,
	}, nil
}
