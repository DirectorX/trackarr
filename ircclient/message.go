package ircclient

import (
	"strings"

	listutils "github.com/l3uddz/trackarr/utils/lists"

	irc "github.com/thoj/go-ircevent"
)

/* Private */

func (c *IRCClient) handleMessage(event *irc.Event) {
	// determine channel name
	channelName := "Unknown"
	if len(event.Arguments) >= 1 && strings.HasPrefix(event.Arguments[0], "#") {
		// we have the channel name
		channelName = event.Arguments[0]
	}

	// ignore messages if not from a known channel / announcer
	if !listutils.StringListContains(c.Tracker.Info.Channels, channelName, false) {
		c.log.Debugf("Ignoring message from %s -> %s", channelName, event.Message())
		return
	} else if !listutils.StringListContains(c.Tracker.Info.Announcers, event.Nick, false) {
		c.log.Debugf("Ignoring message from announcer %s -> %s", event.User, event.Message())
		return
	}

	// clean message
	cleanMessage := c.cleanMessage(event.Message())
	c.log.Tracef("%s -> %s", channelName, cleanMessage)

	// queue message
	if err := c.Processor.QueueLine(channelName, cleanMessage); err != nil {
		c.log.WithError(err).Errorf("Failed queueing line for processing: %q", cleanMessage)
		return
	}
}

func (c IRCClient) cleanMessage(message string) string {
	return messageClean.ReplaceAllString(message, "")
}
