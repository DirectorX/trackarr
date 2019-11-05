package models

import (
	"github.com/asdine/storm/v3"
	"github.com/l3uddz/trackarr/logger"
)

var (
	log = logger.GetLogger("dbm")
)

// Tracker - Model representation of an autodl tracker file
type Tracker struct {
	Name    string `storm:"id,unique"`
	Version string
}

/* Methods */

// NewOrExistingTracker - Return an existing or new tracker
func NewOrExistingTracker(db *storm.DB, name string) (*Tracker, error) {
	var tracker Tracker

	// find existing tracker
	if err := db.One("Name", name, &tracker); err == nil {
		return &tracker, nil
	}

	// create new tracker
	tracker.Name = name
	if err := db.Save(&tracker); err != nil {
		return nil, err
	}

	return &tracker, nil
}
