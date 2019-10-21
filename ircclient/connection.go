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
		if strings.Index(serverString, ":") == -1 {
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
			if strings.Index(port, "+") != -1 {
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

		// handle connection to configured server
		c.log.Infof("Connecting to %s (ssl: %v)", connString, useSsl)
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
	if c.Conn.Connected() {
		c.log.Warn("Disconnecting...")
		c.Conn.Quit()
	} else {
		c.log.Warn("Not connected...")
	}
}
