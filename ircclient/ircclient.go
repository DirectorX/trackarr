package ircclient

import (
	"fmt"
	"github.com/l3uddz/trackarr/autodl/parser"
	"github.com/l3uddz/trackarr/config"
	"github.com/l3uddz/trackarr/logger"
	"github.com/sirupsen/logrus"
	"github.com/thoj/go-ircevent"
)

var (
	log = logger.GetLogger("irc")
)

/* Struct */

type IRCClient struct {
	/* private */
	conn   *irc.Connection
	cfg    *config.TrackerConfiguration
	parser *parser.Parser
	log    *logrus.Entry
	/* public */
}

/* Public */

func Init(p *parser.Parser, c *config.TrackerConfiguration) (*IRCClient, error) {
	log.Tracef("Initializing IRC client for parser: %s", p.Tracker.LongName)

	// set variables
	logName := p.Tracker.LongName
	if p.Tracker.ShortName != nil {
		logName = *p.Tracker.ShortName
	}

	// initialize irc object and irc client
	conn := irc.IRC(c.IRC.Nickname, c.IRC.Nickname)
	client := &IRCClient{
		conn:   conn,
		cfg:    c,
		parser: p,
		log:    logger.GetLogger(logName),
	}

	// set config precedence
	client.setConfigPrecedence()

	// set callbacks
	conn.AddCallback("001", client.handleConnected)
	conn.AddCallback("366", client.handleJoined)
	conn.AddCallback("PRIVMSG", client.handlePrivMsg)

	return client, nil
}

/* Private */

func (c *IRCClient) setConfigPrecedence() {
	// set server from config
	if c.cfg.IRC.Host != nil && c.cfg.IRC.Port != nil {
		serverString := fmt.Sprintf("%s:%s", *c.cfg.IRC.Host, *c.cfg.IRC.Port)
		c.parser.Tracker.Servers = nil
		c.parser.Tracker.Servers = []string{
			serverString,
		}
	}

	// set channels from config
	if len(c.cfg.IRC.Channels) >= 1 {
		c.parser.Tracker.Channels = nil
		c.parser.Tracker.Channels = c.cfg.IRC.Channels
	}

	// set announcers from config
	if len(c.cfg.IRC.Announcers) >= 1 {
		c.parser.Tracker.Announcers = nil
		c.parser.Tracker.Announcers = c.cfg.IRC.Announcers
	}
}
