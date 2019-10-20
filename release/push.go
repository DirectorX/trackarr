package release

import (
	"github.com/imroc/req"
	"github.com/l3uddz/trackarr/config"
	"github.com/l3uddz/trackarr/database"
	"github.com/l3uddz/trackarr/database/models"
	"github.com/l3uddz/trackarr/utils/web"
	"strconv"
	"strings"
)

/* Structs */

type (
	pvrResponse struct {
		Approved   bool
		Rejections []string
	}
	pvrResponses struct {
		Responses []pvrResponse
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

func (r *TrackerRelease) Push(pvr *config.PvrConfiguration) {
	r.Log.Debugf("Pushing: %s (pvr: %s)", r.TorrentName, pvr.Name)

	// prepare request
	pvrRequest := pushRequest{
		Title:            r.TorrentName,
		DownloadUrl:      r.TorrentURL,
		Size:             "0",
		Indexer:          r.Tracker.LongName,
		DownloadProtocol: "torrent",
		Protocol:         "torrent",
		PublishDate:      r.ReleaseTime,
	}

	if r.SizeBytes > 0 {
		pvrRequest.Size = strconv.FormatInt(r.SizeBytes, 10)
	}

	requestUrl := web.JoinURL(pvr.URL, "api/release/push")
	headers := req.Header{
		"X-Api-Key": pvr.ApiKey,
	}

	// send request
	resp, err := web.GetResponse(web.POST, requestUrl, 30, req.BodyJSON(&pvrRequest), headers)
	if err != nil {
		r.Log.WithError(err).Errorf("Failed pushing: %s (pvr: %s)", r.TorrentName, pvr.Name)
		return
	}

	defer resp.Response().Body.Close()

	// validate response
	if resp.Response().StatusCode != 200 {
		r.Log.Errorf("Failed pushing: %s (pvr: %s - response: %q)", r.TorrentName, pvr.Name,
			resp.Response().Status)
		return
	}

	// decode response
	pvrResp := &pvrResponse{}
	if err := resp.ToJSON(&pvrResp); err != nil {
		r.Log.WithError(err).Errorf("Failed decoding push response: %s (pvr: %s)", r.TorrentName, pvr.Name)
		return
	}

	// log result
	r.Log.Infof("Pushed: %s (pvr: %s - approved: %v)", r.TorrentName, pvr.Name, pvrResp.Approved)
	if len(pvrResp.Rejections) > 0 {
		r.Log.Debugf("Push rejected: %s (pvr: %s - reasons: %q)", r.TorrentName, pvr.Name,
			strings.Join(pvrResp.Rejections, ", "))
	}

	// save to database
	r.Log.Tracef("Creating release in database...")

	dbRelease, err := models.NewPushedRelease(database.DB, r.TorrentName, r.TrackerName, pvr.Name, pvrResp.Approved)
	if err != nil {
		r.Log.WithError(err).Errorf("Failed saving release in database: %q", r.TorrentName)
	}
	database.DB.Save(&dbRelease)
}
