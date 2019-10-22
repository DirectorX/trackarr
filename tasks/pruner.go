package tasks

import (
	"github.com/l3uddz/trackarr/database"
	"github.com/l3uddz/trackarr/database/models"
	"time"
)

/* Const */
const CronTaskDatabasePruner = "0 0,6,12,18 * * *"

/* Private */

func taskDatabasePruner() {
	log.Infof("Database: Pruning old releases from database...")
	oldestDate := time.Now().UTC().Truncate(time.Duration(72) * time.Hour)
	rowsCleared := database.DB.Unscoped().Delete(models.PushedRelease{},
		"created_at < ?", oldestDate).RowsAffected
	log.Infof("Database: Pruned %d releases from before %s", rowsCleared, oldestDate)
}
