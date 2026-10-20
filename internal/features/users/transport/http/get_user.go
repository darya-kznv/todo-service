package users_transport_http

import (
	"net/http"

	core_logger "github.com/darya-kznv/todo-service/internal/core/logger"
	core_http_request "github.com/darya-kznv/todo-service/internal/core/transport/http/request"
	core_http_response "github.com/darya-kznv/todo-service/internal/core/transport/http/response"
)

type GetUserResponse UserDTOResponse

// GetUser godoc
// @Summary Получить пользователя
// @Description Получить пользователя по идентификатору
// @Tags users
// @Produce json
// @Param id path string true "ID пользователя" Format(uuid)
// @Success 200 {object} GetUserResponse "Найденный пользователь"
// @Failure 400 {object} core_http_response.ErrorResponse "Некорректный идентификатор"
// @Failure 404 {object} core_http_response.ErrorResponse "Пользователь не найден"
// @Failure 500 {object} core_http_response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /users/{id} [get]
func (h *UsersHTTPHandler) GetUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, err := core_http_request.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID path value",
		)

		return
	}

	user, err := h.usersService.GetUser(ctx, userID)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get user",
		)

		return
	}

	response := GetUserResponse(userDTOFromDomain(user))
	responseHandler.JSONResponse(response, http.StatusOK)
}
