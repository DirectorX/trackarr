package parser

import (
	"regexp"

	"gitlab.com/cloudb0x/trackarr/config"

	"github.com/antchfx/xmlquery"
	"github.com/pkg/errors"
)

/* Private */

func parsePatterns(t *config.TrackerInfo, doc *xmlquery.Node) error {
	// parse line patterns
	if err := parsePattern(t, xmlquery.Find(doc, "//parseinfo/linepatterns/extract"), config.LinePattern); err != nil {
		log.WithError(err).Errorf("Failed parsing tracker linepatterns")
		return errors.Wrap(err, "failed to parse tracker line patterns")
	}

	// parse multiline patterns
	if err :=
		parsePattern(t, xmlquery.Find(doc, "//parseinfo/multilinepatterns/extract"), config.MultiLinePattern); err != nil {
		return errors.Wrap(err, "failed to parse tracker multiline patterns")
	}

	return nil
}

func parsePattern(t *config.TrackerInfo, nodes []*xmlquery.Node, patternType config.MessagePatternType) error {
	patternTypeString := "linepattern"

	if patternType == config.MultiLinePattern {
		patternTypeString = "multilinepattern"
	}

	for _, n := range nodes {
		// parse pattern regex
		regexNode := n.SelectElement("regex")
		lineRegex := regexNode.SelectAttr("value")
		if lineRegex == "" {
			log.Errorf("Failed parsing %q from tracker %s: %s", "value", patternTypeString,
				n.OutputXML(true))
			continue
		}

		optional := false
		if n.SelectAttr("optional") == "true" {
			optional = true
		}

		rxp, err := regexp.Compile(lineRegex)
		if err != nil {
			log.WithError(err).Errorf("Failed compiling tracker %s: %s", patternTypeString, lineRegex)
			continue
		}

		log.Tracef("Found tracker %s (optional: %v): %s", patternTypeString, optional, lineRegex)

		// parse pattern vars
		var lineVars []string
		varsNode := n.SelectElement("vars")
		for _, v := range varsNode.SelectElements("var") {
			varName := v.SelectAttr("name")
			if varName == "" {
				log.Errorf("Failed parsing tracker %s var from: %s", patternTypeString, v.OutputXML(true))
				continue
			}
			log.Tracef("Found tracker %s var: %s", patternTypeString, varName)
			lineVars = append(lineVars, varName)
		}

		// add to list
		if len(lineVars) > 0 {
			switch patternType {
			case config.LinePattern:
				t.LinePatterns = append(t.LinePatterns, config.TrackerPattern{
					PatternType: config.LinePattern,
					Rxp:         rxp,
					Vars:        lineVars,
					Optional:    optional,
				})
			default:
				t.MultiLinePatterns = append(t.MultiLinePatterns, config.TrackerPattern{
					PatternType: config.MultiLinePattern,
					Rxp:         rxp,
					Vars:        lineVars,
					Optional:    optional,
				})
			}
		}
	}

	return nil
}
