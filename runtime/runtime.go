package runtime

import (
	"net/http"

	"gitlab.com/cloudb0x/trackarr/config"
	"gitlab.com/cloudb0x/trackarr/ircclient"
	"gitlab.com/cloudb0x/trackarr/loghook"
	"gitlab.com/cloudb0x/trackarr/tasks"
)

var (
	// State
	Loghook = loghook.NewLoghooker()
	Tracker = make(map[string]*config.TrackerInstance)
	Irc     = make(map[string]*ircclient.IRCClient)
	Web     *http.Server
	Tasks   *tasks.Tasks
)
