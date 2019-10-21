package processor

import (
	"fmt"

	"github.com/l3uddz/trackarr/config"
	"github.com/l3uddz/trackarr/release"
	"github.com/l3uddz/trackarr/utils/maps"

	"github.com/enriquebris/goconcurrentqueue"
	"github.com/pkg/errors"
)

/* Private */

func (p *Processor) processQueue(queue *goconcurrentqueue.FIFO) {
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
			maps.MergeStringMap(&vars, &patternVars)
		}

		if parseFailed {
			// pattern parsing above failed to parse a complete release
			continue
		}

		// finished parsing release lines - process rules
		if err := p.processRules(p.Tracker.Info.LineMatchedRules, &vars); err != nil {
			p.Log.WithError(err).Errorf("Failed processing release lines due to rules failure...")
			continue
		}

		p.Log.Debugf("Finished processing release lines, release vars: %+v", vars)

		// convert parsed release vars to release struct and begin release processing
		if trackerRelease, err := release.FromMap(p.Tracker, p.Log, &vars); err != nil {
			p.Log.WithError(err).Errorf("Failed converting release vars to a release struct...")
		} else {
			// start processing this release
			go trackerRelease.Process()
		}
	}
}

func (p *Processor) nextGoodLine(queue *goconcurrentqueue.FIFO) (string, error) {
	for {
		// pop line from queue
		queuedLine, err := queue.DequeueOrWaitForNextElement()
		if err != nil {
			return "", errors.Wrap(err, "failed dequeuing next line to process")
		}

		// type assert line
		line, ok := queuedLine.(string)
		if !ok {
			return "", fmt.Errorf("failed type asserting dequeued line: %#v", queuedLine)
		}

		// should ignore this line?
		if p.shouldIgnoreLine(line) {
			continue
		}

		return line, nil
	}
}
