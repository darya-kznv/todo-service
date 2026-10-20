package tasks_transport_http

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"

	core_logger "github.com/darya-kznv/todo-service/internal/core/logger"
	core_http_request "github.com/darya-kznv/todo-service/internal/core/transport/http/request"
	core_http_response "github.com/darya-kznv/todo-service/internal/core/transport/http/response"
)

type GetTasksResponse []TaskDTOResponse

// GetTasks godoc
// @Summary Получить задачи
// @Description Получить список задач с фильтрацией и пагинацией
// @Tags tasks
// @Produce json
// @Param limit query int false "Максимальное количество записей"
// @Param offset query int false "Смещение"
// @Param user_id query string false "ID автора" Format(uuid)
// @Success 200 {array} TaskDTOResponse "Список задач"
// @Failure 400 {object} core_http_response.ErrorResponse "Некорректные параметры запроса"
// @Failure 500 {object} core_http_response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /tasks [get]
func (h *TasksHTTPHandler) GetTasks(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	limit, offset, userID, err := getTasksQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get tasks query params")

		return
	}

	taskDomains, err := h.tasksService.GetTasks(ctx, limit, offset, userID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get tasks")

		return
	}

	response := GetTasksResponse(tasksDTOFromDomains(taskDomains))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func getTasksQueryParams(r *http.Request) (*int, *int, *uuid.UUID, error) {
	const (
		limitQueryParamKey        = "limit"
		offsetQueryParamKey       = "offset"
		authorUserIDQueryParamKey = "user_id"
	)

	limit, err := core_http_request.GetInitQueryParams(r, limitQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}

	offset, err := core_http_request.GetInitQueryParams(r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}

	authorUserID, err := core_http_request.GetUUIDQueryParam(r, authorUserIDQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get '%s' query param: %w", authorUserIDQueryParamKey, err)
	}

	return limit, offset, authorUserID, nil
}
