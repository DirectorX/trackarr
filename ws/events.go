package ws

import (
	"github.com/olahol/melody"
)

func socketConnected(s *melody.Session) {
	log.Debugf("Connected: %s", s.Request.RemoteAddr)
}

func socketDisconnected(s *melody.Session) {
	log.Debugf("Disconnected: %s", s.Request.RemoteAddr)
	UnsubscribeAll(s)
}

func socketMessage(s *melody.Session, msg []byte) {
	// unmarshal data
	whMsg := &WebsocketMessage{}
	if err := json.Unmarshal(msg, &whMsg); err != nil {
		log.WithError(err).Errorf("Failed unmarshalling data received: %#v", msg)
		return
	}

	// callbacks registered for this message type?
	callbacks, ok := readCallbacks[whMsg.Type]
	if !ok {
		// there were no callbacks for this message type
		log.Warnf("No read callbacks found for message type: %q", whMsg.Type)
		return
	}

	// iterate callbacks
	for _, callback := range callbacks {
		// ensure callback is in expected format
		callFunc, ok := callback.(func(*melody.Session, *WebsocketMessage))
		if !ok {
			log.Warnf("Failed type asserting read callback function for message type: %v", whMsg.Type)
			continue
		}

		// trigger the callback
		callFunc(s, whMsg)
	}
}
