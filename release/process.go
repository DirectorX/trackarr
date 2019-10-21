package release

import (
	"net/url"

	"github.com/docker/go-units"
	"github.com/imroc/req"
	"github.com/l3uddz/trackarr/config"
	"github.com/l3uddz/trackarr/utils/torrent"
	"github.com/l3uddz/trackarr/utils/web"
	"github.com/pkg/errors"
)

/* Const */
const TorrentFileTimeout = 15

/* Private */

func (r *Release) getProxiedTorrentURL(cookie *string) (string, error) {
	// parse torrent api url
	u, err := url.Parse(web.JoinURL(config.Config.Server.PublicURL, "/api/torrent"))
	if err != nil {
		return "", errors.Wrap(err, "failed parsing public torrent api endpoint")
	}

	// add query params
	q := u.Query()
	q.Set("apikey", config.Config.Server.ApiKey)
	q.Set("url", r.Info.TorrentURL)
	if cookie != nil && *cookie != "" {
		q.Set("cookie", *cookie)
	}
	u.RawQuery = q.Encode()

	r.Log.Tracef("Proxied Torrent URL: %s", u.String())
	return u.String(), nil
}

/* Public */

func (r *Release) Process() {
	bencodeUsed := false

	r.Log.Tracef("Pre-processing: %s", r.Info.TorrentName)

	// convert parsed release size string to bytes (required by pvr)
	if r.Info.SizeString != "" {
		// we have a size string, lets attempt to convert it to bytes
		if byteSize, err := units.FromHumanSize(r.Info.SizeString); err != nil {
			r.Log.WithError(err).Warnf("Failed converting parsed release size %q to bytes", r.Info.SizeString)
			r.Info.SizeBytes = 0
		} else {
			r.Info.SizeBytes = byteSize
			r.Log.Tracef("Converted parsed release size %q to %d bytes", r.Info.SizeString, r.Info.SizeBytes)
		}
	}

	// bencode torrent name and size? (we must enforce this functionality when a bytes size was not determined)
	if r.Tracker.Config.Bencode || r.Info.SizeBytes == 0 {
		// retrieve cookie if set for this tracker
		headers := req.Header{}
		if cookie, ok := r.Tracker.Config.Settings["cookie"]; ok {
			headers["Cookie"] = cookie
		}

		torrentData, err := torrent.GetTorrentDetails(r.Info.TorrentURL, TorrentFileTimeout, headers)
		if err != nil {
			// abort release as we are unable to retrieve the information we need
			return
		}

		// set release information from decoded torrent data
		r.Info.TorrentName = torrentData.Info.Name
		r.Info.SizeBytes = torrentData.Info.Size
		bencodeUsed = true
	}

	r.Log.Debugf("Processing release: %s", r.Info.TorrentName)

	// was bencode used, or does this tracker have a cookie set?
	cookie, hasCookie := r.Tracker.Config.Settings["cookie"]
	if bencodeUsed || hasCookie {
		// we will proxy this torrent via the /api/torrent endpoint
		proxiedURL, err := r.getProxiedTorrentURL(&cookie)
		if err != nil {
			// we must abort this release as pvr's cant grab from these trackers without a cookie
			if hasCookie {
				r.Log.WithError(err).Errorf("Failed building proxied torrent url for: %q, aborting...",
					r.Info.TorrentURL)
				return
			}

			// as this tracker does not require a cookie, we can continue with the original torrent url
			r.Log.WithError(err).Warnf("Failed building proxied torrent url for: %q", r.Info.TorrentURL)
		} else {
			// set the torrent url to the proxied one
			r.Info.TorrentURL = proxiedURL
		}
	}

	// iterate pvr's
	for _, pvr := range config.Pvr {
		// check ignore expressions
		ignore, err := r.shouldIgnore(pvr)
		if err != nil {
			r.Log.WithError(err).Warnf("Failed checking release against ignore expressions for pvr: %q", pvr.Config.Name)
			continue
		}
		if ignore {
			r.Log.Debugf("Ignoring release for pvr: %q", pvr.Config.Name)
			continue
		}

		// check accept expressions
		accept, err := r.shouldAccept(pvr)
		if err != nil {
			r.Log.WithError(err).Warnf("Failed checking release against accept expressions for pvr: %q", pvr.Config.Name)
			continue
		}
		if !accept {
			r.Log.Debugf("Release not accepted for pvr: %q", pvr.Config.Name)
			continue
		}

		// check delay expressions
		delay, err := r.shouldDelay(pvr)
		if err != nil {
			r.Log.WithError(err).Warnf("Failed checking release against delay expressions for pvr: %q", pvr.Config.Name)
			continue
		}

		// push release to pvr
		go func(p *config.PvrConfig, d *int64) {
			r.Push(p, d)
		}(pvr.Config, delay)
	}
}
