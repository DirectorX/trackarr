package tasks

import (
	"github.com/l3uddz/trackarr/logger"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"time"
)

var (
	log       = logger.GetLogger("tasks")
	scheduler *cron.Cron
)

/* Public */

func Init() error {
	scheduler = cron.New()

	/* database tasks */

	// - prune old releases
	if _, err := scheduler.AddFunc(CronTaskDatabasePruner, taskDatabasePruner); err != nil {
		return errors.Wrap(err, "failed initializing task: database pruner")
	}

	return nil
}

func Start() {
	scheduler.Start()
	log.Info("Started scheduler")
}

func Stop() {
	ctx := scheduler.Stop()
	select {
	case <-ctx.Done():
		log.Info("Stopped scheduler")
	case <-time.After(5 * time.Second):
	}
}

func AddTask(expression string, task func()) (cron.EntryID, error) {
	return scheduler.AddFunc(expression, task)
}
