package ircclient

import (
	"crypto/tls"
	"fmt"
	"github.com/pkg/errors"
)

/* Public */

func (c *IRCClient) Start() error {
	connString := ""

	// connect to manually configured parser server
	if c.cfg.IRC.Host != nil {
		// set basic connection information
		connString = fmt.Sprintf("%s:%d", *c.cfg.IRC.Host, c.cfg.IRC.Port)

		// enable ssl if required
		if c.cfg.IRC.TLS {
			c.conn.UseTLS = true
			c.conn.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		}

		// handle connection to configured server
		log.Infof("Connecting to %s (tls: %v)", connString, c.conn.UseTLS)
		if err := c.conn.Connect(connString); err != nil {
			return errors.Wrap(err, "failed connecting to manually configured server")
		}

		// start event loop
		go c.conn.Loop()

		return nil
	}

	// try parser file servers

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
