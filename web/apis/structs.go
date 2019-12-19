package apis

type ErrorResponse struct {
	Error   bool   `json:"error" xml:"error"`
	Message string `json:"message" xml:"message"`
}
