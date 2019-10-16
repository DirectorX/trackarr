package ircclient

import (
	listutils "github.com/l3uddz/trackarr/utils/lists"
	irc "github.com/thoj/go-ircevent"
	"strings"
)

/* Private */

func (c *IRCClient) handleJoined(event *irc.Event) {
	// determine channel name
	channelName := "Unknown"
	if len(event.Arguments) >= 2 && strings.HasPrefix(event.Arguments[1], "#") {
		// we have the channel name
		channelName = event.Arguments[1]
	}
	c.log.Infof("Joined: %s", channelName)

	// is this an announce channel?
	if channelName != "Unknown" && !listutils.StringListContains(c.tracker.Channels, channelName, false) {
		// this is not an announce channel, lets leave.
		c.log.Debugf("Leaving non-announce channel: %s", channelName)
		c.conn.Part(channelName)
	}
}

func (c *IRCClient) handleJoinFailure(event *irc.Event) {
	channelName := "Unknown"
	if len(event.Arguments) >= 2 && strings.HasPrefix(event.Arguments[1], "#") {
		// we have the channel name
		channelName = event.Arguments[1]
	}

	c.log.Warnf("Failed joining: %s, reason: %s", channelName, event.Message())
}
