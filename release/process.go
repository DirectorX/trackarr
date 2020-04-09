package release

import (
	"gitlab.com/cloudb0x/trackarr/cache"
	stringutils "gitlab.com/cloudb0x/trackarr/utils/strings"
	"gitlab.com/cloudb0x/trackarr/utils/tracker"
	"net/url"
	"strings"

	"gitlab.com/cloudb0x/trackarr/config"
	"gitlab.com/cloudb0x/trackarr/utils/torrent"
	"gitlab.com/cloudb0x/trackarr/utils/web"

	"github.com/docker/go-units"
	"github.com/imroc/req"
	"github.com/pkg/errors"
)

/* Const */
const TorrentFileTimeout = 30

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
	bencodeUsed := false
	apiUsed := false
	addedToCache := false

	r.Log.Tracef("Pre-processing: %s", r.Info.TorrentName)

	// replace https torrent urls with http (if ForceHTTP set)
	if r.Tracker.Config.ForceHTTP {
		r.Info.TorrentURL = strings.Replace(r.Info.TorrentURL, "https:", "http:", 1)
	}

	// retrieve api for this tracker (if set)
	trackerApi, _ := tracker.GetApi(r.Tracker)
	if trackerApi != nil {
		// lookup torrent info via the associated api interface
		torrentInfo, err := trackerApi.GetReleaseInfo(r.Info)
		if err != nil {
			// api lookup for torrent failed due to an error
			r.Log.WithError(err).Errorf("Failed looking up missing info via api for torrent: %q", r.Info.TorrentName)
			if !r.Tracker.Config.Bencode.Name && !r.Tracker.Config.Bencode.Size {
				// bencode is disabled so no fallback
				r.Log.Warnf("Aborting push of release as bencode disabled for torrent: %q", r.Info.TorrentName)
				return
			}
			// bencode is enabled, so continue legacy behaviour
		} else if torrentInfo == nil {
			// api lookup failed for some known reason (max login attempts etc..) - fallback to bencode if enabled
			if !r.Tracker.Config.Bencode.Name && !r.Tracker.Config.Bencode.Size {
				// bencode is disabled so no fallback
				r.Log.Warnf("Aborting push of release as bencode disabled for torrent: %q", r.Info.TorrentName)
				return
			}
			// bencode is enabled, so continue legacy behaviour
		} else {
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
				return
			}

			apiUsed = true
		}
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

	// bencode torrent name and size?
	if (r.Tracker.Config.Bencode.Name || r.Tracker.Config.Bencode.Size) && !apiUsed {
		// retrieve cookie if set for this tracker
		headers := req.Header{}
		if cookie, ok := r.Tracker.Config.Settings["cookie"]; ok {
			headers["Cookie"] = cookie
		}

		torrentData, err := torrent.GetTorrentDetails(r.Info.TorrentURL, TorrentFileTimeout, headers)
		if err != nil {
			// abort release as we are unable to retrieve the information we need
			r.Log.WithError(err).Error("Failed decoding details from torrent file")
			return
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

		addedToCache = true
		bencodeUsed = true
	} else if r.Info.SizeBytes == 0 {
		r.Log.Warnf("Failed determining release size for %q as no size parsed from announcement and"+
			" bencode was disabled.", r.Info.TorrentName)
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
		if bencodeUsed || hasCookie || (!bencodeUsed && pvr.HasFileExpressions) {
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
