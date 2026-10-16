package tasks_postgres_repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/darya-kznv/todo-service/internal/core/domain"
	"github.com/google/uuid"
)

const (
	defaultLimit  = 10
	defaultOffset = 0
)

func (r *TasksRepository) GetTasks(
	ctx context.Context,
	limit *int,
	offset *int,
	userID *uuid.UUID,
) ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	limitVal := defaultLimit
	if limit != nil {
		limitVal = *limit
	}

	offsetVal := defaultOffset
	if offset != nil {
		offsetVal = *offset
	}

	var query strings.Builder

	query.WriteString(`
	SELECT id, version, title, description, completed, created_at, completed_at, author_user_id
	FROM todoservice.tasks`)

	args := []any{limitVal, offsetVal}

	if userID != nil {
		query.WriteString(`
	WHERE author_user_id = $3`)
		args = append(args, *userID)
	}

	query.WriteString(`
	ORDER BY created_at DESC, id ASC
	LIMIT $1 OFFSET $2;`)

	rows, err := r.pool.Query(ctx, query.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("select tasks: %w", err)
	}
	defer rows.Close()

	var taskModels []TaskModel

	for rows.Next() {
		var taskModel TaskModel

		if err := rows.Scan(
			&taskModel.ID,
			&taskModel.Version,
			&taskModel.Title,
			&taskModel.Description,
			&taskModel.Completed,
			&taskModel.CreatedAt,
			&taskModel.CompletedAt,
			&taskModel.AuthorUserID,
		); err != nil {
			return nil, fmt.Errorf("scan tasks: %w", err)
		}

		taskModels = append(taskModels, taskModel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	return taskDomainsFromModels(taskModels), nil
}
