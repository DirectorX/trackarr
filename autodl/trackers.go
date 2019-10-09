package autodl

import (
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

	log.Info(body)
	return nil
}
