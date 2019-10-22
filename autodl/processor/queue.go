package processor

import (
	"fmt"

	"github.com/pkg/errors"
)

/* Public */

func (p *Processor) QueueLine(channel string, line string) error {
	// get channel specific queue from queues
	queue, ok := p.queues[channel]
	if !ok {
		return fmt.Errorf("no queue was initialized for channel: %q", channel)
	}

	// add line to queued items
	if err := queue.Enqueue(line); err != nil {
		return errors.Wrapf(err, "failed queueing line for processing")
	}

	p.Log.Tracef("Queued line for processing, queue size for %s: %d", channel, queue.GetLen())
	return nil
}
