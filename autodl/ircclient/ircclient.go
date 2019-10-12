package ircclient

import (
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
	conn    *irc.Connection
	cfg     *config.TrackerConfiguration
	tracker *parser.Parser
	log     *logrus.Entry
	/* public */
}

/* Public */

func Init(p *parser.Parser, c *config.TrackerConfiguration) (*IRCClient, error) {
	log.Tracef("Initializing IRC client for tracker: %s", p.Tracker.LongName)

	// set variables
	logName := p.Tracker.LongName
	if p.Tracker.ShortName != nil {
		logName = *p.Tracker.ShortName
	}

	// initialize irc object and ircclient
	conn := irc.IRC("thebigmuncho", "thebigmuncho")
	client := &IRCClient{
		conn:    conn,
		cfg:     c,
		tracker: p,
		log:     logger.GetLogger(logName),
	}

	// set callbacks
	conn.AddCallback("PRIVMSG", client.handlePrivMsg)

	return client, nil
}



/* Private */
