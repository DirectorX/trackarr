package apis

import (
	"github.com/json-iterator/go"
	"github.com/l3uddz/trackarr/database"
	"github.com/l3uddz/trackarr/database/models"
	"github.com/l3uddz/trackarr/logger"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

/* Public */

func Releases(c echo.Context) error {
	// log
	log := logger.GetLogger("api").WithFields(logrus.Fields{"client": c.RealIP()})

	// parse parameters
	countParam := c.QueryParam("count")
	approvedParam := c.QueryParam("approved")

	releaseCount, err := strconv.Atoi(countParam)
	if err != nil {
		releaseCount = 0
	}

	releaseApproved := false
	if approvedParam == "1" || approvedParam == "true" {
		releaseApproved = true
	}

	log.Debugf("Releases requested, approved: %v", releaseApproved)

	// retrieve releases
	var releases []*models.PushedRelease

	if releaseApproved {
		releases = models.GetLatestApprovedReleases(database.DB, releaseCount)
	} else {
		releases = models.GetLatestPushedReleases(database.DB, releaseCount)
	}

	// return response
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(http.StatusOK)

	enc := json.NewEncoder(c.Response())
	enc.SetEscapeHTML(false)
	return enc.Encode(releases)
}
