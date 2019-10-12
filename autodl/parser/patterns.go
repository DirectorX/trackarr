package parser

import (
	"github.com/antchfx/xmlquery"
	"github.com/pkg/errors"
	"regexp"
)

/* Enum */

type MessagePatternType int

const (
	LinePattern MessagePatternType = iota + 1
	MultiLinePattern
)

/* Struct */

type TrackerPattern struct {
	PatternType MessagePatternType
	Rxp         *regexp.Regexp
	Vars        []string
	Optional    bool
}

/* Private */

func parseTrackerPatterns(doc *xmlquery.Node, tracker *TrackerInfo) error {
	// parse line patterns
	if err := parsePatterns(xmlquery.Find(doc, "//parseinfo/linepatterns/extract"), LinePattern, tracker);
		err != nil {
		log.WithError(err).Errorf("Failed parsing tracker linepatterns")
		return errors.Wrap(err, "failed to parse tracker line patterns")
	}

	// parse multiline patterns
	if err :=
		parsePatterns(xmlquery.Find(doc, "//parseinfo/multilinepatterns/extract"), MultiLinePattern, tracker);
		err != nil {
		return errors.Wrap(err, "failed to parse tracker multiline patterns")
	}

	return nil
}

func parsePatterns(nodes []*xmlquery.Node, patternType MessagePatternType, tracker *TrackerInfo) error {
	patternTypeString := "linepattern"

	if patternType == MultiLinePattern {
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
			case LinePattern:
				tracker.LinePatterns = append(tracker.LinePatterns, TrackerPattern{
					PatternType: LinePattern,
					Rxp:         rxp,
					Vars:        lineVars,
					Optional:    optional,
				})
			default:
				tracker.MultiLinePatterns = append(tracker.MultiLinePatterns, TrackerPattern{
					PatternType: MultiLinePattern,
					Rxp:         rxp,
					Vars:        lineVars,
					Optional:    optional,
				})
			}
		}
	}

	return nil
}
