package processor

/* Private */

func (p *Processor) shouldIgnoreLine(line string) bool {
	// iterate ignore lines
	for _, ignore := range p.Tracker.Info.IgnoreLines {
		if ignore.Rxp.MatchString(line) && ignore.Expected {
			// ignore this message as it matched an ignore pattern
			p.Log.Tracef("Ignoring message as ignore pattern met %q: %s", ignore.Rxp, line)
			return true
		}
	}

	return false
}
