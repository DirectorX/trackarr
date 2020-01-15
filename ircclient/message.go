package ircclient

import (
	"strings"
	"time"

	listutils "gitlab.com/cloudb0x/trackarr/utils/lists"

	irc "github.com/thoj/go-ircevent"
)

/* Private */

func (c *IRCClient) handleMessage(event *irc.Event) {
	// determine channel name
	channelName := "unknown"
	if len(event.Arguments) > 0 && strings.HasPrefix(event.Arguments[0], "#") {
		// we have the channel name
		channelName = event.Arguments[0]
	}

	// ignore messages if not from a known channel / announcer
	if !listutils.StringListContains(c.Tracker.Info.Channels, channelName, false) {
		c.log.Debugf("Ignoring message from channel: %q -> %s", channelName, event.Message())
		return
	}

	if !listutils.StringListContains(c.Tracker.Info.Announcers, event.Nick, false) {
		c.log.Debugf("Ignoring message from nick: %q (u: %s) -> %s", event.Nick, event.User, event.Message())
		return
	}

	// clean message
	cleanMessage := c.cleanMessage(event.Message())
	c.log.Tracef("%s -> %s", channelName, cleanMessage)

	// queue message
	if err := c.Processor.QueueLine(strings.ToLower(channelName), cleanMessage); err != nil {
		c.log.WithError(err).Errorf("Failed queueing line for processing: %q", cleanMessage)
		return
	}

	// update last announced
	c.LastAnnounced.Store(time.Now().UTC().String())
}

func (c IRCClient) cleanMessage(message string) string {
	return messageClean.ReplaceAllString(message, "")
}
