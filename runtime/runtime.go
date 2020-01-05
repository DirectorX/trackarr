package runtime

import (
	"gitlab.com/cloudb0x/trackarr/config"
	"gitlab.com/cloudb0x/trackarr/ircclient"
	"gitlab.com/cloudb0x/trackarr/loghook"
	"gitlab.com/cloudb0x/trackarr/tasks"

	"github.com/labstack/echo/v4"
)

var (
	// State
	Loghook = loghook.NewLoghooker()
	Tracker = make(map[string]*config.TrackerInstance)
	Irc     = make(map[string]*ircclient.IRCClient)
	Web     *echo.Echo
	Tasks   *tasks.Tasks
)
