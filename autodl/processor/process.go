package processor

import (
	"fmt"
	"github.com/enriquebris/goconcurrentqueue"
	"github.com/l3uddz/trackarr/autodl/parser"
	"github.com/l3uddz/trackarr/release"
	"github.com/l3uddz/trackarr/utils/maps"
	"github.com/pkg/errors"
)

/* Private */

func (p *Processor) processQueue(queue *goconcurrentqueue.FIFO) {
	var patterns []parser.TrackerPattern

	// set patterns
	if len(p.Tracker.LinePatterns) > 0 {
		patterns = p.Tracker.LinePatterns
	} else if len(p.Tracker.MultiLinePatterns) > 0 {
		patterns = p.Tracker.MultiLinePatterns
	} else {
		p.Log.Fatalf("Failed determining pattern type for processor...")
		return
	}

	// iterate patterns
	for {
	NewRelease:
		vars := map[string]string{}
		for _, pattern := range patterns {
		RetryPattern:
			// iterate each pattern finding a match
			line, err := p.nextGoodLine(queue)
			if err != nil {
				p.Log.WithError(err).Errorf("Failed dequeuing line to process...")
				goto RetryPattern
			}

			// process line
			p.Log.Debugf("Processing line: %s", line)
			patternVars, err := p.matchPattern(&pattern, line)
			if err != nil {
				p.Log.WithError(err).Errorf("Failed matching pattern, discarding release...")
				goto NewRelease
			}

			// update vars
			maps.MergeStringMap(&vars, &patternVars)
		}

		// finished parsing release lines - process rules
		if err := p.processRules(p.Tracker.LineMatchedRules, &vars); err != nil {
			p.Log.WithError(err).Errorf("Failed processing release lines due to rules failure...")
			goto NewRelease
		}

		p.Log.Debugf("Finished processing release lines, release vars: %+v", vars)

		// convert parsed release vars to release struct and begin release processing
		if trackerRelease, err := release.FromMap(p.Tracker, p.Cfg, p.Log, &vars); err != nil {
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
