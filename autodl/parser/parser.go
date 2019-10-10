package parser

import (
	"github.com/l3uddz/trackarr/logger"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
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
	trackerFilePath := filepath.Join(trackersPath, tracker+".tracker")
	if _, err := os.Stat(trackerFilePath); os.IsNotExist(err) {
		log.WithError(err).Errorf("Failed initializing parser for tracker: %q", tracker)
		return nil, errors.Wrapf(err, "failed to initialize parser for: %q", trackerFilePath)
	}

	// parser tracker file

	return &Parser{trackerName: tracker, trackerFilePath: trackerFilePath}, nil
}

/* Private */
