package ircclient

import (
	"fmt"
	irc "github.com/thoj/go-ircevent"
	"math/rand"
	"time"
)

/* Private */

func (c *IRCClient) handleNickInUse(event *irc.Event) {
	// generate new nick
	rand.Seed(time.Now().UnixNano())
	newNick := fmt.Sprintf("%s%d", c.cfg.IRC.Nickname, rand.Intn(50))
	c.log.Warnf("Nick in use: %s, changing to: %s", c.conn.GetNick(), newNick)

	// change nick
	c.conn.Nick(newNick)
}
