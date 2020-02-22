package processor

import (
	"gitlab.com/cloudb0x/trackarr/config"

	"github.com/sirupsen/logrus"
)

/* Vars */
var ()

/* Structs */

type Processor struct {
	/* public */
	Log     *logrus.Entry
	Tracker *config.TrackerInstance

	/* private */
	queues map[string]chan string
}

/* Public */

func New(log *logrus.Entry, t *config.TrackerInstance) *Processor {
	// initialize queues
	queues := make(map[string]chan string)
	for _, channel := range t.Info.Channels {
		queues[channel] = make(chan string, 128)
	}

	// create processor
	processor := &Processor{
		Log:     log,
		Tracker: t,
		queues:  queues,
	}

	// init queue processors
	for queueName, queue := range processor.queues {
		go func(name string, q chan string) {
			// log start
			log.Tracef("Queue processor started: %q", name)
			// process queue
			processor.processQueue(q)
			// log finish
			log.Tracef("Queue processor finished: %q", name)
		}(queueName, queue)
	}

	return processor
}

/* Private */
