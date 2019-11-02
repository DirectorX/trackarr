package web

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
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
	logMessage := &WebsocketMessage{
		Type: "log",
		Data: &WebsocketLogMessage{
			Time:      fmt.Sprintf("%s", entry.Time.Format(time.RFC3339)),
			Level:     entry.Level.String(),
			Component: strings.TrimSpace(component),
			Message:   entry.Message,
		},
	}

	// broadcast hooked log message
	jsonData, err := logMessage.ToJsonString()
	if err == nil {
		socketWrapper.Broadcast(jsonData)
	}

	return nil
}
