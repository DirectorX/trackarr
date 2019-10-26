package apis

import (
	"encoding/json"
	"fmt"
	"github.com/l3uddz/trackarr/database"
	"github.com/l3uddz/trackarr/database/models"
	"github.com/l3uddz/trackarr/logger"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
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
		return c.JSON(http.StatusBadRequest, &Response{
			Error:   true,
			Message: fmt.Sprintf("Failed parsing required count parameter: %v", err),
		})
	}

	releaseApproved := false
	if approvedParam == "1" || approvedParam == "true" {
		releaseApproved = true
	}

	log.Debugf("%d releases requested, approved: %v", releaseCount, releaseApproved)

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
