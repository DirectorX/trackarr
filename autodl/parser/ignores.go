package parser

import (
	"github.com/antchfx/xmlquery"
	"regexp"
)

/* Struct */
type TrackerIgnore struct {
	Rxp      *regexp.Regexp
	Expected bool
}

/* Private */

func parseTrackerIgnores(doc *xmlquery.Node, tracker *TrackerInfo) error {
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

		log.Debugf("Found tracker ignore (expected: %v): %s", expected, ignoreRegex)

		// compile regex
		rxp, err := regexp.Compile(ignoreRegex)
		if err != nil {
			log.WithError(err).Errorf("Failed compiling tracker ignore: %s", ignoreRegex)
			continue
		}

		// add regex to list
		tracker.IgnoreLines = append(tracker.IgnoreLines, TrackerIgnore{
			Rxp:      rxp,
			Expected: expected,
		})
	}
	return nil
}
