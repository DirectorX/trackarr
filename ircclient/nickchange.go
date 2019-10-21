package ircclient

import (
	"fmt"
	"math/rand"
	"time"

	irc "github.com/thoj/go-ircevent"
)

/* Private */

func (c *IRCClient) handleNickInUse(event *irc.Event) {
	// generate new nick
	rand.Seed(time.Now().UnixNano())
	newNick := fmt.Sprintf("%s%d", c.Tracker.Config.IRC.Nickname, rand.Intn(50))
	c.log.Warnf("Nick in use: %s, changing to: %s", c.Conn.GetNick(), newNick)

	// change nick
	c.Conn.Nick(newNick)
}
