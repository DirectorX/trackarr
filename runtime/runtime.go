package runtime

import (
	"net/http"

	"github.com/l3uddz/trackarr/config"
	"github.com/l3uddz/trackarr/ircclient"
)

var (
	// State
	Tracker = make(map[string]*config.TrackerInstance)
	Pvr     = make(map[string]*config.PvrInstance)
	Irc     = make(map[string]*ircclient.IRCClient)
	Web     *http.Server
)
