package tasks

import (
	"fmt"
	"time"

	"gitlab.com/cloudb0x/trackarr/database"
	"gitlab.com/cloudb0x/trackarr/database/models"
	"gitlab.com/cloudb0x/trackarr/ws"

	"github.com/asdine/storm/v3/q"
	"github.com/asdine/storm/v3"
)

type TaskPruner struct {
	Cron        string
	MaxAgeHours time.Duration
}

/* Private */

func (t *Tasks) taskDatabasePruner() {
	var releases []*models.PushedRelease

	log.Debug("Database: Pruning releases from database...")

	// find releases older than X days
	oldestDate := time.Now().UTC().Add(-t.TaskPruner.MaxAgeHours * time.Hour)

	query := database.DB.Select(q.Lte("CreatedAt", oldestDate))
	if err := query.Find(&releases); err != nil && err != storm.ErrNotFound {
		log.WithError(err).Errorf("Failed finding releases to prune from before: %s", oldestDate)
		return
	}

	// remove found releases
	releasesCount := len(releases)
	if releasesCount > 0 {
		if err := query.Delete(new(models.PushedRelease)); err != nil {
			log.WithError(err).Errorf("Failed pruning %d releases from before: %s", releasesCount, oldestDate)
			return
		}

		// broadcast alert to sockets
		jsonData, err := ws.NewAlert("info", "Database",
			fmt.Sprintf("%d releases pruned", releasesCount))
		if err == nil {
			ws.BroadcastAll(jsonData)
		} else {
			log.WithError(err).Error("Failed creating database pruner websocket alert")
		}
	}

	log.Infof("Database: Pruned %d releases from before %s", releasesCount, oldestDate)
}
