package web

import (
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/olahol/melody"
	"sync"
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

type SocketWrapper struct {
	// private
	m *melody.Melody

	sockets map[*melody.Session]map[string]bool
	topics  map[string]map[*melody.Session]bool
	mtx     *sync.Mutex

	readCallbacks map[string][]interface{}
}

func NewWrapper() *SocketWrapper {
	// init wrapper
	w := &SocketWrapper{
		m: melody.New(),
	}

	w.sockets = make(map[*melody.Session]map[string]bool, 0)
	w.topics = make(map[string]map[*melody.Session]bool, 0)
	w.mtx = &sync.Mutex{}

	w.m.HandleConnect(w.socketConnected)
	w.m.HandleDisconnect(w.socketDisconnected)
	w.m.HandleMessage(w.socketMessage)

	// init default read callbacks
	w.readCallbacks = make(map[string][]interface{}, 0)

	w.AddReadCallback("subscribe", w.callbackSubscribe)

	return w
}

// core

func (w *SocketWrapper) HandlerFunc(context echo.Context) error {
	return w.m.HandleRequest(context.Response().Writer, context.Request())
}

// broadcast

func (w *SocketWrapper) BroadcastAll(data string) {
	_ = w.m.Broadcast([]byte(data))
}

func (w *SocketWrapper) BroadcastTopic(topic string, data string) {
	w.mtx.Lock()
	defer w.mtx.Unlock()

	sockets, ok := w.topics[topic]
	if !ok {
		return
	}

	for socket, _ := range sockets {
		_ = socket.Write([]byte(data))
	}
}

func (w *SocketWrapper) Subscribe(s *melody.Session, topic string) {
	w.mtx.Lock()
	defer w.mtx.Unlock()

	_, ok := w.sockets[s]
	if !ok {
		w.sockets[s] = map[string]bool{}
	}
	w.sockets[s][topic] = true

	_, ok = w.topics[topic]
	if !ok {
		w.topics[topic] = map[*melody.Session]bool{}
	}
	w.topics[topic][s] = true
}

func (w *SocketWrapper) Unsubscribe(s *melody.Session, topic string) {
	w.mtx.Lock()
	defer w.mtx.Unlock()

	_, ok := w.sockets[s]
	if ok {
		delete(w.topics[topic], s)
		if len(w.topics[topic]) == 0 {
			delete(w.topics, topic)
		}
		delete(w.sockets, s)
	}
}

func (w *SocketWrapper) UnsubscribeAll(s *melody.Session) {
	w.mtx.Lock()
	defer w.mtx.Unlock()

	for t := range w.sockets[s] {
		delete(w.topics[t], s)
		if len(w.topics[t]) == 0 {
			delete(w.topics, t)
		}
	}
	delete(w.sockets, s)
}

func (w *SocketWrapper) AddReadCallback(msgType string, callback interface{}) {
	w.readCallbacks[msgType] = append(w.readCallbacks[msgType], callback)
}

// callbacks

func (w *SocketWrapper) callbackSubscribe(s *melody.Session, m *WebsocketMessage) {
	log.Tracef("Processing socket callback for %s: %#v", s.Request.RemoteAddr, m)

	topic, ok := m.Data.(string)
	if !ok {
		return
	}

	w.Subscribe(s, topic)
	log.Debugf("Socket %s subscribed to topic: %q", s.Request.RemoteAddr, topic)
}

/* Private */
func (w SocketWrapper) socketConnected(s *melody.Session) {
	log.Debugf("Socket connected: %s", s.Request.RemoteAddr)
}

func (w *SocketWrapper) socketDisconnected(s *melody.Session) {
	log.Debugf("Socket disconnected: %s", s.Request.RemoteAddr)
	w.UnsubscribeAll(s)
}

func (w *SocketWrapper) socketMessage(s *melody.Session, msg []byte) {
	// unmarshal data
	whMsg := &WebsocketMessage{}
	if err := json.Unmarshal(msg, &whMsg); err != nil {
		log.WithError(err).Errorf("Failed unmarshalling data received on socket: %#v", msg)
		return
	}

	// callbacks registered for this message type?
	callbacks, ok := w.readCallbacks[whMsg.Type]
	if !ok {
		// there were no callbacks for this message type
		log.Warnf("No read callbacks found for socket message type: %q", whMsg.Type)
		return
	}

	// iterate callbacks
	for _, callback := range callbacks {
		// ensure callback is in expected format
		callFunc, ok := callback.(func(*melody.Session, *WebsocketMessage))
		if !ok {
			log.Warnf("Failed type asserting read callback function for socket message type: %v", whMsg.Type)
			continue
		}

		// trigger the callback
		callFunc(s, whMsg)
	}
}
