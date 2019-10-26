package processor

import (
	"github.com/l3uddz/trackarr/config"
	"github.com/l3uddz/trackarr/release"
	"github.com/l3uddz/trackarr/utils/maps"
)

/* Private */

func (p *Processor) processQueue(queue chan string) {
	var patterns []config.TrackerPattern

	// set patterns
	if len(p.Tracker.Info.LinePatterns) > 0 {
		patterns = p.Tracker.Info.LinePatterns
	} else if len(p.Tracker.Info.MultiLinePatterns) > 0 {
		patterns = p.Tracker.Info.MultiLinePatterns
	} else {
		p.Log.Fatalf("Failed determining pattern type for processor...")
		return
	}

	// process lines
	for {
		vars := map[string]string{}
		parseFailed := false

		// iterate each pattern finding a match
		for _, pattern := range patterns {
			line, err := p.nextGoodLine(queue)
			if err != nil {
				p.Log.WithError(err).Errorf("Failed dequeuing line to process, discarding release...")
				parseFailed = true
				break
			}

			// process line
			p.Log.Debugf("Processing line: %s", line)
			patternVars, err := p.matchPattern(&pattern, line)
			if err != nil {
				p.Log.WithError(err).Errorf("Failed matching pattern, discarding release...")
				parseFailed = true
				break
			}

			// update vars
			maps.MergeStringMap(vars, patternVars)
		}

		if parseFailed {
			// pattern parsing above failed to parse a complete release
			continue
		}

		// finished parsing release lines - process rules
		if err := p.processRules(p.Tracker.Info.LineMatchedRules, vars); err != nil {
			p.Log.WithError(err).Errorf("Failed processing release lines due to rules failure...")
			continue
		}

		p.Log.Debugf("Finished processing release lines, release vars: %+v", vars)

		// convert parsed release vars to release struct and begin release processing
		if trackerRelease, err := release.FromMap(p.Tracker, p.Log, vars); err != nil {
			p.Log.WithError(err).Errorf("Failed converting release vars to a release struct...")
		} else {
			// start processing this release
			go func(tr *release.Release) {
				tr.Process()
			}(trackerRelease)
		}
	}
}

func (p *Processor) nextGoodLine(queue chan string) (string, error) {
	for {
		// pop line from queue
		queuedLine := <-queue

		// should ignore this line?
		if p.shouldIgnoreLine(queuedLine) {
			continue
		}

		return queuedLine, nil
	}
}
