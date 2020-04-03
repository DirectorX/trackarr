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
	"time"
)

/* Const */
const (
	ptpTorrentUrl = "https://passthepopcorn.me/torrents.php?torrentid="
	ptpTimeout    = 30
)

/* Var */
type Ptp struct {
	log     *logrus.Entry
	tracker *config.TrackerInstance
}

/* Interface */

func (t *Ptp) GetReleaseInfo(torrentId string) (*TorrentInfo, error) {
	// prepare request
	apiUser, err := maps.GetStringMapValue(t.tracker.Config.Settings, "api_user", false)
	if err != nil {
		t.log.WithError(err).Error("api_user value missing")
		return nil, errors.Wrap(err, "api_user value missing")
	}
	apiKey, err := maps.GetStringMapValue(t.tracker.Config.Settings, "api_key", false)
	if err != nil {
		t.log.WithError(err).Error("api_key value missing")
		return nil, errors.Wrap(err, "api_key value missing")
	}

	headers := req.Header{
		"ApiUser": apiUser,
		"ApiKey":  apiKey,
	}

	// send request
	ptpReleaseAsBytes, err := web.GetBodyBytes(web.GET, fmt.Sprintf("%s%s", ptpTorrentUrl, torrentId), ptpTimeout,
		&web.Retry{
			MaxAttempts:         6,
			ExpectedContentType: "application/json",
			Backoff: backoff.Backoff{
				Jitter: true,
				Min:    3 * time.Second,
				Max:    10 * time.Second,
			}}, headers)
	if err != nil {
		return nil, errors.Wrapf(err, "failed retrieving torrent bytes from: %s", torrentId)
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
		t.log.WithError(err).Errorf("Failed unmarshalling data received: %#v", ptpReleaseAsBytes)
		return nil, err
	}

	t.log.Tracef("Torrent Lookup Response: %+v", ptpInfo)

	// find torrent in parsed response
	for _, v := range ptpInfo.Torrents {
		if v.Id == torrentId {
			return &TorrentInfo{
				Name: v.ReleaseName,
				Size: v.Size,
			}, nil
		}
	}

	return nil, fmt.Errorf("no release found with id: %s", torrentId)
}
