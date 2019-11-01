package web

import (
	"encoding/json"
	"github.com/desertbit/glue"
	"github.com/labstack/echo"
)

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

type Wrapper struct {
	Context echo.Context
	Server  *glue.Server
}

func GlueWrapper() *Wrapper {
	// setup server
	socketServer = glue.NewServer()
	socketServer.OnNewSocket(onNewSocket)

	return &Wrapper{
		Server: socketServer,
	}
}

func (s *Wrapper) HandlerFunc(context echo.Context) error {
	s.Context = context
	s.Server.ServeHTTP(context.Response(), context.Request())
	return nil
}

func onNewSocket(s *glue.Socket) {
	log.Debugf("Socket connected: %s", s.RemoteAddr())

	s.OnClose(func() {
		log.Debugf("Socket closed: %s", s.RemoteAddr())
	})

	s.OnRead(func(data string) {
		// do nothing with read data - eventually we will have a parser
	})
}
