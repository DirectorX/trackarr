package autodl

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/l3uddz/trackarr/utils/web"
)

/* Structs */

// AutodlTracker -- Struct representation of the autodl trackers directory
type AutodlTracker struct {
	Name    string
	Version string
	URL     string
}

/* Vars / Const */
const trackersRepository = "https://github.com/autodl-community/autodl-trackers/tree/master/trackers"

/* Public */

// PullTrackers - Process the autodl-community trackers folder looking for new/changed trackers to pull down
func PullTrackers(trackersPath string) error {
	// retrieve trackers
	trackers, err := getLatestTrackers()
	if err != nil {
		return err
	}

	// iterate trackers looking for new/changes
	for _, trackerData := range *trackers {
		log.Debugf("Processing: %s", trackerData.Name)
	}

	return nil
}

/* Private */

func getLatestTrackers() (*map[string]*AutodlTracker, error) {
	// retrieve trackers page
	log.Infof("Loading latest trackers from %q", trackersRepository)
	body, err := web.GetBody(web.GET, trackersRepository, 30)
	if err != nil {
		return nil, err
	}

	// parse trackers from body
	rxp := regexp.MustCompile(`title="(?P<Name>.+)\.tracker" id="(?P<Version>.+)" href="(?P<URL>.+\.tracker)">.+</a>`)
	matches := rxp.FindAllStringSubmatch(body, -1)

	// build trackers map
	trackers := make(map[string]*AutodlTracker, 0)
	for _, match := range matches {
		// parse tracker from match
		tracker := &AutodlTracker{
			Name:    match[1],
			Version: match[2],
			URL:     fmt.Sprintf("https://raw.githubusercontent.com%s", strings.Replace(match[3], "/blob/", "/", -1)),
		}
		log.Tracef("Tracker: %q - Version: %q - URL: %s", tracker.Name, tracker.Version, tracker.URL)

		// add tracker to map
		trackers[tracker.Name] = tracker
	}
	log.Infof("Found %d trackers in total", len(trackers))
	return &trackers, nil
}
