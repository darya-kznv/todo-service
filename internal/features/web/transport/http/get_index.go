package web_transport_http

import (
	"net/http"

	core_logger "github.com/darya-kznv/todo-service/internal/core/logger"
	core_http_response "github.com/darya-kznv/todo-service/internal/core/transport/http/response"
)

func (h *WebHTTPHandler) GetIndex(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	file, err := h.webService.GetIndex()
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get index page")

		return
	}

	responseHandler.HTMLResponse(file)
}
