package ircclient

import irc "github.com/thoj/go-ircevent"

/* Private */

func (c *IRCClient) handleMessage(event *irc.Event) {
	channelName := event.Arguments[0]
	cleanMessage := c.cleanMessage(event.Message())

	c.log.Tracef("%s: %s", channelName, cleanMessage)
}

func (c IRCClient) cleanMessage(message string) string {
	return c.cleanRxp.ReplaceAllString(message, "")
}
