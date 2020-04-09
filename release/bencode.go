package release

import (
	"github.com/imroc/req"
	"github.com/pkg/errors"
	"gitlab.com/cloudb0x/trackarr/cache"
	"gitlab.com/cloudb0x/trackarr/utils/torrent"
)

func (r *Release) bencodeLookup() (bool, bool, error) {
	// was bencode enabled for this tracker?
	if !r.Tracker.Config.Bencode.Name && !r.Tracker.Config.Bencode.Size {
		return false, false, nil
	}

	// retrieve cookie if set for this tracker
	headers := req.Header{}
	if cookie, ok := r.Tracker.Config.Settings["cookie"]; ok {
		headers["Cookie"] = cookie
	}

	torrentData, err := torrent.GetTorrentDetails(r.Info.TorrentURL, TorrentFileTimeout, headers)
	if err != nil {
		// abort release as we are unable to retrieve the information we need
		return false, false, errors.WithMessage(err, "failed decoding details from torrent file")
	}

	// store parsed torrent files in release
	r.Info.Files = torrentData.Files

	// set release information from decoded torrent data
	if r.Tracker.Config.Bencode.Name {
		// bencode name was set to true
		r.Info.TorrentName = torrentData.Name
	}

	if r.Tracker.Config.Bencode.Size || r.Info.SizeBytes == 0 {
		// bencode size was set to true (or we had no size from parsed release)
		r.Info.SizeBytes = torrentData.Size
	}

	// add torrent to cache
	go cache.AddItem(r.Info.TorrentURL, &cache.CacheItem{
		Name:    r.Info.TorrentName,
		Data:    torrentData.Bytes,
		Release: r.Info,
	})

	return true, true, nil
}
