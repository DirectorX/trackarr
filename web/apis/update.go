package apis

import (
	"github.com/l3uddz/trackarr/config"
	"github.com/l3uddz/trackarr/logger"
	"github.com/l3uddz/trackarr/utils/version"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
)

/* Structs */

type UpdateResponse struct {
	UpdateAvailable bool   `json:"update_available" xml:"update_available"`
	LatestVersion   string `json:"latest_version" xml:"latest_version"`
}

/* Public */

func UpdateStatus(c echo.Context) error {
	// log
	log := logger.GetLogger("api").WithFields(logrus.Fields{"client": c.RealIP()})

	// is there an update available?
	usingLatest, latestVersion := version.IsLatestGitlabVersion(
		"https://gitlab.com/api/v4/projects/15385789/releases", "", config.Build.Version)
	log.Debugf("Latest version: %q, update available: %v", latestVersion, !usingLatest)

	// return response
	return c.JSON(http.StatusOK, &UpdateResponse{
		UpdateAvailable: !usingLatest,
		LatestVersion:   latestVersion,
	})
}
