package processor

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

/* Public */

func (p *Processor) QueueLine(channel string, line string) error {
	// get channel specific queue from queues
	lowerChannel := strings.ToLower(channel)
	queue, ok := p.queues[lowerChannel]
	if !ok {
		return fmt.Errorf("no queue was initialized for channel: %q", lowerChannel)
	}

	// add line to queued items
	if err := queue.Enqueue(line); err != nil {
		return errors.Wrapf(err, "failed queueing line for processing")
	}

	p.log.Tracef("Queued line for processing, queue size for %s: %d", channel, queue.GetLen())
	return nil
}
