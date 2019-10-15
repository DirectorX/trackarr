package processor

import (
	"github.com/l3uddz/trackarr/autodl/parser"
	"github.com/l3uddz/trackarr/config"
	"github.com/l3uddz/trackarr/logger"
	"github.com/sirupsen/logrus"

	"github.com/enriquebris/goconcurrentqueue"
)

/* Vars */
var (
	log = logger.GetLogger("autodl")
)

/* Structs */

type Processor struct {
	/* private */
	log     *logrus.Entry
	tracker *parser.TrackerInfo
	cfg     *config.TrackerConfiguration
	queues  map[string]*goconcurrentqueue.FIFO
}

/* Public */

func New(log *logrus.Entry, tracker *parser.TrackerInfo, config *config.TrackerConfiguration) *Processor {
	// initialize queues
	queues := make(map[string]*goconcurrentqueue.FIFO, 0)
	for _, channel := range tracker.Channels {
		queues[channel] = goconcurrentqueue.NewFIFO()
	}

	return &Processor{
		log:     log,
		tracker: tracker,
		cfg:     config,
		queues:  queues,
	}
}

/* Private */
