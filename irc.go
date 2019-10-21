package main

import (
	"github.com/l3uddz/trackarr/ircclient"
	"github.com/l3uddz/trackarr/runtime"
)

// Initialize IRC clients for trackers
func initIRC() {
	for tName, t := range runtime.Tracker {
		// load irc client
		log.Debugf("Initializing IRC client: %s", tName)
		c, err := ircclient.New(t)
		if err != nil {
			log.WithError(err).Errorf("Failed initializing irc client for tracker: %s", tName)

			continue
		}
		log.Debugf("Initialized irc client: %s", tName)

		// start client
		if err := c.Start(); err != nil {
			log.WithError(err).Errorf("Failed starting irc client for tracker: %s", tName)

			continue
		}

		runtime.Irc[tName] = c
	}
}
