package autodl

import (
	"gitlab.com/cloudb0x/trackarr/autodl/repo"
	"gitlab.com/cloudb0x/trackarr/config"
	"gitlab.com/cloudb0x/trackarr/logger"
	stringutils "gitlab.com/cloudb0x/trackarr/utils/strings"
)

var (
	log = logger.GetLogger("autodl")
)

/* Public */

func Init() error {
	// info log
	log.Infof("Using %s = %q", stringutils.StringLeftJust("TRACKERS", " ", 10), config.Runtime.Trackers)

	// pull the latest autodl-community trackers
	if err := repo.PullTrackers(config.Runtime.Trackers); err != nil {
		return err
	}

	return nil
}
