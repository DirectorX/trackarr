package web

import (
	"encoding/json"
	"github.com/desertbit/glue"
	"github.com/labstack/echo"
)

/* WebsocketMessage */

type WebsocketMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func (whMessage WebsocketMessage) ToJsonString() (string, error) {
	bs, err := json.Marshal(whMessage)
	if err != nil {
		return "{}", err
	}
	return string(bs), err
}

/* Public */

type GlueWrapper struct {
	Context       echo.Context
	Server        *glue.Server
	ReadCallbacks map[string][]interface{}
}

func NewWrapper() *GlueWrapper {
	// init wrapper
	w := &GlueWrapper{
		Server: glue.NewServer(),
	}

	w.ReadCallbacks = make(map[string][]interface{}, 0)

	// init server
	w.Server.OnNewSocket(w.newSocketCreated)
	return w
}

func (w *GlueWrapper) HandlerFunc(context echo.Context) error {
	w.Context = context
	w.Server.ServeHTTP(context.Response(), context.Request())
	return nil
}

func (w *GlueWrapper) Broadcast(data string) {
	for _, sock := range socketWrapper.Server.Sockets() {
		sock.Write(data)
	}
}

func (w *GlueWrapper) AddReadCallback(msgType string, callback interface{}) {
	w.ReadCallbacks[msgType] = append(w.ReadCallbacks[msgType], callback)
}

/* Private */

func (w GlueWrapper) newSocketCreated(s *glue.Socket) {
	log.Debugf("Socket connected: %s", s.RemoteAddr())

	s.OnClose(func() {
		log.Debugf("Socket closed: %s", s.RemoteAddr())
	})

	s.OnRead(func(data string) {
		// unmarshal data
		whMsg := &WebsocketMessage{}
		if err := json.Unmarshal([]byte(data), &whMsg); err != nil {
			log.WithError(err).Errorf("Failed unmarshalling data received on socket: %v", data)
			return
		}

		// callbacks registered for this message type?
		callbacks, ok := w.ReadCallbacks[whMsg.Type]
		if !ok {
			// there were no callbacks for this message type
			log.Warnf("No read callbacks found for socket message type: %v", whMsg.Type)
			return
		}

		// iterate callbacks
		for _, callback := range callbacks {
			// ensure callback is in expected format
			callFunc, ok := callback.(func(*glue.Socket, *WebsocketMessage))
			if !ok {
				log.Warnf("Failed type asserting read callback function for socket message type: %v", whMsg.Type)
				continue
			}

			// trigger the callback
			callFunc(s, whMsg)
		}
	})
}
