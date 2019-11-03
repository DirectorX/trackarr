package release

import (
	"strconv"
	"strings"
	"time"

	"github.com/l3uddz/trackarr/config"
	"github.com/l3uddz/trackarr/database"
	"github.com/l3uddz/trackarr/database/models"
	"github.com/l3uddz/trackarr/utils/web"

	"github.com/imroc/req"
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
		Size             string `json:"size"`
		Indexer          string `json:"indexer"`
		DownloadProtocol string `json:"downloadProtocol"`
		Protocol         string `json:"protocol"`
		PublishDate      string `json:"publishDate"`
	}
)

/* Public */

func (r *Release) Push(pvr *config.PvrConfig, delay *int64) {
	if delay != nil && *delay > 0 {
		r.Log.Debugf("Delaying: %s (pvr: %s)", r.Info.TorrentName, pvr.Name)
		time.Sleep(time.Duration(*delay) * time.Second)
	}
	r.Log.Debugf("Pushing: %s (pvr: %s)", r.Info.TorrentName, pvr.Name)

	// prepare request
	pvrRequest := pushRequest{
		Title:            r.Info.TorrentName,
		DownloadUrl:      r.Info.TorrentURL,
		Size:             "0",
		Indexer:          r.Tracker.Info.LongName,
		DownloadProtocol: "torrent",
		Protocol:         "torrent",
		PublishDate:      r.Info.ReleaseTime,
	}

	if r.Info.SizeBytes > 0 {
		pvrRequest.Size = strconv.FormatInt(r.Info.SizeBytes, 10)
	}

	requestUrl := ""
	if strings.Index(pvr.URL, "/api/") == -1 {
		requestUrl = web.JoinURL(pvr.URL, "api/release/push")
	} else {
		requestUrl = pvr.URL
	}

	headers := req.Header{
		"X-Api-Key": pvr.ApiKey,
	}

	// send request
	resp, err := web.GetResponse(web.POST, requestUrl, 30, req.BodyJSON(&pvrRequest), headers)
	if err != nil {
		r.Log.WithError(err).Errorf("Failed pushing: %s (pvr: %s)", r.Info.TorrentName, pvr.Name)
		return
	}

	defer resp.Response().Body.Close()

	// validate response
	if resp.Response().StatusCode != 200 {
		r.Log.Errorf("Failed pushing: %s (pvr: %s - response: %q)", r.Info.TorrentName, pvr.Name,
			resp.Response().Status)
		return
	}

	// decode response
	pvrResp := &pvrResponse{}
	if err := resp.ToJSON(&pvrResp); err != nil {
		r.Log.WithError(err).Errorf("Failed decoding push response: %s (pvr: %s)", r.Info.TorrentName, pvr.Name)
		return
	}

	// log result
	r.Log.Infof("Pushed: %s (pvr: %s - approved: %v)", r.Info.TorrentName, pvr.Name, pvrResp.Approved)
	if len(pvrResp.Rejections) > 0 {
		r.Log.Tracef("Push rejected: %s (pvr: %s - reasons: %q)", r.Info.TorrentName, pvr.Name,
			strings.Join(pvrResp.Rejections, ", "))
	}

	// save to database
	r.Log.Tracef("Creating release in database...")

	dbRelease, err := models.NewPushedRelease(database.DB, r.Info.TorrentName, r.Info.TrackerName, pvr.Name, pvrResp.Approved)
	if err != nil {
		r.Log.WithError(err).Errorf("Failed saving release in database: %q", r.Info.TorrentName)
	}
	database.DB.Save(&dbRelease)
}
