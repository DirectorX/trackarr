package web

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"time"
)

/* Structs */
type WebsocketLogMessage struct {
	Time      string
	Level     string
	Component string
	Message   string
}

func (whlMessage WebsocketLogMessage) ToJsonBytes() ([]byte, error) {
	bs, err := json.MarshalIndent(whlMessage, "", "  ")
	return bs, err
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

	// create log message struct
	logMessage := &WebsocketLogMessage{
		Time:      fmt.Sprintf("%s", entry.Time.Format(time.RFC3339)),
		Level:     entry.Level.String(),
		Component: component,
		Message:   entry.Message,
	}
	logEmitter.Update(logMessage)
	return nil
}

/* Vars */
var upgrader = websocket.Upgrader{}

/* Public */

func WebsocketLogHandler(c echo.Context) error {
	// initialize websocket
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()
	log.Tracef("Log websocket connection from %s", c.RealIP())

	// initialize logs receiver
	logReceiver := make(chan interface{})
	logEmitter.Subscribe(logReceiver)

	for v := range logReceiver {
		// Retrieve log event
		logEvent := v.(*WebsocketLogMessage)

		// Write log event to websocket
		if whMessage, whErr := logEvent.ToJsonBytes(); whErr == nil {
			err := ws.WriteMessage(websocket.TextMessage, whMessage)
			if err != nil {
				break
			}
		} else {
			break
		}
	}

	log.Tracef("Log websocket disconnection from %s", c.RealIP())
	return nil
}
