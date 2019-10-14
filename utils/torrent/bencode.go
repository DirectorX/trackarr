package torrent

import (
	"github.com/l3uddz/trackarr/logger"
	"github.com/l3uddz/trackarr/utils/web"
	"github.com/pkg/errors"
	"github.com/zeebo/bencode"
)

var (
	log = logger.GetLogger("torrent")
)

/* Public */

func GetTorrentDetails(torrentUrl string, timeout int) (*TorrentFile, error) {
	// retrieve torrent file
	torrentBytes, err := web.GetBodyBytes(web.GET, torrentUrl, timeout)
	if err != nil {
		log.WithError(err).Errorf("Failed retrieving torrent bytes from: %s", torrentUrl)
		return nil, errors.Wrapf(err, "failed retrieving torrent bytes from: %s", torrentUrl)
	}

	// decode torrent data
	tf := &TorrentFile{}
	err = bencode.DecodeBytes(torrentBytes, tf)
	if err != nil {
		log.WithError(err).Errorf("Failed decoding torrent bytes from: %s", torrentUrl)
		return nil, errors.Wrapf(err, "failed decoding torrent bytes from: %s", torrentUrl)
	}

	return tf, nil
}