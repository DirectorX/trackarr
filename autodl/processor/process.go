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
		vars := p.matchPatterns(&p.tracker.LinePatterns, line)
		if len(vars) == 0 {
			// vars were not matched/parsed
			return nil
		}

		// run vars against rules
		if err := p.processRules(p.tracker.LineMatchedRules, &vars); err != nil {
			return err
		}

		// log final vars map
		p.log.Debugf("Vars post linematched processed: %+v", vars)
		return nil

	} else if len(p.tracker.MultiLinePatterns) > 0 {
		// use multi-linepatterns

	} else {
		// unknown??
		return errors.New("unable to determine how to pattern match")
	}

	return nil
}
