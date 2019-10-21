package parser

import (
	"os"
	"path/filepath"

	"github.com/l3uddz/trackarr/config"

	"github.com/pkg/errors"
)

func trackerFilePath(t *config.TrackerInfo) (*string, error) {
	trackerFilePath := filepath.Join(config.Runtime.Trackers, t.Name+".tracker")
	if _, err := os.Stat(trackerFilePath); os.IsNotExist(err) {
		return nil, errors.Wrapf(err, "failed locating tracker file: %q", trackerFilePath)
	}

	return &trackerFilePath, nil
}
