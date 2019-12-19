package apis

import (
	"bytes"
	"fmt"
	"github.com/imroc/req"
	"github.com/jpillora/backoff"
	"github.com/l3uddz/trackarr/cache"
	"github.com/l3uddz/trackarr/logger"
	webutils "github.com/l3uddz/trackarr/utils/web"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

var ()

/* Public */

func Torrent(c echo.Context) error {
	// log
	log := logger.GetLogger("api").WithFields(logrus.Fields{"client": c.RealIP()})

	// parse query params
	url := c.QueryParam("url")
	cookie := c.QueryParam("cookie")

	// validate query params
	if url == "" {
		log.Warn("Torrent request with no URL...")
		return c.String(http.StatusNotAcceptable, "URL was not provided")
	}

	// does this torrent exist in the cache?
	if cacheItem, ok := cache.GetItem(url); ok {
		log.Infof("Torrent requested: %s (cache: %s)", url, cacheItem.Name)
		return c.Stream(http.StatusOK, "application/x-bittorrent", bytes.NewReader(cacheItem.Data))
	}

	// torrent was not in cache, lets return it directly
	log.Infof("Torrent requested: %s", url)

	// set headers
	headers := req.Header{}
	if cookie != "" {
		headers["Cookie"] = cookie
	}

	// retrieve torrent stream
	resp, err := webutils.GetResponse(webutils.GET, url, 30, &webutils.Retry{
		MaxAttempts:         5,
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
