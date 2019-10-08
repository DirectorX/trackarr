package autodl

import (
	"github.com/l3uddz/trackarr/logger"
	stringutils "github.com/l3uddz/trackarr/utils/strings"
)

var (
	log = logger.GetLogger("autodl")
)

/* Public */

func Init(trackersPath string) error {
	// info log
	log.Infof("Using %s = %q", stringutils.StringLeftJust("TRACKERS", " ", 10), trackersPath)

	// pull the latest autodl-community trackers
	if err := PullTrackers(trackersPath); err != nil {
		return err
	}

	return nil
}
