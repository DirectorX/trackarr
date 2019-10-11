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

type Parser struct {
	/* private */
	trackerName     string
	trackerFilePath string

	/* public */
	Tracker TrackerInfo
}

type TrackerInfo struct {
	Settings          []string
	Servers           []TrackerServer
	IgnoreLines       []TrackerIgnore
	LinePatterns      []TrackerPattern
	MultiLinePatterns []TrackerPattern
}

/* Public */

func Init(tracker string, trackersPath string) (*Parser, error) {
	// validate tracker file exists
	trackerFilePath := filepath.Join(trackersPath, tracker+".tracker")
	if _, err := os.Stat(trackerFilePath); os.IsNotExist(err) {
		log.WithError(err).Errorf("Failed initializing parser for tracker: %q", tracker)
		return nil, errors.Wrapf(err, "failed to initialize parser: %q", trackerFilePath)
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

	// parse tracker settings
	if err := parseTrackerSettings(doc, &trackerInfo); err != nil {
		return nil, err
	}

	// parse tracker servers
	if err := parseTrackerServers(doc, &trackerInfo); err != nil {
		return nil, err
	}

	// parse tracker ignore lines
	if err := parseTrackerIgnores(doc, &trackerInfo); err != nil {
		return nil, err
	}

	// parse tracker patterns
	if err := parseTrackerPatterns(doc, &trackerInfo); err != nil {
		return nil, err
	}

	return &Parser{
		trackerName:     tracker,
		trackerFilePath: trackerFilePath,
		Tracker:         trackerInfo,
	}, nil
}

/* Private */
