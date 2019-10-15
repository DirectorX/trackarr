package processor

import (
	"fmt"
	"github.com/l3uddz/trackarr/autodl/parser"
	stringutils "github.com/l3uddz/trackarr/utils/strings"
	"strings"
)

/* Private */

func (p *Processor) matchPattern(pattern *parser.TrackerPattern, line string) (map[string]string, error) {
	results := map[string]string{}
	matches := pattern.Rxp.FindStringSubmatch(line)
	if len(matches) != (len(pattern.Vars) + 1) {
		// pattern did not match the line
		return nil, fmt.Errorf("pattern %q did not match: %s", pattern.Rxp, line)
	}

	// pattern matched - extract vars
	matchPos := 1
	matchCount := len(matches)
	p.log.Tracef("Pattern %q matched with %d groups", pattern.Rxp, matchCount)

	for _, varName := range pattern.Vars {
		// this should not occur - but we must ensure we dont try and access an out of bounds index
		if matchPos > matchCount {
			p.log.Warnf("Failed parsing pattern var %q from match group %d", varName, matchPos)
			continue
		}

		// add var to map
		results[varName] = strings.TrimSpace(matches[matchPos])
		matchPos++
	}

	p.log.Debugf("Found match: %+v", results)
	return results, nil
}

func (p *Processor) matchPatterns(patterns *[]parser.TrackerPattern, line string) map[string]string {
	results := map[string]string{}

	// iterate patterns matching against line
	for _, pattern := range *patterns {
		matches := pattern.Rxp.FindStringSubmatch(line)
		if len(matches) != (len(pattern.Vars) + 1) {
			// pattern did not match the line
			p.log.Tracef("Pattern %q did not match: %s", pattern.Rxp, line)
			continue
		}

		// pattern matched - extract vars
		matchPos := 1
		matchCount := len(matches)
		p.log.Tracef("Pattern %q matched with %d groups", pattern.Rxp, matchCount)

		for _, varName := range pattern.Vars {
			// this should not occur - but we must ensure we dont try and access an out of bounds index
			if matchPos > matchCount {
				p.log.Warnf("Failed parsing pattern var %q from match group %d", varName, matchPos)
				continue
			}

			// add var to map
			results[varName] = strings.TrimSpace(matches[matchPos])
			matchPos++
		}

		p.log.Debugf("Found match: %s", stringutils.JsonifyLax(results))
		break
	}

	return results
}
