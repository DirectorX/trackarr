package database

import (
	"github.com/l3uddz/trackarr/config"
	"github.com/l3uddz/trackarr/database/models"
	"github.com/l3uddz/trackarr/logger"
	stringutils "github.com/l3uddz/trackarr/utils/strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pkg/errors"
)

/* Vars */
var (
	// DB exports database object
	DB *gorm.DB
	// Package logging
	log = logger.GetLogger("db ")
)

/* Public */

// Init - Initialize connection to the database
func Init() error {
	var err error
	DB, err = gorm.Open("sqlite3", config.Runtime.DB)
	if err != nil {
		log.WithError(err).Fatalf("Failed initializing database connection to %q", config.Runtime.DB)
		return errors.Wrap(err, "failed initializing database connection")
	}

	// migrate
	DB.AutoMigrate(
		&models.Tracker{},
		&models.PushedRelease{},
	)

	// log
	log.Infof("Using %s = %q", stringutils.StringLeftJust("DATABASE", " ", 10), config.Runtime.DB)
	return nil
}
