package autodl

import (
	"github.com/l3uddz/trackarr/autodl/repo"
	"github.com/l3uddz/trackarr/config"
	"github.com/l3uddz/trackarr/logger"
	stringutils "github.com/l3uddz/trackarr/utils/strings"
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
