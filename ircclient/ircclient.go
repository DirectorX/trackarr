package ircclient

import (
	"fmt"
	"github.com/l3uddz/trackarr/autodl/parser"
	"github.com/l3uddz/trackarr/autodl/processor"
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
	conn      *irc.Connection
	cfg       *config.TrackerConfiguration
	tracker   *parser.TrackerInfo
	log       *logrus.Entry
	cleanRxp  *regexp.Regexp
	processor *processor.Processor
	/* public */
}

/* Public */

func Init(t *parser.TrackerInfo, c *config.TrackerConfiguration) (*IRCClient, error) {
	// set variables
	logName := t.LongName
	if t.ShortName != nil {
		logName = *t.ShortName
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
		conn.Log.SetOutput(ircLogger.WriterLevel(logrus.TraceLevel))
	} else {
		conn.Log.SetOutput(ioutil.Discard)
	}

	conn.PingFreq = 3 * time.Minute
	conn.Version = "trackarr " + config.Version

	// initialize irc client
	client := &IRCClient{
		conn:      conn,
		cfg:       c,
		tracker:   t,
		log:       ircLogger,
		cleanRxp:  cleanRxp,
		processor: processor.New(ircLogger, t, c),
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
		c.log.Debugf("Using host and port from config: %s:%s", *c.cfg.IRC.Host, *c.cfg.IRC.Port)
		serverString := fmt.Sprintf("%s:%s", *c.cfg.IRC.Host, *c.cfg.IRC.Port)
		c.tracker.Servers = nil
		c.tracker.Servers = []string{
			serverString,
		}
	}

	// set channels from config
	if len(c.cfg.IRC.Channels) >= 1 {
		c.log.Debugf("Using channels from config: %s", strings.Join(c.cfg.IRC.Channels, ", "))
		c.tracker.Channels = nil
		c.tracker.Channels = c.cfg.IRC.Channels
	}

	// set announcers from config
	if len(c.cfg.IRC.Announcers) >= 1 {
		c.log.Debugf("Using announcers from config: %s", strings.Join(c.cfg.IRC.Announcers, ", "))
		c.tracker.Announcers = nil
		c.tracker.Announcers = c.cfg.IRC.Announcers
	}
}
