package apis

import (
	"bytes"
	"fmt"
	"gitlab.com/cloudb0x/trackarr/config"
	"gitlab.com/cloudb0x/trackarr/utils/torrent"
	"net/http"
	"strings"
	"time"

	"gitlab.com/cloudb0x/trackarr/cache"
	"gitlab.com/cloudb0x/trackarr/logger"
	webutils "gitlab.com/cloudb0x/trackarr/utils/web"

	"github.com/imroc/req"
	"github.com/jpillora/backoff"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

/* Const */
const TorrentFileTimeout = 30

/* Public */

func Torrent(c echo.Context) error {
	// parse query params
	url := c.QueryParam("url")
	cookie := c.QueryParam("cookie")
	pvr := c.QueryParam("pvr")

	// validate query params
	if pvr == "" {
		pvr = "Unknown"
	}

	log := logger.GetLogger("api").WithFields(logrus.Fields{"client": c.RealIP(), "pvr": pvr})

	if url == "" {
		log.Warn("Torrent request with no URL...")
		return c.String(http.StatusNotAcceptable, "URL was not provided")
	}

	// does this torrent exist in the cache?
	var cacheItem *cache.CacheItem
	cacheItemPresent := false

	cacheItem, cacheItemPresent = cache.GetItem(url)
	if cacheItemPresent && cacheItem.Data != nil {
		// cache item was found and there was torrent bytes
		// this means we can send the torrent directly, as bencode would have already evaluated against torrent data
		log.Infof("Torrent requested: %s (cache: %s)", url, cacheItem.Name)
		return c.Stream(http.StatusOK, "application/x-bittorrent", bytes.NewReader(cacheItem.Data))
	}

	// torrent was not in cache (or had no bytes to send), lets return it directly
	log.Infof("Torrent requested: %s", url)

	// set headers
	headers := req.Header{}
	if cookie != "" {
		headers["Cookie"] = cookie
	}

	// retrieve pvr instance if set
	if pvr != "Unknown" && cacheItem != nil && cacheItem.Release != nil && len(cacheItem.Release.Files) == 0 {
		// the release has no files - meaning bencode was not previously done / evaluated against
		if pvrInstance, ok := config.Pvr[pvr]; ok && pvrInstance.HasFileExpressions {
			// pvr instance was found, and file expressions were present - we need to re-evaluate before sending torrent
			torrentData, err := torrent.GetTorrentDetails(url, TorrentFileTimeout, headers)
			if err != nil {
				// failed to retrieve torrent data
				log.WithError(err).Error("Failed decoding details from torrent file")
				return c.JSON(http.StatusInternalServerError, &ErrorResponse{
					Error:   true,
					Message: fmt.Sprintf("Failed decoding torrent: %v", err),
				})
			}

			// store parsed torrent files in release
			cacheItem.Release.Files = torrentData.Files
			cacheItem.Data = torrentData.Bytes

			// update/store item in cache (for any future calls within N seconds)
			go cache.AddItem(url, &cache.CacheItem{
				Name:    cacheItem.Name,
				Data:    cacheItem.Data,
				Release: cacheItem.Release,
			})

			// evaluate release against expressions (sweep-two)
			// - check ignore expressions
			ignore, err := pvrInstance.ShouldIgnore(cacheItem.Release, log)
			if err != nil {
				log.WithError(err).Warn("Failed checking release on sweep-two against ignore expressions")
				return c.JSON(http.StatusInternalServerError, &ErrorResponse{
					Error:   true,
					Message: fmt.Sprintf("Failed evaluating ignore expressions on sweep-two for pvr %q: %v", pvrInstance.Config.Name, err),
				})
			}

			if ignore {
				log.WithField("name", cacheItem.Name).Warn("Ignoring approved release after sweep-two of ignore expressions")
				return c.JSON(http.StatusNotFound, &ErrorResponse{
					Error:   true,
					Message: fmt.Sprintf("Ignoring release on sweep-two for pvr: %q", pvrInstance.Config.Name),
				})
			}

			// send release
			log.Debug("Release passed sweep-two of ignore expressions")
			return c.Stream(http.StatusOK, "application/x-bittorrent", bytes.NewReader(cacheItem.Data))
		}
	}

	// retrieve torrent stream
	resp, err := webutils.GetResponse(webutils.GET, url, 30, &webutils.Retry{
		MaxAttempts: 5,
		RetryableStatusCodes: []int{
			504,
		},
		ExpectedContentType: "torrent",
		Backoff: backoff.Backoff{
			Jitter: true,
			Min:    500 * time.Millisecond,
			Max:    3 * time.Second,
		}}, headers)
	if err != nil {
		log.WithError(err).Errorf("Failed retrieving torrent stream: %s", url)
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{
			Error:   true,
			Message: fmt.Sprintf("Failed retrieving torrent: %v", err),
		})
	} else if resp.Response().StatusCode != 200 {
		defer resp.Response().Body.Close()

		log.Errorf("Failed retrieving torrent stream: %s (response: %s)", url, resp.Response().Status)
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{
			Error:   true,
			Message: fmt.Sprintf("Failed retrieving torrent: %s", resp.Response().Status),
		})
	}

	// validate response content-type
	respContentType := resp.Response().Header.Get("Content-Type")
	if respContentType == "" || !strings.Contains(respContentType, "torrent") {
		defer resp.Response().Body.Close()

		log.Errorf("Failed retrieving torrent stream: %s (Content-Type: %s)", url, respContentType)
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{
			Error:   true,
			Message: fmt.Sprintf("Failed retrieving torrent: %s", respContentType),
		})
	}

	// return torrent stream
	return c.Stream(http.StatusOK, "application/x-bittorrent", resp.Response().Body)
}
