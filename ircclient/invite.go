package ircclient

import (
	irc "github.com/thoj/go-ircevent"
	"gitlab.com/cloudb0x/trackarr/utils/lists"
	"time"
)

/* Private */

func (c *IRCClient) handleInvite(event *irc.Event) {
	// parse invite channel
	channel := "Unknown"
	if len(event.Arguments) >= 2 && len(event.Arguments[1]) >= 2 && event.Arguments[1][0] == '#' {
		channel = event.Arguments[1]
	}

	// validate invite to an expected announce channel
	if !lists.StringListContains(c.Tracker.Info.Channels, channel, false) {
		c.log.Warnf("Ignoring invite by %q to channel: %q", event.Nick, channel)
		return
	}

	// is invite from an expected announcer ?
	if !lists.StringListContains(c.Tracker.Info.Announcers, event.Nick, false) {
		c.log.Warnf("Ignoring invite to %q by: %q", channel, event.Nick)
		return
	}

	// sleep 2 seconds before joining
	time.Sleep(2 * time.Second)

	// join announce channel
	c.log.Debugf("Invited by %q, joining: %q", event.Nick, channel)
	c.Conn.Join(channel)
}
