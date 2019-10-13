package processor

import (
	"github.com/l3uddz/trackarr/autodl/parser"
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
	/* private */
	log     *logrus.Entry
	tracker *parser.TrackerInfo
	cfg     *config.TrackerConfiguration
	lines   []string
}

/* Public */

func New(log *logrus.Entry, tracker *parser.TrackerInfo, config *config.TrackerConfiguration) *Processor {
	return &Processor{
		log:     log,
		tracker: tracker,
		cfg:     config,
		lines:   []string{},
	}
}

/* Private */
