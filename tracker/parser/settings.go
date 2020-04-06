package parser

import (
	"errors"
	"fmt"
	"strings"

	"gitlab.com/cloudb0x/trackarr/config"
	listutils "gitlab.com/cloudb0x/trackarr/utils/lists"

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

		// is this a textbox setting?
		if strings.EqualFold(settingName, "textbox") {
			if s := n.SelectAttr("name"); s == "" {
				return fmt.Errorf("failed parsing tracker setting: %s", n.OutputXML(true))
			} else {
				settingName = s
			}
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
