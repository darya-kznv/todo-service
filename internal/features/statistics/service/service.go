package statistics_service

import (
	"context"
	"time"

	"github.com/darya-kznv/todo-service/internal/core/domain"
	"github.com/google/uuid"
)

type StatisticsService struct {
	statisticsRepository StatisticsRepository
}

type StatisticsRepository interface {
	GetStatistics(
		ctx context.Context,
		userID *uuid.UUID,
		from *time.Time,
		to *time.Time,
	) (domain.Statistics, error)
}

func NewStatisticsService(statisticsRepository StatisticsRepository) *StatisticsService {
	return &StatisticsService{statisticsRepository: statisticsRepository}
}
