package ircclient

import (
	"crypto/tls"
	"github.com/pkg/errors"
	"net"
	"strings"
)

/* Public */

func (c *IRCClient) Start() error {

	// iterate servers
	for _, serverString := range c.parser.Tracker.Servers {
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
				log.WithError(err).Errorf("Failed splitting port from: %s", serverString)
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
			c.conn.UseTLS = true
			c.conn.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		} else {
			// disable ssl
			c.conn.UseTLS = false
			c.conn.TLSConfig = nil
		}

		// handle connection to configured server
		log.Infof("Connecting to %s (ssl: %v)", connString, useSsl)
		if err := c.conn.Connect(connString); err != nil {
			log.WithError(err).Errorf("failed connecting to server: %s", connString)
			continue
		}

		// start event loop
		go c.conn.Loop()

		return nil
	}

	return errors.New("failed connecting to an irc server")
}

func (c *IRCClient) Stop() {
	if c.conn.Connected() {
		c.log.Warn("Disconnecting...")
		c.conn.Quit()
	} else {
		c.log.Warn("Not connected...")
	}
}
