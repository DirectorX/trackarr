package ptp

import (
	"encoding/json"
	"fmt"
	"github.com/imroc/req"
	"github.com/jpillora/backoff"
	"github.com/pkg/errors"
	"gitlab.com/cloudb0x/trackarr/config"
	"gitlab.com/cloudb0x/trackarr/logger"
	"gitlab.com/cloudb0x/trackarr/utils/maps"
	"gitlab.com/cloudb0x/trackarr/utils/web"
	"time"
)

//Contains all the meta-info data from the original torrent file
type ReleaseInfo struct {
	Page          string // "Details",
	Result        string // "OK",
	GroupId       string // "753",
	Name          string // "The Godfather: Part III",
	Year          string // "1990",
	CoverImage    string // "https://ptpimg.me/4n8r45.jpg",
	AuthKey       string // "2658d23b373c717ed499fb8355928479",
	PassKey       string // "bp0axqzxyo2cwswk0phvyoqn0gifz6tl",
	TorrentId     string //: "766438",
	ImdbId        string // "0099674",
	ImdbRating    string // "7.6",
	ImdbVoteCount int64  // 339853,
	Torrents      []TorrentInfo
}
type TorrentInfo struct {
	Id            string // "1368",
	InfoHash      string // "FD499DC35AEC8383E795DF2A8A42CC65C814DF85",
	Quality       string // "Standard Definition",
	Source        string // "DVD",
	Container     string // "AVI",
	Codec         string // "XviD",
	Resolution    string // "608x352",
	Size          string // "1468434432",
	Scene         bool   // false,
	UploadTime    string //: "2008-10-23 01:15:04",
	Snatched      string // "542",
	Seeders       string // "11",
	Leechers      string //: "0",
	ReleaseName   string // "The Godfather Part 3",
	ReleaseGroup  string // null,
	Checked       bool   // true,
	GoldenPopcorn bool   // false
}

var (
	// Logging
	log = logger.GetLogger("ptp")
)

const Url = "https://passthepopcorn.me/torrents.php?torrentid="

func GetReleaseDetails(torrentId string, tracker *config.TrackerInstance, timeout int) (*TorrentInfo, error) {
	apiUser, err := maps.GetStringMapValue(tracker.Config.Settings, "api_user", false)
	if err != nil {
		log.WithError(err).Error("PTP Api User value missing")
		return nil, errors.Wrap(err, "PTP Api User value missing")
	}
	apiKey, err := maps.GetStringMapValue(tracker.Config.Settings, "api_key", false)
	if err != nil {
		log.WithError(err).Error("PTP Api Key value missing")
		return nil, errors.Wrap(err, "PTP Api Key value missing")
	}

	headers := req.Header{
		"ApiUser": apiUser,
		"ApiKey":  apiKey,
	}

	// retrieve releaseBytes
	ptpReleaseAsBytes, err := web.GetBodyBytes(web.GET, fmt.Sprintf("%s%s", Url, torrentId), timeout, &web.Retry{
		MaxAttempts:         6,
		ExpectedContentType: "application/json",
		Backoff: backoff.Backoff{
			Jitter: true,
			Min:    10 * time.Second,
			Max:    60 * time.Second,
		}}, headers)
	if err != nil {
		return nil, errors.Wrapf(err, "failed retrieving torrent bytes from: %s", torrentId)
	}
	var ptpInfo ReleaseInfo
	if err := json.Unmarshal(ptpReleaseAsBytes, &ptpInfo); err != nil {
		log.WithError(err).Errorf("Failed unmarshalling data received: %#v", ptpReleaseAsBytes)
		return nil, err
	}
	log.Tracef("New PTP release : %s", ptpInfo)
	for _, v := range ptpInfo.Torrents {
		if v.Id == torrentId {
			return &v, nil
		}
	}
	return nil, errors.Errorf("No release found with id %s", torrentId)
}
