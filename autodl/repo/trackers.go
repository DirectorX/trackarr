package repo

import (
	"fmt"
	"github.com/l3uddz/trackarr/logger"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/l3uddz/trackarr/database"
	models "github.com/l3uddz/trackarr/database/models"
	stringutils "github.com/l3uddz/trackarr/utils/strings"
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
var (
	log = logger.GetLogger("autodl")
)

const trackersRepository = "https://github.com/autodl-community/autodl-trackers/tree/master/trackers"

/* Public */

// PullTrackers - Process all available trackers looking for new/changed trackers to pull
func PullTrackers(trackersPath string) error {
	// ensure tracker directory exists
	if _, err := os.Stat(trackersPath); os.IsNotExist(err) {
		if err := os.Mkdir(trackersPath, 0700); err != nil {
			log.WithError(err).Errorf("Failed to create tracker directory: %q", trackersPath)
			return errors.Wrap(err, "failed creating tracker directory")
		} else {
			log.Infof("Created tracker directory: %q", trackersPath)
		}
	}

	// retrieve trackers
	trackers, err := getAvailableTrackers()
	if err != nil {
		return err
	}

	// process found trackers
	trackerPulls := 0
	trackerErrors := 0
	for _, trackerData := range *trackers {
		log.Tracef("Processing tracker: %s", trackerData.Name)

		// retrieve tracker from database
		tracker, err := models.NewOrExistingTracker(database.DB, trackerData.Name)
		if err != nil {
			log.WithError(err).Errorf("Failed retrieving tracker from database: %q", trackerData.Name)
			return errors.Wrap(err, "failed retrieving tracker from database")
		}

		// grab tracker if required
		trackerPath := filepath.Join(trackersPath, trackerData.Name+".tracker")
		if _, err := os.Stat(trackerPath); os.IsNotExist(err) || tracker.Version != trackerData.Version {
			// the tracker file did not exist, or we were using an old version, we must download it
			log.Infof("Pulling tracker: %s -> %q", trackerData.Name, trackerPath)

			if err := pullTracker(trackerData.URL, trackerPath); err != nil {
				// failed to pull this tracker
				trackerErrors++
				continue
			}

			// update tracker in database
			tracker.Version = trackerData.Version
			database.DB.Save(&tracker)

			trackerPulls++
		} else {
			log.Tracef("No pull required for tracker: %s", trackerData.Name)
		}
	}

	if trackerPulls > 0 || trackerErrors > 0 {
		log.Infof("Pulled %d %s with %d %s", trackerPulls, stringutils.Pluralize("tracker", trackerPulls),
			trackerErrors, stringutils.Pluralize("failure", trackerErrors))
	} else {
		log.Infof("Trackers are up to date")
	}

	return nil
}

/* Private */

// getAvailableTrackers - Retrieve all available trackers from autodl-community repository
func getAvailableTrackers() (*map[string]*AutodlTracker, error) {
	// retrieve trackers page
	log.Infof("Finding available trackers from: %s", trackersRepository)
	body, err := web.GetBody(web.GET, trackersRepository, 30)
	if err != nil {
		return nil, err
	}

	// parse trackers from body
	rxp := regexp.MustCompile(
		`title="(?P<Name>.+)\.tracker" id="(?P<Version>.+)" href="(?P<URL>.+\.tracker)">.+</a>`)
	matches := rxp.FindAllStringSubmatch(body, -1)

	// build trackers map
	trackers := make(map[string]*AutodlTracker, 0)
	for _, match := range matches {
		// parse tracker from match
		tracker := &AutodlTracker{
			Name:    match[1],
			Version: match[2],
			URL: fmt.Sprintf("https://raw.githubusercontent.com%s",
				strings.Replace(match[3], "/blob/", "/", -1)),
		}
		log.Tracef("Available tracker: %q - Version: %q - URL: %s", tracker.Name, tracker.Version, tracker.URL)

		// add tracker to map
		trackers[tracker.Name] = tracker
	}
	log.Infof("Found %d trackers", len(trackers))
	return &trackers, nil
}

// pullTracker - Download a tracker and save to specified path
func pullTracker(url string, path string) error {
	// download tracker
	trackerData, err := web.GetBody(web.GET, url, 30)
	if err != nil {
		log.WithError(err).Errorf("Failed pulling tracker: %s", url)
		return errors.Wrap(err, "failed downloading tracker")
	}

	// TODO: validate tracker is in expected XML format that we are able to parse later on

	// save to tracker file
	file, err := os.Create(path)
	if err != nil {
		log.WithError(err).Errorf("Failed creating tracker: %q", path)
		return errors.Wrap(err, "failed creating tracker file")
	}
	defer file.Close()

	if _, err := file.WriteString(trackerData); err != nil {
		log.WithError(err).Errorf("Failed writing tracker: %q", path)
		return errors.Wrap(err, "failed writing tracker file")
	}

	return nil
}
