package tasks_transport_http

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/darya-kznv/todo-service/internal/core/domain"
	core_logger "github.com/darya-kznv/todo-service/internal/core/logger"
	core_http_request "github.com/darya-kznv/todo-service/internal/core/transport/http/request"
	core_http_response "github.com/darya-kznv/todo-service/internal/core/transport/http/response"
)

type CreateTaskRequest struct {
	Title        string    `json:"title" validate:"required,min=1,max=100" example:"Write API documentation"`
	Description  *string   `json:"description" validate:"omitempty,min=1,max=1000" example:"Document every endpoint"`
	AuthorUserID uuid.UUID `json:"author_user_id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
}

type CreateTaskResponse TaskDTOResponse

// CreateTask godoc
// @Summary Создать задачу
// @Description Создать новую задачу
// @Tags tasks
// @Accept json
// @Produce json
// @Param request body CreateTaskRequest true "Тело запроса"
// @Success 201 {object} CreateTaskResponse "Созданная задача"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 404 {object} core_http_response.ErrorResponse "Not found"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /tasks [post]
func (h *TasksHTTPHandler) CreateTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var request CreateTaskRequest

	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")

		return
	}

	taskDomain := domain.NewTaskUninitialized(
		request.Title,
		request.Description,
		request.AuthorUserID,
	)

	taskDomain, err := h.tasksService.CreateTask(ctx, taskDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create task")

		return
	}

	response := taskDTOFromDomain(taskDomain)
	responseHandler.JSONResponse(response, http.StatusCreated)
}
