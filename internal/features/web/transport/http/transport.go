package web_transport_http

import (
	"net/http"

	"github.com/darya-kznv/todo-service/internal/core/domain"
	core_http_server "github.com/darya-kznv/todo-service/internal/core/transport/http/server"
)

type WebService interface {
	GetIndex() (domain.File, error)
}

type WebHTTPHandler struct {
	webService WebService
}

func NewWebHTTPHandler(webService WebService) *WebHTTPHandler {
	return &WebHTTPHandler{webService: webService}
}

func (h *WebHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{{
		Method:  http.MethodGet,
		Path:    "/{$}",
		Handler: h.GetIndex,
	}}
}
