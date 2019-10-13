package ircclient

import (
	irc "github.com/thoj/go-ircevent"
	"strings"
)

/* Private */

func (c *IRCClient) handleJoined(event *irc.Event) {
	channelName := "Unknown"
	if len(event.Arguments) >= 1 && strings.HasPrefix(event.Arguments[1], "#") {
		// we have the channel name
		channelName = event.Arguments[1]
	}
	c.log.Infof("Joined: %s", channelName)
}

func (c *IRCClient) handleJoinFailure(event *irc.Event) {
	channelName := "Unknown"
	if len(event.Arguments) >= 1 && strings.HasPrefix(event.Arguments[1], "#") {
		// we have the channel name
		channelName = event.Arguments[1]
	}

	c.log.Warnf("Failed joining: %s, reason: %s", channelName, event.Message())
}
