package ircclient

/* Public */

func (c *IRCClient) Start() {
	c.log.Info("Starting connection...")
}

func (c *IRCClient) Stop() {
	c.log.Warn("Stopping connection")
}
