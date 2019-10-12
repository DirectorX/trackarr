package parser

import (
	"github.com/antchfx/xmlquery"
	"github.com/l3uddz/trackarr/logger"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	log = logger.GetLogger("autodl")
)

/* Struct */

type TrackerInfo struct {
	LongName   string
	ShortName  *string
	Settings   []string
	Servers    []string
	Channels   []string
	Announcers []string

	IgnoreLines       []TrackerIgnore
	LinePatterns      []TrackerPattern
	MultiLinePatterns []TrackerPattern

	LineMatchedRules *xmlquery.Node
}

/* Public */

func Parse(tracker string, trackersPath string) (*TrackerInfo, error) {
	// validate tracker file exists
	trackerFilePath := filepath.Join(trackersPath, tracker+".tracker")
	if _, err := os.Stat(trackerFilePath); os.IsNotExist(err) {
		log.WithError(err).Errorf("Failed locating tracker file: %q", tracker)
		return nil, errors.Wrapf(err, "failed locating tracker file: %q", trackerFilePath)
	}

	// read tracker file
	trackerData, err := ioutil.ReadFile(trackerFilePath)
	if err != nil {
		log.WithError(err).Errorf("Failed reading tracker file: %q", trackerFilePath)
		return nil, errors.Wrapf(err, "failed reading tracker file: %q", trackerFilePath)
	}

	// parse tracker info
	trackerInfo := TrackerInfo{}

	// get tracker doc root for xpath queries
	doc, err := xmlquery.Parse(strings.NewReader(string(trackerData)))
	if err != nil {
		log.WithError(err).Errorf("Failed parsing doc root of tracker file: %q", trackerFilePath)
		return nil, errors.Wrap(err, "failed parsing tracker file doc root")
	}

	// parse tracker details
	if err := parseTrackerDetails(doc, &trackerInfo, tracker); err != nil {
		return nil, err
	}

	// parse tracker settings
	if err := parseTrackerSettings(doc, &trackerInfo); err != nil {
		return nil, err
	}
	log.Debugf("Parsed %d tracker settings", len(trackerInfo.Settings))

	// parse tracker servers
	if err := parseTrackerServers(doc, &trackerInfo); err != nil {
		return nil, err
	}
	log.Debugf("Parsed %d tracker servers, %d channels and %d announcers", len(trackerInfo.Servers),
		len(trackerInfo.Channels), len(trackerInfo.Announcers))

	// parse tracker ignore lines
	if err := parseTrackerIgnores(doc, &trackerInfo); err != nil {
		return nil, err
	}
	log.Debugf("Parsed %d tracker ignore lines", len(trackerInfo.IgnoreLines))

	// parse tracker patterns
	if err := parseTrackerPatterns(doc, &trackerInfo); err != nil {
		return nil, err
	}
	log.Debugf("Parsed %d tracker linepatterns / %d multilinepatterns",
		len(trackerInfo.LinePatterns),
		len(trackerInfo.MultiLinePatterns))

	// parse tracker match rules
	if err := parseTrackerRules(doc, &trackerInfo); err != nil {
		return nil, err
	}

	return &trackerInfo, nil
}

/* Private */
