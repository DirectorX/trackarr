package processor

import (
	"gitlab.com/cloudb0x/trackarr/config"
	"gitlab.com/cloudb0x/trackarr/release"
	"gitlab.com/cloudb0x/trackarr/utils/maps"

	"github.com/pkg/errors"
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
		p.Log.Fatal("Failed determining pattern type for processor...")
		return
	}

	patternsSize := len(patterns)

	// process lines
	for {
		vars := map[string]string{}
		parseFailed := false
		linePatternParsed := false

		// iterate each pattern finding a match
		for pos, pattern := range patterns {
			line, err := p.nextGoodLine(queue)
			if err != nil {
				// if an error occurred, the only possible cause is due to the channel being closed
				p.Log.WithError(err).Error("Failed dequeuing line to process, processor shutting down...")
				return
			}

			// process line
			p.Log.Debugf("Processing line: %s", line)
			patternVars, err := p.matchPattern(&pattern, line)
			if err != nil {
				parseFailed = true

				// try next pattern (for LinePattern type only)
				if pattern.PatternType == config.LinePattern && (pos+1) < patternsSize {
					// try the next pattern
					p.Log.WithError(err).Trace("Failed matching pattern, trying next...")
					continue
				}

				// multi-line pattern failed or all line patterns failed.
				p.Log.WithError(err).Error("Failed matching pattern, discarding release...")
				break

			} else {
				// update vars
				maps.MergeStringMap(vars, patternVars)

				// flag that a LinePattern was parsed
				if pattern.PatternType == config.LinePattern {
					linePatternParsed = true
					// line patterns need only one match before proceeding with processing rules
					break
				}

				// multi-line patterns must continue until all patterns matched before processing rules
			}
		}

		if parseFailed && !linePatternParsed {
			// pattern parsing failed to parse a complete release
			continue
		}

		// finished parsing release lines - process rules
		if err := p.processRules(p.Tracker.Info.LineMatchedRules, vars); err != nil {
			p.Log.WithError(err).Error("Failed processing release lines due to rules failure...")
			continue
		}

		p.Log.Debugf("Finished processing release lines, release vars: %+v", vars)

		// convert parsed release vars to release struct and begin release processing
		if trackerRelease, err := release.FromMap(p.Tracker, p.Log, vars); err != nil {
			p.Log.WithError(err).Error("Failed converting release vars to a release struct...")
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
		queuedLine, ok := <-queue
		if !ok {
			// the channel has been closed
			return "", errors.New("line queue has been closed")
		}

		// should ignore this line?
		if p.shouldIgnoreLine(queuedLine) {
			continue
		}

		return queuedLine, nil
	}
}
