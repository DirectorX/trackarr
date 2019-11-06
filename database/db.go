package database

import (
	"github.com/l3uddz/trackarr/config"
	"github.com/l3uddz/trackarr/logger"
	stringutils "github.com/l3uddz/trackarr/utils/strings"

	"github.com/asdine/storm/v3"
	"github.com/pkg/errors"
)

/* Vars */
var (
	// DB exports database object
	DB *storm.DB
	// Package logging
	log = logger.GetLogger("db ")
)

/* Public */

// Init - Initialize connection to the database
func Init() error {
	var err error
	DB, err = storm.Open(config.Runtime.DB)
	if err != nil {
		log.WithError(err).Fatalf("Failed initializing database connection to %q", config.Runtime.DB)
		return errors.Wrap(err, "failed initializing database connection")
	}

	// log
	log.Infof("Using %s = %q", stringutils.StringLeftJust("DATABASE", " ", 10), config.Runtime.DB)
	return nil
}
