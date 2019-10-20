package ircclient

import (
	irc "github.com/thoj/go-ircevent"
	"strings"
	"time"
)

/* Private */

func (c *IRCClient) handleConnected(event *irc.Event) {
	identified := false

	// send commands
	for _, command := range c.cfg.IRC.Commands {
		cmdToSend := strings.Join(command, " ")
		cmdToSend = strings.TrimLeft(cmdToSend, "/")

		c.log.Debugf("Connected, sending command: %s", cmdToSend)
		c.conn.SendRaw(cmdToSend)

		if strings.Contains(strings.ToLower(cmdToSend), "identify") {
			identified = true
		}

		// sleep a second per command
		time.Sleep(1 * time.Second)
	}

	// join channels
	if !identified {
		// as we have not tried to identify, join announce channels rather than wait for +r mode
		c.log.Debugf("Connected, joining: %s", strings.Join(c.tracker.Channels, ", "))
		for _, channel := range c.tracker.Channels {
			c.conn.Join(channel)
		}
	}
}
