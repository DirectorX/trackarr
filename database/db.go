package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/l3uddz/trackarr/database/models"
	"github.com/l3uddz/trackarr/logger"
	stringutils "github.com/l3uddz/trackarr/utils/strings"
)

/* Vars */
var (
	// DB exports database object
	DB *gorm.DB
	// Package logging
	log = logger.GetLogger("db ")
)

/* Public */

// Init - Init connection to the database
func Init(dbPath string) error {
	var err error
	DB, err = gorm.Open("sqlite3", dbPath)
	if err != nil {
		log.WithError(err).Fatalf("Failed initializing database connection to %q", dbPath)
		return err
	}

	// migrate
	DB.AutoMigrate(
		&models.Tracker{},
	)

	// log
	log.Infof("Using %s = %q", stringutils.StringLeftJust("DATABASE", " ", 10), dbPath)
	return nil
}
