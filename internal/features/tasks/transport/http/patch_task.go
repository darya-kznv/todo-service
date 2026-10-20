package tasks_transport_http

import (
	"fmt"
	"net/http"

	"github.com/darya-kznv/todo-service/internal/core/domain"
	core_errors "github.com/darya-kznv/todo-service/internal/core/errors"
	core_logger "github.com/darya-kznv/todo-service/internal/core/logger"
	core_http_request "github.com/darya-kznv/todo-service/internal/core/transport/http/request"
	core_http_response "github.com/darya-kznv/todo-service/internal/core/transport/http/response"
	core_http_types "github.com/darya-kznv/todo-service/internal/core/transport/http/types"
)

type PatchTaskRequest struct {
	Title       core_http_types.Nullable[string] `json:"title" swaggertype:"string" example:"Publish API documentation"`
	Description core_http_types.Nullable[string] `json:"description" swaggertype:"string" example:"Review and publish Swagger UI"`
	Completed   core_http_types.Nullable[bool]   `json:"completed" swaggertype:"boolean" example:"true"`
}

func (r *PatchTaskRequest) Validate() error {
	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf("'Title' can't be NULL: %w", core_errors.ErrInvalidArgument)
		}

		titleLength := len([]rune(*r.Title.Value))
		if titleLength < 1 || titleLength > 100 {
			return fmt.Errorf(
				"`Title` must be between 1 and 100 symbols: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if r.Description.Set && r.Description.Value != nil {
		descriptionLength := len([]rune(*r.Description.Value))
		if descriptionLength < 1 || descriptionLength > 1000 {
			return fmt.Errorf(
				"`Description` must be between 1 and 1000 symbols: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if r.Completed.Set && r.Completed.Value == nil {
		return fmt.Errorf("'Completed' can't be NULL: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}

type PatchTaskResponse TaskDTOResponse

// PatchTask godoc
// @Summary Обновить задачу
// @Description Частично обновить задачу по идентификатору
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "ID задачи" Format(uuid)
// @Param request body PatchTaskRequest true "Изменяемые поля"
// @Success 200 {object} PatchTaskResponse "Обновлённая задача"
// @Failure 400 {object} core_http_response.ErrorResponse "Некорректный запрос"
// @Failure 404 {object} core_http_response.ErrorResponse "Задача не найдена"
// @Failure 409 {object} core_http_response.ErrorResponse "Конфликт версий"
// @Failure 500 {object} core_http_response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /tasks/{id} [patch]
func (h *TasksHTTPHandler) PatchTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	taskID, err := core_http_request.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get taskID path value")

		return
	}

	var request PatchTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")

		return
	}

	taskPatch := domain.NewTaskPatch(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		request.Completed.ToDomain(),
	)

	task, err := h.tasksService.PatchTask(ctx, taskID, taskPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch task")

		return
	}

	response := PatchTaskResponse(taskDTOFromDomain(task))
	responseHandler.JSONResponse(response, http.StatusOK)
}
