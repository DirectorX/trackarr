package processor

import (
	"github.com/l3uddz/trackarr/autodl/parser"
	"github.com/l3uddz/trackarr/logger"
	"github.com/sirupsen/logrus"
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
	lines   []string
}

/* Public */

func New(log *logrus.Entry, tracker *parser.TrackerInfo) *Processor {
	return &Processor{
		log:     log,
		tracker: tracker,
		lines:   []string{},
	}
}

/* Private */
