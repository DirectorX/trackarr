package parser

import (
	"github.com/l3uddz/trackarr/config"

	"github.com/antchfx/xmlquery"
	"github.com/pkg/errors"
)

/* Private */

func parseDetails(t *config.TrackerInfo, doc *xmlquery.Node) error {
	// find trackerinfo element
	trackerInfo := xmlquery.FindOne(doc, "/trackerinfo")
	if trackerInfo == nil {
		return errors.New("failed parsing trackerinfo")
	}

	// parse details
	shortName := trackerInfo.SelectAttr("shortName")
	longName := trackerInfo.SelectAttr("longName")

	switch shortName {
	case "":
		log.Warnf("Failed to parse tracker %q from: %s", "shortName", trackerInfo.OutputXML(true))
	default:
		t.ShortName = &shortName
		log.Tracef("Found tracker short name: %s", *t.ShortName)
	}

	switch longName {
	case "":
		log.Warnf("Failed to parse tracker %q from: %s", "longName", trackerInfo.OutputXML(true))
		t.LongName = t.Name
	default:
		t.LongName = longName
		log.Tracef("Found tracker long name: %s", t.LongName)
	}

	return nil
}
