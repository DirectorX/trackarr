package apis

import (
	"github.com/l3uddz/trackarr/logger"
	"github.com/l3uddz/trackarr/runtime"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
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
