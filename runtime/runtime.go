package runtime

import (
	"github.com/l3uddz/trackarr/loghook"
	"net/http"

	"github.com/l3uddz/trackarr/config"
	"github.com/l3uddz/trackarr/ircclient"
)

var (
	// State
	Tracker = make(map[string]*config.TrackerInstance)
	Pvr     = &config.Pvr
	Irc     = make(map[string]*ircclient.IRCClient)
	Web     *http.Server
	Loghook = loghook.NewLoghooker()
)
