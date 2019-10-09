package models

import (
	"github.com/jinzhu/gorm"
	"github.com/l3uddz/trackarr/logger"
)

var (
	log = logger.GetLogger("dbm")
)

// Tracker - Model representation of an autodl tracker file
type Tracker struct {
	gorm.Model
	Name    string `sql:"type:varchar(256);unique;not null"`
	Version string `sql:"type:varchar(256);not null"`
}

/* Methods */

// NewOrExistingTracker - Return an existing or new tracker
func NewOrExistingTracker(db *gorm.DB, name string) (*Tracker, error) {
	tracker := &Tracker{}

	if err := db.FirstOrInit(&tracker, Tracker{Name: name}).Error; err != nil {
		log.WithError(err).Errorf("Failed to unexpectedly tracker with name: %q", name)
		return nil, err
	}

	return tracker, nil
}
