package apis

import (
	"net/http"

	"gitlab.com/cloudb0x/trackarr/config"
	"gitlab.com/cloudb0x/trackarr/logger"
	"gitlab.com/cloudb0x/trackarr/version"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

/* Structs */

type UpdateResponse struct {
	UpdateAvailable bool   `json:"update_available" xml:"update_available"`
	LatestVersion   string `json:"latest_version" xml:"latest_version"`
	CurrentVersion  string `json:"current_version" xml:"current_version"`
}

/* Public */

func UpdateStatus(c echo.Context) error {
	// log
	log := logger.GetLogger("api").WithFields(logrus.Fields{"client": c.RealIP()})

	// is there an update available?
	usingLatest, latestVersion := version.Trackarr.IsLatest()
	log.Debugf("Latest version: %q, update available: %v", latestVersion, !usingLatest)

	// return response
	return c.JSON(http.StatusOK, &UpdateResponse{
		UpdateAvailable: !usingLatest,
		LatestVersion:   latestVersion.String(),
		CurrentVersion:  config.Build.Version,
	})
}
