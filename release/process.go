package release

import (
	"github.com/docker/go-units"
	"gitlab.com/cloudb0x/trackarr/cache"
	"net/url"
	"strings"

	"gitlab.com/cloudb0x/trackarr/config"
	"gitlab.com/cloudb0x/trackarr/utils/web"

	"github.com/pkg/errors"
)

/* Private */

func (r *Release) getProxiedTorrentURL(cookie *string, pvr string) (string, error) {
	// parse torrent api url
	u, err := url.Parse(web.JoinURL(config.Config.Server.PublicURL, "/api/torrent"))
	if err != nil {
		return "", errors.Wrap(err, "failed parsing public torrent api endpoint")
	}

	// add query params
	q := u.Query()
	q.Set("apikey", config.Config.Server.ApiKey)
	q.Set("url", r.Info.TorrentURL)
	q.Set("pvr", pvr)
	if cookie != nil && *cookie != "" {
		q.Set("cookie", *cookie)
	}
	u.RawQuery = q.Encode()

	r.Log.Tracef("Proxied Torrent URL: %s", u.String())
	return u.String(), nil
}

/* Public */

func (r *Release) Process() {
	var err error
	usedBencode := false
	usedApi := false
	addedToCache := false

	r.Log.Tracef("Pre-processing: %s", r.Info.TorrentName)

	// replace https torrent urls with http (if ForceHTTP set)
	if r.Tracker.Config.ForceHTTP {
		r.Info.TorrentURL = strings.Replace(r.Info.TorrentURL, "https:", "http:", 1)
	}

	// api lookup if available for this tracker
	usedApi, err = r.apiLookup()
	if err != nil {
		// abort processing this release as an error occurred and bencode was disabled
		return
	}

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

	// bencode lookup if enabled and api was not used successfully before
	if !usedApi {
		usedBencode, addedToCache, err = r.bencodeLookup()
		if err != nil {
			r.Log.WithError(err).Error("Failed performing bencode lookup")
			return
		}
	}

	// was a size available?
	if r.Info.SizeBytes == 0 {
		r.Log.Warnf("Failed determining release size for %q, using 0 byte size...", r.Info.TorrentName)
	}

	r.Log.Debugf("Processing release: %s", r.Info.TorrentName)

	// iterate pvr's
	for _, pvr := range config.Pvr {
		// check ignore expressions
		ignore, err := pvr.ShouldIgnore(r.Info, r.Log)
		if err != nil {
			r.Log.WithError(err).Warnf("Failed checking release against ignore expressions for pvr: %q", pvr.Config.Name)
			continue
		}
		if ignore {
			r.Log.Debugf("Ignoring release for pvr: %q", pvr.Config.Name)
			continue
		}

		// check accept expressions
		accept, err := pvr.ShouldAccept(r.Info, r.Log)
		if err != nil {
			r.Log.WithError(err).Warnf("Failed checking release against accept expressions for pvr: %q", pvr.Config.Name)
			continue
		}
		if !accept {
			r.Log.Debugf("Release not accepted for pvr: %q", pvr.Config.Name)
			continue
		}

		// check delay expressions
		delay, err := pvr.ShouldDelay(r.Info, r.Log)
		if err != nil {
			r.Log.WithError(err).Warnf("Failed checking release against delay expressions for pvr: %q", pvr.Config.Name)
			continue
		}

		// add item to cache if not added already
		if !addedToCache {
			// store item in cache to be used by second-sweep
			go cache.AddItem(r.Info.TorrentURL, &cache.CacheItem{
				Name:    r.Info.TorrentName,
				Data:    nil,
				Release: r.Info,
			})

			addedToCache = true
		}

		// was bencode used / tracker requires a cookie / bencode should be used later on (to evaluate against torrent file data)
		torrentUrl := r.Info.TorrentURL
		cookie, hasCookie := r.Tracker.Config.Settings["cookie"]
		if usedBencode || hasCookie || (!usedBencode && pvr.HasFileExpressions) {
			// we will proxy this torrent via the /api/torrent endpoint
			proxiedURL, err := r.getProxiedTorrentURL(&cookie, pvr.Config.Name)
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
				// set the torrent url to the proxied url
				torrentUrl = proxiedURL
			}
		}

		// push release to pvr
		go func(p *config.PvrConfig, d *int64, url *string) {
			r.Push(p, d, url)
		}(pvr.Config, delay, &torrentUrl)
	}
}
