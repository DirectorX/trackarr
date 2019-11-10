package ws

/* Message */

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

/* Alert */

type WebsocketAlert struct {
	Level   string `json:"level"`
	Title   string `json:"title"`
	Message string `json:"msg"`
}

func NewAlert(level string, title string, message string) (string, error) {
	whMsg := &WebsocketMessage{
		Type: "alert",
		Data: &WebsocketAlert{
			Level:   level,
			Title:   title,
			Message: message,
		},
	}
	return whMsg.ToJsonString()
}
