package tasks

import (
	"time"

	"gitlab.com/cloudb0x/trackarr/config"
	"gitlab.com/cloudb0x/trackarr/logger"

	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
)

var (
	log = logger.GetLogger("tasks")
)

type Tasks struct {
	scheduler  *cron.Cron
	TaskPruner *TaskPruner
}

/* Public */

// New instance
func New() *Tasks {
	return &Tasks{
		TaskPruner: &TaskPruner{
			Cron:        "0 0,6,12,18 * * *",
			MaxAgeHours: time.Duration(config.Config.Database.MaxAgeHours),
		},
	}
}

func (t *Tasks) Init() error {
	t.scheduler = cron.New()

	/* database tasks */
	// - prune old releases
	if _, err := t.scheduler.AddFunc(t.TaskPruner.Cron, t.taskDatabasePruner); err != nil {
		return errors.Wrap(err, "failed initializing task: database pruner")
	}

	return nil
}

func (t *Tasks) Start() {
	t.scheduler.Start()
	log.Info("Started scheduler")
}

func (t *Tasks) Stop() {
	ctx := t.scheduler.Stop()
	select {
	case <-ctx.Done():
		log.Info("Stopped scheduler")
	case <-time.After(5 * time.Second):
		log.Warn("Timed out waiting for scheduled jobs to finish")
	}
}

func (t *Tasks) AddTask(expression string, task func()) (cron.EntryID, error) {
	return t.scheduler.AddFunc(expression, task)
}
