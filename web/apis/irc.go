package apis

import (
	"net/http"

	"gitlab.com/cloudb0x/trackarr/logger"
	"gitlab.com/cloudb0x/trackarr/runtime"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

/* Struct */

type IrcTrackerStatus struct {
	Connected     bool   `json:"connected"`
	LastJoined    string `json:"last_joined"`
	LastAnnounced string `json:"last_announced"`
}

/* Public */

func IrcStatus(c echo.Context) error {
	// log
	log := logger.GetLogger("api").WithFields(logrus.Fields{"client": c.RealIP()})

	// build map of irc statuses
	clientStatuses := map[string]IrcTrackerStatus{}
	for _, client := range runtime.Irc {
		clientStatuses[client.Tracker.Name] = IrcTrackerStatus{
			Connected:     client.Conn.Connected(),
			LastJoined:    client.LastJoined.Load(),
			LastAnnounced: client.LastAnnounced.Load(),
		}
	}

	log.Debugf("%d irc client statuses requested", len(clientStatuses))

	// return response
	return c.JSON(http.StatusOK, &clientStatuses)
}
