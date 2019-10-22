package database

import (
	"github.com/l3uddz/trackarr/database/models"
	"github.com/matryer/runner"
	"time"
)

var (
	pTask *runner.Task
)

func StartPruner() {
	pTask = runner.Go(prunerTask)
}

func StopPruner() {
	if pTask.Running() {
		pTask.Stop()
		select {
		case <-pTask.StopChan():
			log.Info("Stopped task: Pruner")
		case <-time.After(5 * time.Second):
			log.Warn("Failed stopping task: Pruner")
		}
	}
}

/* Private */

func prunerTask(shouldStop runner.S) error {
	intervalHours := 12

	log.Info("Started task: Pruner")

	for {
		// clear releases older than 3 days
		oldestDate := time.Now().UTC().Truncate(time.Duration(72) * time.Hour)
		rowsCleared := DB.Unscoped().Delete(models.PushedRelease{}, "created_at < ?", oldestDate).RowsAffected
		log.Infof("Pruned %d releases from before :%s", rowsCleared, oldestDate)

		// go into sleep cycle
		waitCount := (60 * 60) * intervalHours
		for waitedSeconds := 0; waitedSeconds < waitCount; waitedSeconds++ {
			if shouldStop() {
				return nil
			}
			time.Sleep(1 * time.Second)
		}
	}
}
