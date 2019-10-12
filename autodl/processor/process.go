package processor

import "github.com/pkg/errors"

/* Public */

func (p *Processor) ProcessLine(line string) error {
	// should we ignore this line
	if p.shouldIgnoreLine(line) {
		return nil
	}

	// process line matching patterns
	if len(p.tracker.LinePatterns) > 0 {
		// use linepatterns
		_ = p.matchPatterns(&p.tracker.LinePatterns, line)
	} else if len(p.tracker.MultiLinePatterns) > 0 {
		// use multi-linepatterns

	} else {
		// unknown??
		log.Errorf("Unsure how to pattern match: %s", line)
		return errors.New("unable to determine how to pattern match")
	}

	return nil
}
