package parser

import (
	"github.com/l3uddz/trackarr/logger"
	"github.com/pkg/errors"
	"os"
)

var (
	log = logger.GetLogger("autodl")
)

/* Struct */

type Parser struct {
	/* privates */
	trackerName     string
	trackerFilePath string
}

/* Public */

func Init(tracker string, trackersPath string) (*Parser, error) {
	// validate tracker file exists
	if _, err := os.Stat(trackersPath); os.IsNotExist(err) {
		log.WithError(err).Errorf("Failed initializing parser for tracker: %q", tracker)
		return nil, errors.Wrapf(err, "failed to initialize parser for: %q", trackersPath)
	}

	// parser tracker file

	return &Parser{trackerName: tracker, trackerFilePath: trackersPath}, nil
}

/* Private */
