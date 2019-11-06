package tasks

import (
	"github.com/asdine/storm/q"
	"github.com/asdine/storm/v3"
	"github.com/l3uddz/trackarr/database"
	"github.com/l3uddz/trackarr/database/models"
	"time"
)

/* Const */
const CronTaskDatabasePruner = "0 0,6,12,18 * * *"

/* Private */

func taskDatabasePruner() {
	var releases []*models.PushedRelease

	log.Debugf("Database: Pruning releases from database...")

	// find releases older than X days
	oldestDate := time.Now().UTC().Add(-time.Duration(72) * time.Hour)

	query := database.DB.Select(q.Lte("CreatedAt", oldestDate))
	if err := query.Find(&releases); err != nil && err != storm.ErrNotFound {
		log.WithError(err).Errorf("Failed finding releases to prune from before: %s", oldestDate)
		return
	}

	// remove found releases
	releasesCount := len(releases)
	if releasesCount > 0 {
		if err := query.Delete(new(models.PushedRelease)); err != nil {
			log.WithError(err).Errorf("Failed pruning %d releases from before: %s", oldestDate)
			return
		}
	}

	log.Infof("Database: Pruned %d releases from before %s", releasesCount, oldestDate)
}
