package paths

import (
	"os"
	"path/filepath"

	"gitlab.com/cloudb0x/trackarr/logger"

	"github.com/sirupsen/logrus"
)

var log = logger.GetLogger("paths")

func GetCurrentBinaryPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.WithFields(logrus.Fields{"prefix": "GetCurrentBinaryPath"}).
			WithError(err).
			Error("Failed to retrieve current binary path")

		// get current working dir
		if dir, err = os.Getwd(); err != nil {
			log.WithFields(logrus.Fields{"prefix": "GetCurrentBinaryPath"}).
				WithError(err).
				Error("Failed to retrieve current working path")
			os.Exit(1)
		}
	}
	return dir
}
