package parser

import (
	"errors"
	"github.com/antchfx/xmlquery"
	listutils "github.com/l3uddz/trackarr/utils/lists"
	"strings"
)

/* Private */

func parseTrackerSettings(doc *xmlquery.Node, tracker *TrackerInfo) error {
	skipSettings := []string{
		"description",
		"cookie_description",
	}

	for _, n := range xmlquery.Find(doc, "//settings/*[name()]") {
		// strip gazelle_ prefix
		settingName := strings.Replace(n.Data, "gazelle_", "", -1)

		// skip specific settings
		if listutils.StringListContains(skipSettings, settingName, true) {
			log.Tracef("Skipping tracker setting: %q", settingName)
			continue
		}

		log.Tracef("Found tracker setting: %q", settingName)

		// add setting to list
		tracker.Settings = append(tracker.Settings, settingName)
	}

	// were settings parsed?
	if len(tracker.Settings) == 0 {
		return errors.New("failed parsing tracker settings")
	}

	return nil
}
