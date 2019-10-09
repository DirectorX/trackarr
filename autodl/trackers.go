package autodl

import (
	"regexp"

	"github.com/l3uddz/trackarr/utils/web"
)

/* Vars / Const */
const trackersRepository = "https://github.com/autodl-community/autodl-trackers/tree/master/trackers"

/* Public */

// PullTrackers - Process the autodl-community trackers folder looking for new/changed trackers to pull down
func PullTrackers(trackersPath string) error {
	log.Debugf("Pulling latest trackers from %q", trackersRepository)
	_ = getLatestTrackers()
	return nil
}

/* Private */

func getLatestTrackers() error {
	// retrieve trackers page
	body, err := web.GetBody(web.GET, trackersRepository, 30)
	if err != nil {
		return err
	}

	// parse trackers from body
	rxp := regexp.MustCompile(`title="(?P<Name>.+)\.tracker" id="(?P<Version>.+)" href="(?P<URL>.+\.tracker)">.+</a>`)
	matches := rxp.FindAllStringSubmatch(body, -1)

	// iterate through matches
	trackers := 0
	for _, match := range matches {
		trackers++
		log.Infof("Tracker: %s - Value: %s", match[1], match[2])
	}
	log.Infof("Found %d trackers", trackers)
	return nil
}
