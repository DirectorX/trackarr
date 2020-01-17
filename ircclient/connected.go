package ircclient

import (
	"strings"
	"time"

	irc "github.com/thoj/go-ircevent"
)

/* Private */

func (c *IRCClient) handleConnected(event *irc.Event) {
	identified := false

	// reset LastJoined
	c.LastJoined.Store("")

	// send commands
	for _, command := range c.Tracker.Config.IRC.Commands {
		cmdToSend := strings.TrimLeft(command, "/")

		c.log.Debugf("Connected, sending command: %s", cmdToSend)
		c.Conn.SendRaw(cmdToSend)

		if strings.Contains(strings.ToLower(cmdToSend), "identify") {
			identified = true
		}

		// sleep a second per command
		time.Sleep(1 * time.Second)
	}

	// join channels
	if !identified {
		// as we have not tried to identify, join announce channels rather than wait for +r mode
		c.log.Debugf("Connected, joining: %s", strings.Join(c.Tracker.Info.Channels, ", "))
		for _, channel := range c.Tracker.Info.Channels {
			c.Conn.Join(channel)
		}
	}
}
