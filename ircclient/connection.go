package ircclient

import (
	"crypto/tls"
	"net"
	"strings"

	"github.com/pkg/errors"
)

/* Public */

func (c *IRCClient) Start() error {

	// iterate servers
	for _, serverString := range c.Tracker.Info.Servers {
		connString := ""
		useSsl := false

		// was a port specified ?
		if !strings.Contains(serverString, ":") {
			// there was no : so assume default port 6667
			connString = serverString + ":6667"
		} else {
			// split port from serverString
			host, port, err := net.SplitHostPort(serverString)
			if err != nil {
				c.log.WithError(err).Errorf("Failed splitting port from: %s", serverString)
				continue
			}

			// was the port specified with ssl, e.g. +6697 ?
			if strings.Contains(port, "+") {
				useSsl = true
				port = strings.Replace(port, "+", "", -1)
			}

			// build connString
			connString = host + ":" + port
		}

		// set connection settings
		if useSsl {
			// enable ssl
			c.Conn.UseTLS = true
			c.Conn.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		} else {
			// disable ssl
			c.Conn.UseTLS = false
			c.Conn.TLSConfig = nil
		}

		if c.Tracker.Config.IRC.Sasl.User != "" && c.Tracker.Config.IRC.Sasl.Pass != "" {
			// enable sasl authentication
			c.Conn.UseSASL = true
			c.Conn.SASLLogin = c.Tracker.Config.IRC.Sasl.User
			c.Conn.SASLPassword = c.Tracker.Config.IRC.Sasl.Pass
		}

		// handle connection to configured server
		c.log.Infof("Connecting to %s (ssl: %v / sasl: %v)", connString, useSsl, c.Conn.UseSASL)
		if err := c.Conn.Connect(connString); err != nil {
			c.log.WithError(err).Errorf("failed connecting to server: %s", connString)
			continue
		}

		// start event loop
		go c.Conn.Loop()

		return nil
	}

	return errors.New("failed connecting to an irc server")
}

func (c *IRCClient) Stop() {
	// TODO: close line queue channel
	if c.Conn.Connected() {
		c.log.Warn("Disconnecting...")
		c.Conn.Quit()
	} else {
		c.log.Warn("Not connected...")
	}
}
