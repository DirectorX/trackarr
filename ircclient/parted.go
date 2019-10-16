package ircclient

import (
	irc "github.com/thoj/go-ircevent"
	"strings"
)

/* Private */

func (c *IRCClient) handleParted(event *irc.Event) {
	if event.Nick != c.conn.GetNick() {
		// we are not interested in parted messages for other parties
		return
	}

	// determine channel name
	channelName := "Unknown"
	if len(event.Arguments) >= 1 && strings.HasPrefix(event.Arguments[0], "#") {
		// we have the channel name
		channelName = event.Arguments[0]
	}

	c.log.Infof("Parted: %s", channelName)
}
