package statistics_transport_http

import (
	"fmt"
	"net/http"
	"time"

	core_logger "github.com/darya-kznv/todo-service/internal/core/logger"
	core_http_request "github.com/darya-kznv/todo-service/internal/core/transport/http/request"
	core_http_response "github.com/darya-kznv/todo-service/internal/core/transport/http/response"
	"github.com/google/uuid"
)

const (
	userIDQueryParamKey = "user_id"
	fromQueryParamKey   = "from"
	toQueryParamKey     = "to"
)

type GetStatisticsResponse struct {
	TasksCreated               int      `json:"tasks_created" example:"10"`
	TasksCompleted             int      `json:"tasks_completed" example:"8"`
	TasksCompletedRate         *float64 `json:"tasks_completed_rate" example:"0.8"`
	TasksAverageCompletionTime *string  `json:"tasks_average_completion_time" example:"2h30m0s"`
}

// GetStatistics godoc
// @Summary Получить статистику
// @Description Получить статистику задач с фильтрацией по автору и периоду
// @Tags statistics
// @Produce json
// @Param user_id query string false "ID автора" Format(uuid)
// @Param from query string false "Начало периода" Format(date)
// @Param to query string false "Конец периода" Format(date)
// @Success 200 {object} GetStatisticsResponse "Статистика задач"
// @Failure 400 {object} core_http_response.ErrorResponse "Некорректные параметры запроса"
// @Failure 500 {object} core_http_response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /statistics [get]
func (h *StatisticsHTTPHandler) GetStatistics(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, from, to, err := getStatisticsQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get statistics query params")

		return
	}

	statistics, err := h.statisticsService.GetStatistics(ctx, userID, from, to)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get statistics")

		return
	}

	responseHandler.JSONResponse(toDTOFromDomain(statistics), http.StatusOK)
}

func getStatisticsQueryParams(r *http.Request) (*uuid.UUID, *time.Time, *time.Time, error) {
	userID, err := core_http_request.GetUUIDQueryParam(r, userIDQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get '%s' query param: %w", userIDQueryParamKey, err)
	}

	from, err := core_http_request.GetDateQueryParam(r, fromQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get '%s' query param: %w", fromQueryParamKey, err)
	}

	to, err := core_http_request.GetDateQueryParam(r, toQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get '%s' query param: %w", toQueryParamKey, err)
	}

	return userID, from, to, nil
}
