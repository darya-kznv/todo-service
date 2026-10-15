package tasks_postgres_repository

import (
	"time"

	"github.com/google/uuid"
)

type TaskModel struct {
	ID           uuid.UUID
	Version      int
	Title        string
	Description  *string
	Completed    bool
	CompletedAt  *time.Time
	CreatedAt    time.Time
	AuthorUserID uuid.UUID
}
