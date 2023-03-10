package processor

import (
	"fmt"
	"strings"

	"gitlab.com/cloudb0x/trackarr/config"
)

/* Private */

func (p *Processor) matchPattern(pattern *config.TrackerPattern, line string) (map[string]string, error) {
	results := map[string]string{}
	matches := pattern.Rxp.FindStringSubmatch(line)
	if len(matches) != (len(pattern.Vars) + 1) {
		// pattern did not match the line
		return nil, fmt.Errorf("pattern %q did not match: %s", pattern.Rxp, line)
	}

	// pattern matched - extract vars
	matchPos := 1
	matchCount := len(matches)
	p.Log.Tracef("Pattern %q matched with %d groups", pattern.Rxp, matchCount)

	for _, varName := range pattern.Vars {
		// this should not occur - but we must ensure we dont try and access an out of bounds index
		if matchPos > matchCount {
			p.Log.Warnf("Failed parsing pattern var %q from match group %d", varName, matchPos)
			continue
		}

		// add var to map
		results[varName] = strings.TrimSpace(matches[matchPos])
		matchPos++
	}

	p.Log.Tracef("Found match: %+v", results)
	return results, nil
}
