package models

import (
	"time"

	stringutils "gitlab.com/cloudb0x/trackarr/utils/strings"

	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
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
	return stringutils.TimeDiffDurationString(time.Now().UTC(), r.CreatedAt, true)
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
		return nil, err
	}

	return release, nil
}

func GetLatestPushedReleases(db *storm.DB, count int) []*PushedRelease {
	releases := make([]*PushedRelease, 0)

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
	releases := make([]*PushedRelease, 0)

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
