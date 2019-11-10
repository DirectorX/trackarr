package ws

import "github.com/olahol/melody"

/* Public */

func AddReadCallback(msgType string, callback interface{}) {
	readCallbacks[msgType] = append(readCallbacks[msgType], callback)
}

/* Private */

func callbackSubscribe(s *melody.Session, m *WebsocketMessage) {
	log.Tracef("Processing callback for %s: %#v", s.Request.RemoteAddr, m)

	topic, ok := m.Data.(string)
	if !ok {
		return
	}

	Subscribe(s, topic)
	log.Debugf("%s subscribed to topic: %q", s.Request.RemoteAddr, topic)
}

func callbackUnsubscribe(s *melody.Session, m *WebsocketMessage) {
	log.Tracef("Processing callback for %s: %#v", s.Request.RemoteAddr, m)

	topic, ok := m.Data.(string)
	if !ok {
		return
	}

	Unsubscribe(s, topic)
	log.Debugf("%s unsubscribed from topic: %q", s.Request.RemoteAddr, topic)
}
