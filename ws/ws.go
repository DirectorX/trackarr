package ws

import (
	"sync"

	jsoniter "github.com/json-iterator/go"
	"gitlab.com/cloudb0x/trackarr/logger"
	"github.com/olahol/melody"
)

var (
	log  = logger.GetLogger("websocket")
	json = jsoniter.ConfigCompatibleWithStandardLibrary

	// core
	m *melody.Melody

	sockets map[*melody.Session]map[string]bool
	topics  map[string]map[*melody.Session]bool
	mtx     *sync.Mutex

	readCallbacks map[string][]interface{}
)

/* Public */

func Init() error {
	// init wrapper
	m = melody.New()

	sockets = make(map[*melody.Session]map[string]bool)
	topics = make(map[string]map[*melody.Session]bool)
	mtx = &sync.Mutex{}

	m.HandleConnect(socketConnected)
	m.HandleDisconnect(socketDisconnected)
	m.HandleMessage(socketMessage)

	// init default read callbacks
	readCallbacks = make(map[string][]interface{})

	AddReadCallback("subscribe", callbackSubscribe)
	AddReadCallback("unsubscribe", callbackUnsubscribe)

	return nil
}
