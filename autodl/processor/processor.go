package processor

import (
	"strings"

	"github.com/l3uddz/trackarr/config"
	"github.com/l3uddz/trackarr/logger"

	"github.com/enriquebris/goconcurrentqueue"
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
	queues map[string]*goconcurrentqueue.FIFO
}

/* Public */

func New(log *logrus.Entry, t *config.TrackerInstance) *Processor {
	// initialize queues
	queues := make(map[string]*goconcurrentqueue.FIFO)
	for _, channel := range t.Info.Channels {
		queues[strings.ToLower(channel)] = goconcurrentqueue.NewFIFO()
	}

	// create processor
	processor := &Processor{
		Log:     log,
		Tracker: t,
		queues:  queues,
	}

	// init queue processors
	for _, queue := range processor.queues {
		go func(q *goconcurrentqueue.FIFO) {
			processor.processQueue(q)
		}(queue)
	}

	return processor
}

/* Private */
