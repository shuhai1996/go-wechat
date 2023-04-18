package socket

type Response struct {
	Type string `json:"type"`
	Username string `json:"username"`
	Message string `json:"message"`
}