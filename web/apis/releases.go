package apis

import (
	"encoding/json"
	"fmt"
	"github.com/l3uddz/trackarr/database"
	"github.com/l3uddz/trackarr/database/models"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

/* Public */

func Releases(c echo.Context) error {
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
	return json.NewEncoder(c.Response()).Encode(releases)
}
