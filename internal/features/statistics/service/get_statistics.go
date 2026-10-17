package statistics_service

import (
	"context"
	"fmt"
	"time"

	"github.com/darya-kznv/todo-service/internal/core/domain"
	core_errors "github.com/darya-kznv/todo-service/internal/core/errors"
	"github.com/google/uuid"
)

func (s *StatisticsService) GetStatistics(
	ctx context.Context,
	userID *uuid.UUID,
	from *time.Time,
	to *time.Time,
) (domain.Statistics, error) {
	if from != nil && to != nil {
		if to.Before(*from) || to.Equal(*from) {
			return domain.Statistics{}, fmt.Errorf(
				"`to` must be after `from`: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	tasks, err := s.statisticsRepository.GetTasks(ctx, userID, from, to)
	if err != nil {
		return domain.Statistics{}, fmt.Errorf("get tasks from repository: %w", err)
	}

	statistics := calcStatistics(tasks)

	return statistics, nil
}

func calcStatistics(tasks []domain.Task) domain.Statistics {
	tasksCreated := len(tasks)

	if tasksCreated == 0 {
		return domain.Statistics{
			TasksCreated:               0,
			TasksCompleted:             0,
			TasksCompletedRate:         nil,
			TasksAverageCompletionTime: nil,
		}
	}

	tasksCompleted := 0
	var totalDuration time.Duration
	for _, task := range tasks {
		if task.Completed {
			tasksCompleted++
		}

		complationDuration := task.CompletionDuration()
		if complationDuration != nil {
			totalDuration += *complationDuration
		}
	}

	tasksCompletedRate := float64(tasksCompleted) / float64(tasksCreated) * 100

	var averageCompletionTime *time.Duration
	if tasksCompleted > 0 && totalDuration != 0 {
		avg := totalDuration / time.Duration(tasksCompleted)

		averageCompletionTime = &avg
	}

	return domain.Statistics{
		TasksCreated:               tasksCreated,
		TasksCompleted:             tasksCompleted,
		TasksCompletedRate:         &tasksCompletedRate,
		TasksAverageCompletionTime: averageCompletionTime,
	}
}
