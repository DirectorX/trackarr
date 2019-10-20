package release

import (
	"github.com/docker/go-units"
	"github.com/imroc/req"
	"github.com/l3uddz/trackarr/utils/torrent"
)

/* Const */
const TorrentFileTimeout = 15

/* Public */

func (r *TrackerRelease) Process() {
	bencodeUsed := false

	r.Log.Tracef("Pre-processing: %s", r.TorrentName)

	// convert parsed release size string to bytes (required by pvr)
	if r.SizeString != "" {
		// we have a size string, lets attempt to convert it to bytes
		if byteSize, err := units.FromHumanSize(r.SizeString); err != nil {
			r.Log.WithError(err).Warnf("Failed converting parsed release size %q to bytes", r.SizeString)
			r.SizeBytes = 0
		} else {
			r.SizeBytes = byteSize
			r.Log.Tracef("Converted parsed release size %q to %d bytes", r.SizeString, r.SizeBytes)
		}
	}

	// bencode torrent name and size? (we must enforce this functionality when a bytes size was not determined)
	if r.Cfg.Bencode || r.SizeBytes == 0 {
		// retrieve cookie if set for this tracker
		headers := req.Header{}
		if cookie, ok := r.Cfg.Config["cookie"]; ok {
			headers["Cookie"] = cookie
		}

		torrentData, err := torrent.GetTorrentDetails(r.TorrentURL, TorrentFileTimeout, headers)
		if err != nil {
			// abort release as we are unable to retrieve the information we need
			return
		}

		// set release information from decoded torrent data
		r.TorrentName = torrentData.Info.Name
		r.SizeBytes = torrentData.Info.Size
		bencodeUsed = true
	}

	r.Log.Debugf("Processing: %s", r.TorrentName)

	// TODO: was bencode used? if so we should alter the url so the cached torrent is grabbed (on approve)
	if bencodeUsed {
		// update torrenturl to use the /api/torrent endpoint
	}

	// iterate pvr's
	for pvr, expressions := range pvrExpressions {
		// check ignore expressions
		ignore, err := r.shouldIgnore(pvr, &expressions)
		if err != nil {
			r.Log.WithError(err).Warnf("Failed checking release against ignore expressions for pvr: %q", pvr.Name)
			continue
		}

		if ignore {
			r.Log.Debugf("Ignoring release for pvr: %q", pvr.Name)
			continue
		}

		// check accept expressions
		accept, err := r.shouldAccept(pvr, &expressions)
		if err != nil {
			r.Log.WithError(err).Warnf("Failed checking release against accept expressions for pvr: %q", pvr.Name)
			continue
		}

		if !accept {
			r.Log.Debugf("Release not accepted for pvr: %q", pvr.Name)
			continue
		}

		// push release to pvr
		pvrObj := pvr
		go r.Push(pvrObj)
	}
}
