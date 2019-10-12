package ircclient

import (
	listutils "github.com/l3uddz/trackarr/utils/lists"
	irc "github.com/thoj/go-ircevent"
)

/* Private */

func (c *IRCClient) handleMessage(event *irc.Event) {
	channelName := event.Arguments[0]

	// ignore messages if not from an known channel / announcer
	if !listutils.StringListContains(c.tracker.Channels, channelName, false) {
		log.Tracef("Ignoring message from channel %s -> %s", channelName, event.Message())
		return
	} else if !listutils.StringListContains(c.tracker.Announcers, event.User, false) {
		log.Tracef("Ignoring message from announcer %s -> %s", event.User, event.Message())
		return
	}

	// clean message
	cleanMessage := c.cleanMessage(event.Message())
	c.log.Tracef("%s -> %s", channelName, cleanMessage)

	// process message
	if err := c.processor.ProcessLine(cleanMessage); err != nil {
		c.log.WithError(err).Errorf("Failed processing line from %s -> %s", channelName, cleanMessage)
		return
	}
}

func (c IRCClient) cleanMessage(message string) string {
	return c.cleanRxp.ReplaceAllString(message, "")
}
