package processor

import (
	"fmt"
)

/* Public */

func (p *Processor) QueueLine(channel string, line string) error {
	// get channel specific queue from queues
	queue, ok := p.queues[channel]
	if !ok {
		return fmt.Errorf("no queue was initialized for channel: %q", channel)
	}

	// add line to queued items
	p.Log.Tracef("Adding line to queue: %s", line)
	queue <- line
	p.Log.Tracef("Queued line for processing: %s", line)
	return nil
}
