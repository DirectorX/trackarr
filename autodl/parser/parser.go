package parser

import (
	"encoding/xml"
	"github.com/antchfx/xmlquery"
	"github.com/l3uddz/trackarr/logger"
	listutils "github.com/l3uddz/trackarr/utils/lists"
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
	/* privates */
	trackerName     string
	trackerFilePath string

	/* public */
	Tracker *TrackerInfo
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

	// decode tracker file
	trackerInfo := &TrackerInfo{}
	if err := xml.Unmarshal([]byte(trackerData), &trackerInfo); err != nil {
		log.WithError(err).Errorf("Failed decoding tracker file: %q", trackerFilePath)
		return nil, errors.Wrapf(err, "failed decoding tracker file: %q", trackerFilePath)
	}

	// get tracker doc root for xpath queries
	doc, err := xmlquery.Parse(strings.NewReader(string(trackerData)))
	if err != nil {
		log.WithError(err).Errorf("Failed parsing doc root of tracker file: %q", trackerFilePath)
		return nil, errors.Wrap(err, "failed parsing tracker file doc root")
	}

	// parse required tracker settings
	if err := parseTrackerSettings(doc, trackerInfo); err != nil {
		return nil, err
	}

	return &Parser{
		trackerName:     tracker,
		trackerFilePath: trackerFilePath,
		Tracker:         trackerInfo,
	}, nil
}

/* Private */

func parseTrackerSettings(doc *xmlquery.Node, tracker *TrackerInfo) error {
	skipSettings := []string{
		"description",
	}

	for _, n := range xmlquery.Find(doc, "//settings/*[name()]") {
		// strip gazelle_ prefix
		settingName := strings.Replace(n.Data, "gazelle_", "", -1)

		// skip specific settings
		if listutils.StringListContains(skipSettings, settingName, true) {
			log.Debugf("Skipping tracker setting: %q", settingName)
			continue
		}

		log.Debugf("Found tracker setting: %q", settingName)

		// add setting to list
		tracker.Settings = append(tracker.Settings, settingName)
	}

	// were settings parsed?
	if len(tracker.Settings) == 0 {
		return errors.New("failed parsing tracker settings")
	}

	return nil
}
