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
	"golang.org/x/time/rate"
	"net/url"
	"time"
)

/* Const */
const (
	btnApiUrl               = "https://api.broadcasthe.net/"
	btnTimeout              = 60
	btnApiRateLimitDuration = time.Hour
	btnApiRateLimit         = 150
)

/* Struct */
type Btn struct {
	log         *logrus.Entry
	tracker     *config.TrackerInstance
	rl          *rate.Limiter
	postRequest btnRequest
	headers     req.Header
}

type (
	btnRequest struct {
		Id      int32     `json:"id"`
		JsonRPC string    `json:"jsonrpc"`
		Method  string    `json:"method"`
		Params  [2]string `json:"params"`
	}
)

/* Private */

func newBtn(tracker *config.TrackerInstance) (Interface, error) {
	log := log.WithField("api", tracker.Name)

	// validate required tracker settings available
	apiKey, err := maps.GetStringMapValue(tracker.Config.Settings, "api_key", false)
	if err != nil {
		return nil, errors.WithMessage(err, "api_key setting missing")
	}

	// return api instance
	return &Btn{
		log:     log,
		tracker: tracker,
		headers: req.Header{
			"Content-Type": "application/json",
		},
		postRequest: btnRequest{
			Id:      1,
			JsonRPC: "2.0",
			Method:  "getTorrentById",
			Params:  [2]string{apiKey},
		},
		rl: web.GetRateLimiter(tracker.Name, btnApiRateLimit, btnApiRateLimitDuration),
	}, nil
}

/* Interface */

func (t *Btn) GetReleaseInfo(torrent *config.ReleaseInfo) (*TorrentInfo, error) {
	//extract torrentId
	torrentUrl, err := url.Parse(torrent.TorrentURL)
	if err != nil {
		return nil, fmt.Errorf("malformed URL: %s", torrent.TorrentURL)
	}

	// validate torrentId
	torrentId := torrentUrl.Query().Get("id")
	if torrentId == "" {
		return nil, fmt.Errorf("missing mandatory torrentId: %#v", torrent)
	}
	t.log.Tracef("Extracted TorrentID: %s", torrentId)

	// prepare request
	t.postRequest.Params[1] = torrentId
	t.log.Tracef("Request Body as JSON : %#v", req.BodyJSON(t.postRequest))

	// send request
	btnReleaseAsBytes, err := web.GetBodyBytes(web.POST, btnApiUrl, btnTimeout, req.BodyJSON(t.postRequest),
		&web.Retry{
			MaxAttempts:         6,
			ExpectedContentType: "application/json",
			Backoff: backoff.Backoff{
				Jitter: true,
				Min:    2 * time.Second,
				Max:    6 * time.Second,
			}}, t.headers, t.rl)
	if err != nil {
		return nil, errors.Wrapf(err, "failed retrieving torrent bytes from: %s", torrentId)
	}
	t.log.Tracef("Raw API Response : %s", string(btnReleaseAsBytes))

	// parse response
	var btnInfo struct {
		Id     int32
		Result struct {
			TorrentID   string // "1289304",
			Size        string // "1942229405",
			ReleaseName string // "Shahs.of.Sunset.S08E09.iNTERNAL.1080p.WEB.h264-TRUMP",
		}
	}

	if err := json.Unmarshal(btnReleaseAsBytes, &btnInfo); err != nil {
		t.log.WithError(err).Errorf("Failed unmarshalling response: %#v", btnReleaseAsBytes)
		return nil, errors.Wrap(err, "failed unmarshalling response")
	}
	t.log.Tracef("GetReleaseInfo Response: %+v", btnInfo)

	return &TorrentInfo{
		Name: btnInfo.Result.ReleaseName,
		Size: btnInfo.Result.Size,
	}, nil
}
