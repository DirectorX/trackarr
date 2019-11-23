package loghook

import (
	"strings"
	"sync"
	"time"

	"github.com/l3uddz/trackarr/logger"
	"github.com/l3uddz/trackarr/ws"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

/* Var */
var (
	log = logger.GetLogger("loghook")
)

/* Struct */
type Loghooker struct {
	running bool
	hooked  bool
	wg      sync.WaitGroup
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
		queue: make(chan *logrus.Entry, 128),
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
	if l.running {
		return errors.New("loghooker has already started")
	}

	go l.processor()

	if !l.hooked {
		logrus.AddHook(l)
		l.hooked = true
	}

	l.running = true
	return nil
}

func (l *Loghooker) Stop() error {
	if !l.running {
		return errors.New("loghooker has not started")
	}

	l.running = false
	l.wg.Wait()
	close(l.queue)

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
