package statistics_postgres_repository

import (
	"time"

	"github.com/darya-kznv/todo-service/internal/core/domain"
)

type StatisticsModel struct {
	TasksCreated                int64
	TasksCompleted              int64
	TasksCompletedRate          *float64
	AverageCompletionTimeSecond *float64
}

func statisticsDomainFromModel(statistics StatisticsModel) domain.Statistics {
	var averageCompletionTime *time.Duration
	if statistics.AverageCompletionTimeSecond != nil {
		duration := time.Duration(
			*statistics.AverageCompletionTimeSecond * float64(time.Second),
		)
		averageCompletionTime = &duration
	}

	return domain.NewStatistics(
		int(statistics.TasksCreated),
		int(statistics.TasksCompleted),
		statistics.TasksCompletedRate,
		averageCompletionTime,
	)
}
