package tracker

import (
	"github.com/l3uddz/trackarr/config"
	"github.com/l3uddz/trackarr/logger"
	"github.com/l3uddz/trackarr/runtime"
	"github.com/l3uddz/trackarr/tracker/parser"

	"github.com/pkg/errors"
)

var (
	log = logger.GetLogger("tracker")
)

/* Public */

func Init() error {
	log.Infof("Initializing trackers...")

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
