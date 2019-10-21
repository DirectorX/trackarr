package parser

import (
	"errors"
	"strings"

	"github.com/l3uddz/trackarr/config"

	"github.com/antchfx/xmlquery"
)

/* Private */

func parseServers(t *config.TrackerInfo, doc *xmlquery.Node) error {
	for _, n := range xmlquery.Find(doc, "//servers/server") {
		// parse server hosts
		serverNames := n.SelectAttr("serverNames")
		if serverNames == "" {
			log.Errorf("Failed parsing %q from tracker server: %s", "serverNames", n.OutputXML(true))
			continue
		}

		serverHosts := strings.Split(serverNames, ",")
		if len(serverHosts) < 1 {
			log.Errorf("Failed parsing %q from tracker server: %s", "serverNames", n.OutputXML(true))
			continue
		}
		log.Tracef("Found tracker server hosts: %s", strings.Join(serverHosts, ", "))

		// parse server channels
		channelNames := n.SelectAttr("channelNames")
		if channelNames == "" {
			log.Errorf("Failed parsing %q from tracker server: %s", "channelNames", n.OutputXML(true))
			continue
		}

		serverChannels := strings.Split(channelNames, ",")
		if len(serverChannels) < 1 {
			log.Errorf("Failed parsing %q from tracker server: %s", "channelNames", n.OutputXML(true))
			continue
		}
		log.Tracef("Found tracker server channels: %s", strings.Join(serverChannels, ", "))

		// parse server announcers
		announcerNames := n.SelectAttr("announcerNames")
		if announcerNames == "" {
			log.Errorf("Failed parsing %q from tracker server: %s", "announcerNames", n.OutputXML(true))
			continue
		}

		serverAnnouncers := strings.Split(announcerNames, ",")
		if len(serverAnnouncers) < 1 {
			log.Errorf("Failed parsing %q from tracker server: %s", "announcerNames", n.OutputXML(true))
			continue
		}
		log.Tracef("Found tracker server announcers: %s", strings.Join(serverAnnouncers, ", "))

		// add parsed details to lists
		t.Servers = append(t.Servers, serverHosts...)
		t.Channels = append(t.Channels, serverChannels...)
		t.Announcers = append(t.Announcers, serverAnnouncers...)

	}

	// were servers parsed?
	if len(t.Servers) == 0 {
		return errors.New("failed parsing tracker servers")
	}

	return nil
}
