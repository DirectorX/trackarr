package main

import (
	"gitlab.com/cloudb0x/trackarr/runtime"
)

func startupChecks() {
	// Were there connected clients?
	if connectedClients := len(runtime.Irc); connectedClients < 1 {
		log.Fatal("Failed to establish a connection to any of the enabled trackers...")
	} else {
		log.Infof("Connected to %d trackers!", connectedClients)
	}
}
