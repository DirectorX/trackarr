package processor

import (
	"github.com/l3uddz/trackarr/config"
	"github.com/l3uddz/trackarr/logger"

	"github.com/sirupsen/logrus"
)

/* Vars */
var (
	log = logger.GetLogger("autodl")
)

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
	for _, queue := range processor.queues {
		go func(q chan string) {
			processor.processQueue(q)
		}(queue)
	}

	return processor
}

/* Private */
