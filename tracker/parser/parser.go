package parser

import (
	"io/ioutil"
	"strings"

	"github.com/l3uddz/trackarr/config"
	"github.com/l3uddz/trackarr/logger"

	"github.com/antchfx/xmlquery"
	"github.com/pkg/errors"
)

var (
	log = logger.GetLogger("tracker_parser")
)

/* Public */

func New(name string) *config.TrackerInfo {
	p := &config.TrackerInfo{}
	p.Name = name

	return p
}

/* Private */

func Parse(t *config.TrackerInfo) error {
	trackerFilePath, err := trackerFilePath(t)
	if err != nil {
		return err
	}

	// read tracker file
	trackerData, err := ioutil.ReadFile(*trackerFilePath)
	if err != nil {
		log.WithError(err).Errorf("Failed reading tracker file: %q", *trackerFilePath)
		return errors.Wrapf(err, "failed reading tracker file: %q", *trackerFilePath)
	}

	// get tracker doc root for xpath queries
	doc, err := xmlquery.Parse(strings.NewReader(string(trackerData)))
	if err != nil {
		log.WithError(err).Errorf("Failed parsing doc root of tracker file: %q", trackerFilePath)
		return errors.Wrap(err, "failed parsing tracker file doc root")
	}

	// parse tracker details
	if err := parseDetails(t, doc); err != nil {
		return err
	}

	// parse tracker settings
	if err := parseSettings(t, doc); err != nil {
		return err
	}
	log.Debugf("Parsed %d tracker settings", len(t.Settings))

	// parse tracker servers
	if err := parseServers(t, doc); err != nil {
		return err
	}
	log.Debugf("Parsed %d tracker servers, %d channels and %d announcers", len(t.Servers),
		len(t.Channels), len(t.Announcers))

	// parse tracker ignore lines
	if err := parseIgnores(t, doc); err != nil {
		return err
	}
	log.Debugf("Parsed %d tracker ignore lines", len(t.IgnoreLines))

	// parse tracker patterns
	if err := parsePatterns(t, doc); err != nil {
		return err
	}
	log.Debugf("Parsed %d tracker linepatterns / %d multilinepatterns",
		len(t.LinePatterns),
		len(t.MultiLinePatterns))

	// parse tracker match rules
	if err := parseRules(t, doc); err != nil {
		return err
	}

	return nil
}
