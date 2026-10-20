package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/darya-kznv/todo-service/internal/core/logger"
	core_http_request "github.com/darya-kznv/todo-service/internal/core/transport/http/request"
	core_http_response "github.com/darya-kznv/todo-service/internal/core/transport/http/response"
)

type GetTaskResponse TaskDTOResponse

// GetTask godoc
// @Summary Получить задачу
// @Description Получить задачу по идентификатору
// @Tags tasks
// @Produce json
// @Param id path string true "ID задачи" Format(uuid)
// @Success 200 {object} GetTaskResponse "Найденная задача"
// @Failure 400 {object} core_http_response.ErrorResponse "Некорректный идентификатор"
// @Failure 404 {object} core_http_response.ErrorResponse "Задача не найдена"
// @Failure 500 {object} core_http_response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /tasks/{id} [get]
func (h *TasksHTTPHandler) GetTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	taskID, err := core_http_request.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get taskID path value")

		return
	}

	task, err := h.tasksService.GetTask(ctx, taskID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get task")

		return
	}

	response := GetTaskResponse(taskDTOFromDomain(task))
	responseHandler.JSONResponse(response, http.StatusOK)
}
