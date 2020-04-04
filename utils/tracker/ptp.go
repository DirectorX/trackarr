package tracker

import (
	"fmt"
	"github.com/imroc/req"
	"github.com/jpillora/backoff"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gitlab.com/cloudb0x/trackarr/config"
	"gitlab.com/cloudb0x/trackarr/utils/maps"
	"gitlab.com/cloudb0x/trackarr/utils/web"
	"go.uber.org/ratelimit"
	"time"
)

/* Const */
const (
	ptpTorrentUrl   = "https://passthepopcorn.me/torrents.php"
	ptpTimeout      = 30
	ptpApiRateLimit = 1
)

/* Struct */
type Ptp struct {
	log     *logrus.Entry
	tracker *config.TrackerInstance
	headers req.Header
	rl      *ratelimit.Limiter
}

/* Private */

func newPtp(tracker *config.TrackerInstance) (Interface, error) {
	log := log.WithField("api", tracker.Name)

	// validate required tracker settings available
	apiUser, err := maps.GetStringMapValue(tracker.Config.Settings, "api_user", false)
	if err != nil {
		return nil, errors.WithMessage(err, "api_user setting missing")
	}

	apiKey, err := maps.GetStringMapValue(tracker.Config.Settings, "api_key", false)
	if err != nil {
		return nil, errors.WithMessage(err, "api_key setting missing")
	}

	// return api instance
	return &Ptp{
		log:     log,
		tracker: tracker,
		headers: req.Header{
			"ApiUser": apiUser,
			"ApiKey":  apiKey,
		},
		rl: web.GetRateLimiter(tracker.Name, ptpApiRateLimit),
	}, nil
}

/* Interface */

func (t *Ptp) GetReleaseInfo(torrent *config.ReleaseInfo) (*TorrentInfo, error) {
	// validate torrent has required information
	if torrent.TorrentId == "" {
		return nil, fmt.Errorf("missing mandatory torrentId: %#v", torrent)
	}

	// send request
	ptpReleaseAsBytes, err := web.GetBodyBytes(web.GET, ptpTorrentUrl, ptpTimeout, req.QueryParam{
		"torrentid": torrent.TorrentId,
	}, &web.Retry{
		MaxAttempts:         5,
		ExpectedContentType: "application/json",
		Backoff: backoff.Backoff{
			Jitter: true,
			Min:    2 * time.Second,
			Max:    6 * time.Second,
		}}, t.headers, t.rl)
	if err != nil {
		return nil, errors.Wrapf(err, "failed retrieving torrent info bytes for: %s", torrent.TorrentId)
	}

	// parse response
	var ptpInfo struct {
		Torrents []struct {
			Id          string // "1368",
			Size        string // "1468434432",
			ReleaseName string // "The Godfather Part 3",
		}
	}

	if err := json.Unmarshal(ptpReleaseAsBytes, &ptpInfo); err != nil {
		t.log.WithError(err).Errorf("Failed unmarshalling response: %#v", string(ptpReleaseAsBytes))
		return nil, errors.Wrap(err, "failed unmarshalling response")
	}

	t.log.Tracef("GetReleaseInfo Response: %+v", ptpInfo)

	// find torrent in parsed response
	for _, v := range ptpInfo.Torrents {
		if v.Id == torrent.TorrentId {
			return &TorrentInfo{
				Name:     v.ReleaseName,
				Category: "",
				Size:     v.Size,
			}, nil
		}
	}

	return nil, fmt.Errorf("no release found with id: %s", torrent.TorrentId)
}
