package tracker

import (
	"fmt"
	"gitlab.com/cloudb0x/trackarr/config"
	"gitlab.com/cloudb0x/trackarr/logger"
	"strings"
)

/* Interface */
type Interface interface {
	GetReleaseInfo(string) (*TorrentInfo, error)
}

/* Struct */
type TorrentInfo struct {
	Name string
	Size string
}

/* Var */
var (
	// Logging
	log = logger.GetLogger("tracker")
)

/* Public */
func GetApi(tracker *config.TrackerInstance) (Interface, error) {
	switch strings.ToLower(tracker.Name) {
	case "passthepopcorn":
		return &Ptp{
			log:     log.WithField("api", tracker.Name),
			tracker: tracker,
		}, nil

	default:
		break
	}

	return nil, fmt.Errorf("api not implemented for tracker: %q", tracker.Name)
}
