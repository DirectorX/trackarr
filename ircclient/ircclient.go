package ircclient

import (
	"fmt"
	"github.com/l3uddz/trackarr/autodl/parser"
	"github.com/l3uddz/trackarr/config"
	"github.com/l3uddz/trackarr/logger"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/thoj/go-ircevent"
	"io/ioutil"
	"regexp"
	"strings"
	"time"
)

var (
	log = logger.GetLogger("irc")
)

/* Const */
const RegexMessageClean = `\x0f|\x1f|\x02|\x03(?:[\d]{1,2}(?:,[\d]{1,2})?)?`

/* Struct */

type IRCClient struct {
	/* private */
	conn     *irc.Connection
	cfg      *config.TrackerConfiguration
	parser   *parser.Parser
	log      *logrus.Entry
	cleanRxp *regexp.Regexp
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

	cleanRxp, err := regexp.Compile(RegexMessageClean)
	if err != nil {
		log.WithError(err).Errorf("Failed compiling message clean regex??")
		return nil, errors.Wrap(err, "failed compiling message clean regex")
	}

	// initialize irc object
	conn := irc.IRC(c.IRC.Nickname, c.IRC.Nickname)

	// set base irc object settings
	ircLogger := logger.GetLogger(logName)
	if c.Verbose {
		conn.Debug = true
		conn.Log.SetOutput(ircLogger.Writer())
	} else {
		conn.Log.SetOutput(ioutil.Discard)
	}

	conn.PingFreq = 3 * time.Minute

	// initialize irc client
	client := &IRCClient{
		conn:     conn,
		cfg:      c,
		parser:   p,
		log:      ircLogger,
		cleanRxp: cleanRxp,
	}

	// set config precedence
	client.setConfigPrecedence()

	// set callbacks
	conn.AddCallback("001", client.handleConnected)
	conn.AddCallback("366", client.handleJoined)
	conn.AddCallback("PRIVMSG", client.handleMessage)

	return client, nil
}

/* Private */

func (c *IRCClient) setConfigPrecedence() {
	// set server from config
	if c.cfg.IRC.Host != nil && c.cfg.IRC.Port != nil {
		log.Debugf("Using host and port from tracker config: %s:%s", *c.cfg.IRC.Host, *c.cfg.IRC.Port)
		serverString := fmt.Sprintf("%s:%s", *c.cfg.IRC.Host, *c.cfg.IRC.Port)
		c.parser.Tracker.Servers = nil
		c.parser.Tracker.Servers = []string{
			serverString,
		}
	}

	// set channels from config
	if len(c.cfg.IRC.Channels) >= 1 {
		log.Debugf("Using channels from tracker config: %s", strings.Join(c.cfg.IRC.Channels, ", "))
		c.parser.Tracker.Channels = nil
		c.parser.Tracker.Channels = c.cfg.IRC.Channels
	}

	// set announcers from config
	if len(c.cfg.IRC.Announcers) >= 1 {
		log.Debugf("Using announcers from tracker config: %s", strings.Join(c.cfg.IRC.Announcers, ", "))
		c.parser.Tracker.Announcers = nil
		c.parser.Tracker.Announcers = c.cfg.IRC.Announcers
	}
}
