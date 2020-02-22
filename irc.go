package main

import (
	"gitlab.com/cloudb0x/trackarr/ircclient"
	"gitlab.com/cloudb0x/trackarr/runtime"
)

// Initialize IRC clients for trackers
func initIRC() {
	for tName, t := range runtime.Tracker {
		// load irc client
		log.Debugf("Initializing irc client: %s", tName)
		t2 := t
		c, err := ircclient.New(t2)
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
