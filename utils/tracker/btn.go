package tracker

import (
	"encoding/json"
	"fmt"
	"github.com/imroc/req"
	"github.com/jpillora/backoff"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gitlab.com/cloudb0x/trackarr/config"
	"gitlab.com/cloudb0x/trackarr/utils/maps"
	"gitlab.com/cloudb0x/trackarr/utils/web"
	"net/url"
	"time"
)

/* Const */
const (
	btnApiUrl       = "https://api.broadcasthe.net/"
	btnTimeout      = 60
	btnApiRateLimit = 150
)

/* Var */
type Btn struct {
	log     *logrus.Entry
	tracker *config.TrackerInstance
}

type (
	btnRequest struct {
		Id      int32     `json:"id"`
		JsonRPC string    `json:"jsonrpc"`
		Method  string    `json:"method"`
		Params  [2]string `json:"params"`
	}
)

/* Interface */

func (t *Btn) GetReleaseInfo(torrent *config.ReleaseInfo) (*TorrentInfo, error) {
	t.log.Tracef("BTN Announce releaseInfo: %#v", torrent)
	//extract torrentId
	torrentUrl, err := url.Parse(torrent.TorrentURL)
	if err != nil {
		t.log.WithError(err).Error("Malformed UR")
		return nil, errors.Wrap(err, "Malformed URL")
	}
	torrentId := torrentUrl.Query().Get("id")
	if torrentId == "" {
		return nil, fmt.Errorf("missing mandatory torrentId: %#v", torrent)
	}
	t.log.Tracef("Extracted TorrentID: %s", torrentId)

	// prepare request
	apiKey, err := maps.GetStringMapValue(t.tracker.Config.Settings, "api_key", false)
	if err != nil {
		t.log.WithError(err).Error("api_key setting missing")
		return nil, errors.Wrap(err, "api_key setting missing")
	}

	headers := req.Header{
		"Content-Type": "application/json",
	}

	postRequest := btnRequest{
		Id:      1,
		JsonRPC: "2.0",
		Method:  "getTorrentById",
		Params:  [2]string{apiKey, torrentId},
	}

	t.log.Tracef("Request Body as JSON : %#v", req.BodyJSON(postRequest))
	// send request
	btnReleaseAsBytes, err := web.GetBodyBytes(web.POST, btnApiUrl, btnTimeout, req.BodyJSON(postRequest),
		&web.Retry{
			MaxAttempts:         6,
			ExpectedContentType: "application/json",
			Backoff: backoff.Backoff{
				Jitter: true,
				Min:    30 * time.Second,
				Max:    60 * time.Second,
			}}, web.GetRateLimiter(t.tracker.Name, btnApiRateLimit), headers)

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
		t.log.WithError(err).Errorf("Failed unmarshalling data received: %#v", btnReleaseAsBytes)
		return nil, err
	}

	t.log.Tracef("API GetReleaseInfo Response: %+v", btnInfo)
	return &TorrentInfo{
		Name: btnInfo.Result.ReleaseName,
		Size: btnInfo.Result.Size,
	}, nil
}
