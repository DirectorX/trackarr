package ircclient

import (
	listutils "github.com/l3uddz/trackarr/utils/lists"
	irc "github.com/thoj/go-ircevent"
)

/* Private */

func (c *IRCClient) handleMessage(event *irc.Event) {
	channelName := event.Arguments[0]

	// ignore messages if not from a known channel / announcer
	if !listutils.StringListContains(c.tracker.Channels, channelName, false) {
		c.log.Tracef("Ignoring message from channel %s -> %s", channelName, event.Message())
		return
	} else if !listutils.StringListContains(c.tracker.Announcers, event.Nick, false) {
		c.log.Tracef("Ignoring message from announcer %s -> %s", event.User, event.Message())
		return
	}

	// clean message
	cleanMessage := c.cleanMessage(event.Message())
	c.log.Tracef("%s -> %s", channelName, cleanMessage)

	// queue message
	if err := c.processor.QueueLine(channelName, cleanMessage); err != nil {
		c.log.WithError(err).Errorf("Failed queueing line for processing: %q", cleanMessage)
		return
	}
}

func (c IRCClient) cleanMessage(message string) string {
	return c.cleanRxp.ReplaceAllString(message, "")
}
