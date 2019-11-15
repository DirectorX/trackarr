package loghook

import (
	"github.com/l3uddz/trackarr/logger"
	"github.com/l3uddz/trackarr/ws"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.uber.org/atomic"
	"strings"
	"time"
)

/* Var */
var (
	log = logger.GetLogger("loghook")
)

/* Struct */
type Loghooker struct {
	running *atomic.Bool
	hooked  *atomic.Bool
	queue   chan *logrus.Entry
}

type WebsocketLogMessage struct {
	Time      string `json:"time"`
	Level     string `json:"level"`
	Component string `json:"component"`
	Message   string `json:"message"`
}

/* Public */
func NewLoghooker() *Loghooker {
	return &Loghooker{
		running: atomic.NewBool(false),
		hooked:  atomic.NewBool(false),
		queue:   make(chan *logrus.Entry, 128),
	}
}

func (l *Loghooker) Push(entry *logrus.Entry) error {
	select {
	case l.queue <- entry:
		break
	default:
		// dont log the error as it will just trigger another push that will fail
		return errors.New("failed adding log entry to queue as it was full")

	}

	return nil
}

func (l *Loghooker) Start() error {
	if l.running.Load() {
		return errors.New("loghooker has already been started")
	}

	if l.queue == nil {
		l.queue = make(chan *logrus.Entry, 128)
	}

	go l.processor()

	if !l.hooked.Load() {
		logrus.AddHook(l)
		l.hooked.Store(true)
	}

	l.running.Store(true)
	return nil
}

func (l *Loghooker) Stop() error {
	if !l.running.Load() {
		return errors.New("loghooker has not been started")
	}

	close(l.queue)
	l.queue = nil

	l.running.Store(false)
	return nil
}

/* Private */
func (l *Loghooker) processor() {
	for {
		// pop log from queue
		entry, ok := <-l.queue
		if !ok {
			break
		}

		// get component from log entry
		component := ""
		if prefixValue, ok := entry.Data["prefix"]; ok {
			component = prefixValue.(string)
		}

		// create websocket message
		logMessage := &ws.WebsocketMessage{
			Type: "log",
			Data: &WebsocketLogMessage{
				Time:      entry.Time.Format(time.RFC3339),
				Level:     entry.Level.String(),
				Component: strings.TrimSpace(component),
				Message:   entry.Message,
			},
		}

		// broadcast hooked log message
		jsonData, err := logMessage.ToJsonString()
		if err == nil {
			ws.BroadcastTopic("logs", jsonData)
		}
	}

	log.Info("Logs queue processor finished")
}
