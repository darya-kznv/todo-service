package statistics_postgres_repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/darya-kznv/todo-service/internal/core/domain"
	"github.com/google/uuid"
)

func (r *StatisticsRepository) GetStatistics(
	ctx context.Context,
	userID *uuid.UUID,
	from *time.Time,
	to *time.Time,
) (domain.Statistics, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var queryBuilder strings.Builder
	queryBuilder.WriteString(`
	SELECT
		COUNT(*),
		COUNT(*) FILTER (WHERE completed),
		(
			COUNT(*) FILTER (WHERE completed) * 100.0
			/ NULLIF(COUNT(*), 0)
		)::DOUBLE PRECISION,
		(
			AVG(EXTRACT(EPOCH FROM completed_at - created_at))
			FILTER (WHERE completed AND completed_at IS NOT NULL)
		)::DOUBLE PRECISION
	FROM todoservice.tasks`)

	args := make([]any, 0, 3)
	conditions := make([]string, 0, 3)

	if userID != nil {
		conditions = append(conditions, fmt.Sprintf("author_user_id = $%d", len(args)+1))
		args = append(args, *userID)
	}

	if from != nil {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", len(args)+1))
		args = append(args, *from)
	}

	if to != nil {
		conditions = append(conditions, fmt.Sprintf("created_at < $%d", len(args)+1))
		args = append(args, *to)
	}

	if len(conditions) > 0 {
		queryBuilder.WriteString("\n\tWHERE ")
		queryBuilder.WriteString(strings.Join(conditions, " AND "))
	}

	queryBuilder.WriteString(";")

	var statisticsModel StatisticsModel

	err := r.pool.QueryRow(ctx, queryBuilder.String(), args...).Scan(
		&statisticsModel.TasksCreated,
		&statisticsModel.TasksCompleted,
		&statisticsModel.TasksCompletedRate,
		&statisticsModel.AverageCompletionTimeSecond,
	)
	if err != nil {
		return domain.Statistics{}, fmt.Errorf("scan statistics: %w", err)
	}

	return statisticsDomainFromModel(statisticsModel), nil
}
