package parser

import (
	"regexp"

	"gitlab.com/cloudb0x/trackarr/config"

	"github.com/antchfx/xmlquery"
)

/* Private */

func parseIgnores(t *config.TrackerInfo, doc *xmlquery.Node) error {
	for _, n := range xmlquery.Find(doc, "//parseinfo/ignore/regex") {
		// parse ignore regex
		ignoreRegex := n.SelectAttr("value")
		if ignoreRegex == "" {
			log.Errorf("Failed parsing %q from tracker ignore: %s", "value", n.OutputXML(true))
			continue
		}

		expected := false
		if n.SelectAttr("expected") != "false" {
			expected = true
		}

		log.Tracef("Found tracker ignore (expected: %v): %s", expected, ignoreRegex)

		// compile regex
		rxp, err := regexp.Compile(ignoreRegex)
		if err != nil {
			log.WithError(err).Errorf("Failed compiling tracker ignore: %s", ignoreRegex)
			continue
		}

		// add regex to list
		t.IgnoreLines = append(t.IgnoreLines, config.TrackerIgnore{
			Rxp:      rxp,
			Expected: expected,
		})
	}
	return nil
}
