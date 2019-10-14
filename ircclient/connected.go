package ircclient

import (
	irc "github.com/thoj/go-ircevent"
	"strings"
)

/* Private */

func (c *IRCClient) handleConnected(event *irc.Event) {
	// send commands
	for _, command := range c.cfg.IRC.Commands {
		cmdToSend := strings.Join(command, " ")
		cmdToSend = strings.TrimLeft(cmdToSend, "/")

		c.log.Debugf("Connected, sending command: %s", cmdToSend)
		c.conn.SendRaw(cmdToSend)
	}

	// join channels
	c.log.Debugf("Connected, joining: %s", strings.Join(c.tracker.Channels, ", "))
	for _, channel := range c.tracker.Channels {
		c.conn.Join(channel)
	}
}
