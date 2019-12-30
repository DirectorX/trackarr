package ircclient

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"time"

	"gitlab.com/cloudb0x/trackarr/autodl/processor"
	"gitlab.com/cloudb0x/trackarr/config"
	"gitlab.com/cloudb0x/trackarr/logger"

	"github.com/sirupsen/logrus"
	irc "github.com/thoj/go-ircevent"
)

var (
	log = logger.GetLogger("irc")

	messageClean *regexp.Regexp
)

/* Const */
const regexMessageClean = `\x0f|\x1f|\x02|\x03(?:[\d]{1,2}(?:,[\d]{1,2})?)?`

/* Struct */

type IRCClient struct {
	Conn      *irc.Connection
	Tracker   *config.TrackerInstance
	Processor *processor.Processor
	// Private
	log *logrus.Entry
}

/* init */

func init() {
	var err error
	if messageClean, err = regexp.Compile(regexMessageClean); err != nil {
		log.WithError(err).Fatal("Failed compiling message clean regex??")
	}
}

/* Public */

func New(t *config.TrackerInstance) (*IRCClient, error) {
	// set variables
	logName := t.Info.LongName
	if t.Info.ShortName != nil {
		logName = *t.Info.ShortName
	}

	// initialize irc object
	conn := irc.IRC(t.Config.IRC.Nickname, t.Config.IRC.Nickname)

	// set base irc object settings
	ircLogger := logger.GetLogger(logName)
	if t.Config.IRC.Verbose {
		conn.Debug = true
		conn.Log.SetOutput(ircLogger.WriterLevel(logrus.TraceLevel))
	} else {
		conn.Log.SetOutput(ioutil.Discard)
	}

	conn.PingFreq = 3 * time.Minute
	conn.Timeout = 15 * time.Second
	conn.Version = "trackarr " + config.Build.Version

	// initialize irc client
	client := &IRCClient{
		Conn:      conn,
		Tracker:   t,
		Processor: processor.New(ircLogger, t),
		// Private
		log: ircLogger,
	}
	// set IRC connection in the Tracker struct
	client.Tracker.IRC = conn

	// set config precedence
	client.setConfigPrecedence()

	// set callbacks
	// - connected
	conn.AddCallback("001", client.handleConnected)
	// - mode
	conn.AddCallback("MODE", client.handleMode)
	// - join
	conn.AddCallback("366", client.handleJoined)
	conn.AddCallback("448", client.handleJoinFailure)
	conn.AddCallback("475", client.handleJoinFailure)
	conn.AddCallback("477", client.handleJoinFailure)
	// - parted
	conn.AddCallback("PART", client.handleParted)
	// - nick change
	conn.ClearCallback("433")
	conn.AddCallback("433", client.handleNickInUse)
	// - message
	conn.AddCallback("PRIVMSG", client.handleMessage)

	return client, nil
}

/* Private */

func (c *IRCClient) setConfigPrecedence() {
	// set server from config
	if c.Tracker.Config.IRC.Host != nil && c.Tracker.Config.IRC.Port != nil {
		c.log.Debugf("Using host and port from config: %s:%s", *c.Tracker.Config.IRC.Host, *c.Tracker.Config.IRC.Port)
		c.Tracker.Info.Servers = []string{
			fmt.Sprintf("%s:%s", *c.Tracker.Config.IRC.Host, *c.Tracker.Config.IRC.Port),
		}
	}

	// set channels from config
	if len(c.Tracker.Config.IRC.Channels) >= 1 {
		c.log.Debugf("Using channels from config: %s", strings.Join(c.Tracker.Config.IRC.Channels, ", "))
		c.Tracker.Info.Channels = nil
		c.Tracker.Info.Channels = c.Tracker.Config.IRC.Channels
	}

	// set announcers from config
	if len(c.Tracker.Config.IRC.Announcers) >= 1 {
		c.log.Debugf("Using announcers from config: %s", strings.Join(c.Tracker.Config.IRC.Announcers, ", "))
		c.Tracker.Info.Announcers = nil
		c.Tracker.Info.Announcers = c.Tracker.Config.IRC.Announcers
	}
}
