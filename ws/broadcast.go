package ws

func BroadcastAll(data string) {
	_ = m.Broadcast([]byte(data))
}

func BroadcastTopic(topic string, data string) {
	mtx.Lock()
	defer mtx.Unlock()

	sockets, ok := topics[topic]
	if !ok {
		return
	}

	for socket := range sockets {
		_ = socket.Write([]byte(data))
	}
}
