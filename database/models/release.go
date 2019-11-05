package models

import (
	"github.com/asdine/storm/q"
	"github.com/asdine/storm/v3"
	stringutils "github.com/l3uddz/trackarr/utils/strings"
	"time"
)

type PushedRelease struct {
	Id          int `storm:"id,increment"`
	Name        string
	TrackerName string
	PvrName     string
	Approved    bool
	CreatedAt   time.Time
}

/* Struct Methods */
func (r PushedRelease) DurationSinceCreated() string {
	return stringutils.TimeDiffDurationString(time.Now(), r.CreatedAt, true)
}

/* Methods */
func NewPushedRelease(db *storm.DB, name string, trackerName string, pvrName string, approved bool) (*PushedRelease, error) {
	release := &PushedRelease{
		Name:        name,
		TrackerName: trackerName,
		PvrName:     pvrName,
		Approved:    approved,
		CreatedAt:   time.Now().UTC(),
	}

	if err := db.Save(release); err != nil {
		log.WithError(err).Errorf("Failed unexpectedly finding existing pushed release with name: %q", name)
		return nil, err
	}

	return release, nil
}

func GetLatestPushedReleases(db *storm.DB, count int) []*PushedRelease {
	var releases []*PushedRelease

	if count == 0 {
		if err := db.All(&releases); err != nil {
			log.WithError(err).Error("Failed retrieving all pushed releases from database...")
		}
	} else {
		if err := db.All(&releases, storm.Limit(count)); err != nil {
			log.WithError(err).Errorf("Failed retrieving %d pushed releases from database...", count)
		}
	}

	return releases
}

func GetLatestApprovedReleases(db *storm.DB, count int) []*PushedRelease {
	var releases []*PushedRelease

	if count == 0 {
		if err := db.Select(q.Eq("Approved", true)).Find(&releases); err != nil && err != storm.ErrNotFound {
			log.WithError(err).Error("Failed retrieving all approved releases from database...")
		}
	} else {
		if err := db.Select(q.Eq("Approved", true)).Limit(count).Find(&releases); err != nil && err != storm.ErrNotFound {
			log.WithError(err).Errorf("Failed retrieving %d approved releases from database...", count)
		}
	}

	return releases
}
