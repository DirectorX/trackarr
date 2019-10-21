package parser

import (
	"errors"
	"strings"

	listutils "github.com/l3uddz/trackarr/utils/lists"
	"github.com/l3uddz/trackarr/config"

	"github.com/antchfx/xmlquery"
)

/* Private */

func parseSettings(t *config.TrackerInfo, doc *xmlquery.Node) error {
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
		t.Settings = append(t.Settings, settingName)
	}

	// were settings parsed?
	if len(t.Settings) == 0 {
		return errors.New("failed parsing tracker settings")
	}

	return nil
}
