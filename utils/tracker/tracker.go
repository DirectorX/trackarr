package tracker

import (
	"fmt"
	"gitlab.com/cloudb0x/trackarr/config"
	"gitlab.com/cloudb0x/trackarr/logger"
	"strings"
	"sync"
)

/* Interface */
type Interface interface {
	GetReleaseInfo(info *config.ReleaseInfo) (*TorrentInfo, error)
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
	// API Interfaces
	apiInterfaces map[string]Interface
	mtx           sync.Mutex
)

/* Public */
func GetApi(tracker *config.TrackerInstance) (Interface, error) {
	// acquire lock
	mtx.Lock()
	defer mtx.Unlock()

	// determine tracker name to check
	trackerName := strings.ToLower(tracker.Name)
	if tracker.Info.LongName != "" {
		trackerName = strings.ToLower(tracker.Info.LongName)
	}

	// ensure tracker api map is initialized
	if apiInterfaces == nil {
		apiInterfaces = make(map[string]Interface)
		log.Trace("Initialized tracker apiInterfaces map")
	}

	// api already initialized?
	if api, ok := apiInterfaces[trackerName]; ok {
		// return already initialized api
		return api, nil
	}

	// get appropriate api interface
	var api Interface
	switch trackerName {
	case "passthepopcorn":
		api = &Ptp{
			log:     log.WithField("api", trackerName),
			tracker: tracker,
		}

		log.Debugf("Initialized API Interface for tracker: %q", trackerName)

	default:
		return nil, fmt.Errorf("api not implemented for tracker: %q", trackerName)

	}

	apiInterfaces[trackerName] = api
	return api, nil
}
