package database

import (
	"gitlab.com/cloudb0x/trackarr/config"
	"gitlab.com/cloudb0x/trackarr/logger"
	stringutils "gitlab.com/cloudb0x/trackarr/utils/strings"

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
		return errors.Wrapf(err, "failed initializing database connection to %q", config.Runtime.DB)
	}

	// log
	log.Infof("Using %s = %q", stringutils.StringLeftJust("DATABASE", " ", 10), config.Runtime.DB)
	return nil
}
