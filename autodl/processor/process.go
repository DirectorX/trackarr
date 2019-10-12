package processor

/* Public */

func (p *Processor) ProcessLine(line string) error {
	// should we ignore this line
	if p.shouldIgnoreLine(line) {
		return nil
	}

	// process line matching patterns

	return nil
}
