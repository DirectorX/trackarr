package apis

type Response struct {
	Error   bool   `json:"error" xml:"error"`
	Message string `json:"message" xml:"message"`
}
