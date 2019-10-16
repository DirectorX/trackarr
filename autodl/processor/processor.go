package processor

import (
	"github.com/l3uddz/trackarr/autodl/parser"
	"github.com/l3uddz/trackarr/config"
	"github.com/l3uddz/trackarr/logger"
	"github.com/sirupsen/logrus"
	"strings"

	"github.com/enriquebris/goconcurrentqueue"
)

/* Vars */
var (
	log = logger.GetLogger("autodl")
)

/* Structs */

type Processor struct {
	/* public */
	Log     *logrus.Entry
	Tracker *parser.TrackerInfo
	Cfg     *config.TrackerConfiguration

	/* private */
	queues  map[string]*goconcurrentqueue.FIFO
}

/* Public */

func New(log *logrus.Entry, tracker *parser.TrackerInfo, config *config.TrackerConfiguration) *Processor {
	// initialize queues
	queues := make(map[string]*goconcurrentqueue.FIFO, 0)
	for _, channel := range tracker.Channels {
		queues[strings.ToLower(channel)] = goconcurrentqueue.NewFIFO()
	}

	// create processor
	processor := &Processor{
		Log:     log,
		Tracker: tracker,
		Cfg:     config,
		queues:  queues,
	}

	// init queue processors
	for _, queue := range processor.queues {
		go processor.processQueue(queue)
	}

	return processor
}

/* Private */
