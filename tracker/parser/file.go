package parser

import (
	"fmt"
	"github.com/l3uddz/trackarr/config"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

func trackerFilePath(t *config.TrackerInfo) (*string, error) {
	// TODO: remove this once viper.SetCaseInsensitive is public
	trackerName := ""

	err := filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// skip dirs
			if info.IsDir() {
				return nil
			}

			if trackerName != "" {
				return nil
			}

			if strings.ToLower(info.Name()) == (strings.ToLower(t.Name) + ".tracker") {
				trackerName = info.Name()
			}

			return nil
		})

	if err != nil {
		return nil, errors.Wrapf(err, "failed finding trackers in: %s", config.Runtime.Trackers)
	}

	if trackerName == "" {
		return nil, fmt.Errorf("failed finding tracker file with name: %s.tracker", t.Name)
	}

	trackerFilePath := filepath.Join(config.Runtime.Trackers, trackerName)
	if _, err := os.Stat(trackerFilePath); os.IsNotExist(err) {
		return nil, errors.Wrapf(err, "failed locating tracker file: %q", trackerFilePath)
	}

	return &trackerFilePath, nil
}
