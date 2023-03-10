package release

import (
	"gitlab.com/cloudb0x/trackarr/utils/defaults"
	"strings"
	"time"

	"gitlab.com/cloudb0x/trackarr/config"
	"gitlab.com/cloudb0x/trackarr/database"
	"gitlab.com/cloudb0x/trackarr/database/models"
	"gitlab.com/cloudb0x/trackarr/utils/web"
	"gitlab.com/cloudb0x/trackarr/ws"

	"github.com/imroc/req"
	"github.com/jpillora/backoff"
)

/* Structs */

type (
	pvrResponse struct {
		Approved   bool
		Rejections []string
	}
	pushRequest struct {
		Title            string `json:"title"`
		DownloadUrl      string `json:"downloadUrl"`
		Size             int64  `json:"size"`
		Indexer          string `json:"indexer"`
		DownloadProtocol string `json:"downloadProtocol"`
		Protocol         string `json:"protocol"`
		PublishDate      string `json:"publishDate"`
	}
)

/* Public */

func (r *Release) Push(pvr *config.PvrConfig, delay *int64, torrentUrl *string) {
	if delay != nil && *delay > 0 {
		r.Log.Debugf("Delaying: %s (pvr: %s)", r.Info.TorrentName, pvr.Name)
		time.Sleep(time.Duration(*delay) * time.Second)
	}
	r.Log.Debugf("Pushing: %s (pvr: %s)", r.Info.TorrentName, pvr.Name)

	// prepare request
	pvrRequest := pushRequest{
		Title:            r.Info.TorrentName,
		DownloadUrl:      *torrentUrl,
		Size:             0,
		Indexer:          r.Tracker.Info.LongName,
		DownloadProtocol: "torrent",
		Protocol:         "torrent",
		PublishDate:      r.Info.ReleaseTime,
	}

	if r.Info.SizeBytes > 0 {
		pvrRequest.Size = r.Info.SizeBytes
	}

	requestUrl := ""
	if !strings.Contains(pvr.URL, "/api/") {
		requestUrl = web.JoinURL(pvr.URL, "api/v3/release/push")
	} else {
		requestUrl = pvr.URL
	}

	headers := req.Header{
		"X-Api-Key": pvr.ApiKey,
	}

	// send request
	resp, err := web.GetResponse(web.POST, requestUrl, defaults.GetOrDefaultInt(&pvr.Timeout, 30),
		req.BodyJSON(&pvrRequest), &web.Retry{
			MaxAttempts: 6,
			RetryableStatusCodes: []int{
				504,
			},
			Backoff: backoff.Backoff{
				Jitter: true,
				Min:    500 * time.Millisecond,
				Max:    10 * time.Second,
			},
		}, headers)
	if err != nil {
		r.Log.WithError(err).Errorf("Failed pushing: %s (pvr: %s)", r.Info.TorrentName, pvr.Name)
		return
	}

	defer web.DrainAndClose(resp.Response().Body)

	// validate response
	if resp.Response().StatusCode != 200 {
		r.Log.Errorf("Failed pushing: %s (pvr: %s - response: %q)", r.Info.TorrentName, pvr.Name,
			resp.Response().Status)
		return
	}

	// decode response
	pvrResp := make([]pvrResponse, 0)
	if err := resp.ToJSON(&pvrResp); err != nil {
		r.Log.WithError(err).Errorf("Failed decoding push response: %s (pvr: %s)", r.Info.TorrentName, pvr.Name)
		return
	} else if len(pvrResp) < 1 {
		r.Log.Errorf("Failed processing push response: %s (pvr: %s)", r.Info.TorrentName, pvr.Name)
		return
	}

	// log result
	r.Log.Infof("Pushed: %s (pvr: %s - approved: %v)", r.Info.TorrentName, pvr.Name, pvrResp[0].Approved)
	if len(pvrResp[0].Rejections) > 0 {
		r.Log.Tracef("Push rejected: %s (pvr: %s - reasons: %q)", r.Info.TorrentName, pvr.Name,
			strings.Join(pvrResp[0].Rejections, ", "))
	}

	// save to database
	dbRelease, err := models.NewPushedRelease(database.DB, r.Info.TorrentName, r.Info.TrackerName, pvr.Name,
		pvrResp[0].Approved)
	if err != nil {
		r.Log.WithError(err).Errorf("Failed saving release in database: %q", r.Info.TorrentName)
		return
	}

	// broadcast release to releases topic
	pushMessage := &ws.WebsocketMessage{
		Type: "release",
		Data: dbRelease,
	}

	jsonData, err := pushMessage.ToJsonString()
	if err == nil {
		ws.BroadcastTopic("releases", jsonData)
	}
}
