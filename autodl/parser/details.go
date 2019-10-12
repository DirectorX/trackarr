package parser

import (
	"github.com/antchfx/xmlquery"
	"github.com/pkg/errors"
)

/* Private */

func parseTrackerDetails(doc *xmlquery.Node, tracker *TrackerInfo, trackerName string) error {
	// find trackerinfo element
	trackerInfo := xmlquery.FindOne(doc, "/trackerinfo")
	if trackerInfo == nil {
		return errors.New("failed parsing trackerinfo")
	}

	// parse details
	shortName := trackerInfo.SelectAttr("shortName")
	longName := trackerInfo.SelectAttr("longName")

	if shortName == "" {
		log.Debugf("Failed to parse tracker %q from: %s", "shortName", trackerInfo.OutputXML(true))
	} else {
		tracker.ShortName = &shortName
		log.Debugf("Found tracker short name: %s", *tracker.ShortName)
	}

	if longName == "" {
		log.Debugf("Failed to parse tracker %q from: %s", "longName", trackerInfo.OutputXML(true))
		tracker.LongName = trackerName
	} else {
		tracker.LongName = longName
		log.Debugf("Found tracker long name: %s", tracker.LongName)
	}

	return nil
}
