package apis

import (
	"net/http"

	"gitlab.com/cloudb0x/trackarr/logger"
	"gitlab.com/cloudb0x/trackarr/runtime"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

/* Public */

func IrcStatus(c echo.Context) error {
	// log
	log := logger.GetLogger("api").WithFields(logrus.Fields{"client": c.RealIP()})

	// build map of irc statuses
	clientStatuses := map[string]bool{}
	for _, client := range runtime.Irc {
		clientStatuses[client.Tracker.Name] = client.Conn.Connected()
	}

	log.Debugf("%d irc client statuses requested", len(clientStatuses))

	// return response
	return c.JSON(http.StatusOK, &clientStatuses)
}
