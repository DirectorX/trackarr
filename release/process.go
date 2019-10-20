package release

import (
	"github.com/docker/go-units"
	"github.com/imroc/req"
	"github.com/l3uddz/trackarr/config"
	"github.com/l3uddz/trackarr/utils/torrent"
	"github.com/l3uddz/trackarr/utils/web"
	"github.com/pkg/errors"
	"net/url"
)

/* Const */
const TorrentFileTimeout = 15

/* Private */

func (r TrackerRelease) getProxiedTorrentURL(cookie *string) (string, error) {
	// parse torrent api url
	u, err := url.Parse(web.JoinURL(config.Config.Server.PublicURL, "/api/torrent"))
	if err != nil {
		return "", errors.Wrap(err, "failed parsing public torrent api endpoint")
	}

	// add query params
	q := u.Query()
	q.Set("apikey", config.Config.Server.ApiKey)
	q.Set("url", r.TorrentURL)
	if cookie != nil && *cookie != "" {
		q.Set("cookie", *cookie)
	}
	u.RawQuery = q.Encode()

	r.Log.Tracef("Proxied Torrent URL: %s", u.String())
	return u.String(), nil
}

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

	// was bencode used, or does this tracker have a cookie set?
	cookie, hasCookie := r.Cfg.Config["cookie"]
	if bencodeUsed || hasCookie {
		// we will proxy this torrent via the /api/torrent endpoint
		proxiedURL, err := r.getProxiedTorrentURL(&cookie)
		if err != nil {
			// we must abort this release as pvr's cant grab from these trackers without a cookie
			if hasCookie {
				r.Log.WithError(err).Errorf("Failed building proxied torrent url for: %q, aborting...",
					r.TorrentURL)
				return
			}

			// as this tracker does not require a cookie, we can continue with the original torrent url
			r.Log.WithError(err).Warnf("Failed building proxied torrent url for: %q", r.TorrentURL)
		} else {
			// set the torrent url to the proxied one
			r.TorrentURL = proxiedURL
		}
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
