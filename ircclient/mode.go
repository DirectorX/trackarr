package ircclient

import (
	irc "github.com/thoj/go-ircevent"
	"strings"
	"time"
)

/* Private */

func (c *IRCClient) handleMode(event *irc.Event) {
	if !strings.Contains(event.Raw, c.conn.GetNick()) || !strings.Contains(event.Raw, "+r") {
		// the raw message did not contain our nick and a +r
		return
	}

	// sleep 2 seconds before joining after +r
	time.Sleep(2 * time.Second)

	// join announce channels
	c.log.Debugf("Identified, joining: %s", strings.Join(c.tracker.Channels, ", "))
	for _, channel := range c.tracker.Channels {
		c.conn.Join(channel)
	}
}
