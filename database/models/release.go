package models

import (
	"time"

	stringutils "github.com/l3uddz/trackarr/utils/strings"

	"github.com/jinzhu/gorm"
)

type PushedRelease struct {
	gorm.Model
	Name        string `sql:"type:text;not null"`
	TrackerName string `sql:"type:varchar(256);not null"`
	PvrName     string `sql:"type:text;not null"`
	Approved    bool   `sql:"type:bool;not null,DEFAULT:false"`
}

/* Struct Methods */
func (r PushedRelease) DurationSinceCreated() string {
	return stringutils.TimeDiffDurationString(time.Now(), r.CreatedAt, true)
}

/* Methods */
func NewPushedRelease(db *gorm.DB, name string, trackerName string, pvrName string, approved bool) (*PushedRelease, error) {
	release := &PushedRelease{
		Name:        name,
		TrackerName: trackerName,
		PvrName:     pvrName,
		Approved:    approved,
	}

	if err := db.FirstOrInit(&release, PushedRelease{Name: name, TrackerName: trackerName, PvrName: pvrName}).Error; err != nil {
		log.WithError(err).Errorf("Failed unexpectedly finding existing pushed release with name: %q", name)
		return nil, err
	}

	return release, nil
}

func GetLatestPushedReleases(db *gorm.DB, count int) []*PushedRelease {
	var releases []*PushedRelease

	if count == 0 {
		db.Order("id desc").Find(&releases)

	} else {
		db.Limit(count).Order("id desc").Find(&releases)
	}
	return releases
}

func GetLatestApprovedReleases(db *gorm.DB, count int) []*PushedRelease {
	var releases []*PushedRelease

	if count == 0 {
		db.Where("approved = 1").Order("id desc").Find(&releases)
	} else {
		db.Where("approved = 1").Limit(count).Order("id desc").Find(&releases)
	}
	return releases
}
