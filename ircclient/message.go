package ircclient

import irc "github.com/thoj/go-ircevent"

/* Private */

func (c *IRCClient) handlePrivMsg(event *irc.Event) {
	channelName := event.Arguments[0]

	c.log.Infof("Private message from %s: %s", channelName, event.Message())
}
