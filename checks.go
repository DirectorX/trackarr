package main

import (
	"gitlab.com/cloudb0x/trackarr/runtime"
)

func startupChecks() {
	// Were there connected clients?
	if connectedClients := len(runtime.Irc); connectedClients < 1 {
		// Alert user that no connections were established
		log.Warn("No tracker connections were established...")
	} else {
		log.Infof("Connected to %d trackers!", connectedClients)
	}
}
