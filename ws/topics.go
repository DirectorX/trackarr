package ws

import "github.com/olahol/melody"

/* Public - credits: https://github.com/mahmud-ridwan/tonesa/blob/master/hub/hub.go */

func Subscribe(s *melody.Session, topic string) {
	mtx.Lock()
	defer mtx.Unlock()

	_, ok := sockets[s]
	if !ok {
		sockets[s] = map[string]bool{}
	}
	sockets[s][topic] = true

	_, ok = topics[topic]
	if !ok {
		topics[topic] = map[*melody.Session]bool{}
	}
	topics[topic][s] = true
}

func Unsubscribe(s *melody.Session, topic string) {
	mtx.Lock()
	defer mtx.Unlock()

	_, ok := sockets[s]
	if ok {
		delete(topics[topic], s)
		if len(topics[topic]) == 0 {
			delete(topics, topic)
		}
		delete(sockets, s)
	}
}

func UnsubscribeAll(s *melody.Session) {
	mtx.Lock()
	defer mtx.Unlock()

	for t := range sockets[s] {
		delete(topics[t], s)
		if len(topics[t]) == 0 {
			delete(topics, t)
		}
	}
	delete(sockets, s)
}
