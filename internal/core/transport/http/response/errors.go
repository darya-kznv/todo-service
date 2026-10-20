package core_http_response

type ErrorResponse struct {
	Message string `json:"message" example:"failed to process request"`
	Error   string `json:"error" example:"invalid argument"`
}
