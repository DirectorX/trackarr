package tracker

import (
	"gitlab.com/cloudb0x/trackarr/config"
	"gitlab.com/cloudb0x/trackarr/logger"
	"gitlab.com/cloudb0x/trackarr/runtime"
	"gitlab.com/cloudb0x/trackarr/tracker/parser"

	"github.com/pkg/errors"
)

var (
	log = logger.GetLogger("tracker")
)

/* Public */

func Init() error {
	log.Info("Initializing trackers...")

	for trackerName, t := range config.Config.Trackers {
		// skip disabled trackers
		if !t.Enabled {
			log.Debugf("Skipping disabled tracker: %s", trackerName)

			continue
		}

		t2 := t
		trackerInstance := &config.TrackerInstance{
			Name:   trackerName,
			Config: &t2,
			Info:   parser.New(trackerName),
		}

		if err := parser.Parse(trackerInstance.Info); err != nil {
			return errors.Wrapf(err, "parsing tracker %s", trackerName)
		}
		log.Debugf("Parsed tracker: %s", trackerName)

		// validate required config settings were set for this tracker
		settingsFilled := true
		for _, trackerSetting := range trackerInstance.Info.Settings {
			if _, ok := trackerInstance.Config.Settings[trackerSetting]; !ok {
				log.Warnf("Skipping tracker %s, missing config setting: %q", trackerName, trackerSetting)
				settingsFilled = false
				break
			}
		}

		if !settingsFilled {
			// there were missing config settings that were required by this tracker
			continue
		}

		runtime.Tracker[trackerName] = trackerInstance
	}

	return nil
}
