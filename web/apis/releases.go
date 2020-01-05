package apis

import (
	"net/http"
	"strconv"

	"gitlab.com/cloudb0x/trackarr/database"
	"gitlab.com/cloudb0x/trackarr/database/models"
	"gitlab.com/cloudb0x/trackarr/logger"

	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
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
