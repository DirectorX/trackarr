package release

import "github.com/l3uddz/trackarr/utils/torrent"

/* Const */
const TorrentFileTimeout = 15

/* Public */

func (r *TrackerRelease) Process() {
	// bencode torrent name and size?
	if r.Cfg.Bencode {
		torrentData, err := torrent.GetTorrentDetails(r.TorrentURL, TorrentFileTimeout)
		if err != nil {
			// abort release as we are unable to retrieve the information we need
			return
		}

		// set release information from decoded torrent data
		r.TorrentName = torrentData.Info.Name
		r.TorrentSizeBytes = &torrentData.Info.Size
	}

	r.Log.Debugf("Processing release: %s", r.TorrentName)

}
