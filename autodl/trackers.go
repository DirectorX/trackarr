package autodl

/* Vars / Const */
const trackersRepository = "https://github.com/autodl-community/autodl-trackers/tree/master/trackers"

/* Public */

func PullTrackers(trackersPath string) error {
	log.Debugf("Pulling latest trackers from %q", trackersRepository)
	return nil
}
