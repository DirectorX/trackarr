package processor

import (
	"fmt"
	"github.com/enriquebris/goconcurrentqueue"
	"github.com/l3uddz/trackarr/autodl/parser"
	"github.com/l3uddz/trackarr/utils/maps"
	stringutils "github.com/l3uddz/trackarr/utils/strings"
	"github.com/pkg/errors"
)

/* Private */

func (p *Processor) processQueue(queue *goconcurrentqueue.FIFO) {
	var patterns []parser.TrackerPattern

	// set patterns
	if len(p.tracker.LinePatterns) > 0 {
		patterns = p.tracker.LinePatterns
	} else if len(p.tracker.MultiLinePatterns) > 0 {
		patterns = p.tracker.MultiLinePatterns
	} else {
		p.log.Fatalf("Failed determining pattern type for processor...")
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
				p.log.WithError(err).Errorf("Failed dequeueing line to process...")
				goto RetryPattern
			}

			// process line
			p.log.Debugf("Processing line: %s", line)
			patternVars, err := p.matchPattern(&pattern, line)
			if err != nil {
				p.log.WithError(err).Errorf("Failed matching pattern, discarding release...")
				goto NewRelease
			}

			// update vars
			maps.MergeStringMap(&vars, &patternVars)
		}

		// finished parsing release - process rules
		if err := p.processRules(p.tracker.LineMatchedRules, &vars); err != nil {
			p.log.WithError(err).Errorf("failed processing release due to rules failure...")
			goto NewRelease
		}

		// push release
		p.log.Debugf("Finished processing: %s", stringutils.JsonifyLax(vars))
	}
}

func (p *Processor) nextGoodLine(queue *goconcurrentqueue.FIFO) (string, error) {
	for {
		// pop line from queue
		queuedLine, err := queue.DequeueOrWaitForNextElement()
		if err != nil {
			return "", errors.Wrap(err, "failed dequeueing next line to process")
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
