package web

import (
	"strings"
	"time"

	"github.com/l3uddz/trackarr/ws"
	"github.com/sirupsen/logrus"
)

/* Structs */

type WebsocketLogMessage struct {
	Time      string `json:"time"`
	Level     string `json:"level"`
	Component string `json:"component"`
	Message   string `json:"message"`
}

/* Logrus Hook */
type WebsocketLogHook struct{}

func (hook *WebsocketLogHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
		logrus.TraceLevel,
	}
}

func (hook *WebsocketLogHook) Fire(entry *logrus.Entry) error {
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

	return nil
}
