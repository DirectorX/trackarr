package release

import (
	"github.com/pkg/errors"
	stringutils "gitlab.com/cloudb0x/trackarr/utils/strings"
	"gitlab.com/cloudb0x/trackarr/utils/tracker"
)

func (r *Release) apiLookup() (bool, error) {
	// retrieve tracker api
	trackerApi, _ := tracker.GetApi(r.Tracker)
	if trackerApi == nil {
		// api not implemented for this tracker
		return false, nil
	}

	// lookup torrent info via the associated api interface
	torrentInfo, err := trackerApi.GetReleaseInfo(r.Info)
	if err != nil {
		// api lookup for torrent failed due to an error
		r.Log.WithError(err).Errorf("Failed looking up missing info via api for torrent: %q", r.Info.TorrentName)
		if !r.Tracker.Config.Bencode.Name && !r.Tracker.Config.Bencode.Size {
			// bencode is disabled so no fallback
			r.Log.Warnf("Aborting push of release as bencode disabled for torrent: %q", r.Info.TorrentName)
			return false, errors.New("api lookup failed and bencode disabled")
		}

		// bencode is enabled, so continue legacy behaviour
		return false, nil
	}

	if torrentInfo == nil {
		// api lookup failed for some known reason (max login attempts etc..) - fallback to bencode if enabled
		if !r.Tracker.Config.Bencode.Name && !r.Tracker.Config.Bencode.Size {
			// bencode is disabled so no fallback
			r.Log.Warnf("Aborting push of release as bencode disabled for torrent: %q", r.Info.TorrentName)
			return false, errors.New("api lookup failed and bencode disabled")
		}

		// bencode is enabled, so continue legacy behaviour
		return false, nil
	}

	// api lookup was successful, process response
	r.Log.Debugf("Retrieved torrent info via api: %+v", torrentInfo)

	// set info from api lookup
	r.Info.TorrentName = stringutils.NewOrExisting(&torrentInfo.Name, r.Info.TorrentName)
	r.Info.Category = stringutils.NewOrExisting(&torrentInfo.Category, r.Info.Category)
	r.Info.SizeString = stringutils.NewOrExisting(&torrentInfo.Size, r.Info.SizeString)

	// validate required information was parsed
	if r.Info.SizeString == "" && (!r.Tracker.Config.Bencode.Name && !r.Tracker.Config.Bencode.Size) {
		// no size string was retrieved from api, however, bencode is disabled, so we cannot proceed
		r.Log.Warnf("Aborting push of release as api response was incomplete and bencode disabled"+
			" for torrent: %q", r.Info.TorrentName)
		return false, errors.New("api lookup response incomplete and bencode disabled")
	}

	return true, nil
}
