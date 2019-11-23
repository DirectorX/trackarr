package runtime

import (
	"github.com/l3uddz/trackarr/loghook"
	"github.com/l3uddz/trackarr/tasks"
	"net/http"

	"github.com/l3uddz/trackarr/config"
	"github.com/l3uddz/trackarr/ircclient"
)

var (
	// State
	Loghook = loghook.NewLoghooker()
	Tracker = make(map[string]*config.TrackerInstance)
	Irc     = make(map[string]*ircclient.IRCClient)
	Web     *http.Server
	Tasks   *tasks.Tasks
)
