package ircclient

import (
	irc "github.com/thoj/go-ircevent"
	"strings"
)

/* Private */

func (c *IRCClient) handleConnected(event *irc.Event) {
	c.log.Debugf("Connected, joining: %s", strings.Join(c.tracker.Channels, ", "))
	for _, channel := range c.tracker.Channels {
		c.conn.Join(channel)
	}
}
