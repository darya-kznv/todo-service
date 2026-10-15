package domain

import (
	"fmt"
	"time"

	core_errors "github.com/darya-kznv/todo-service/internal/core/errors"
	"github.com/google/uuid"
)

type Task struct {
	ID           uuid.UUID
	Version      int
	Title        string
	Description  *string
	Completed    bool
	CompletedAt  *time.Time
	CreatedAt    time.Time
	AuthorUserID uuid.UUID
}

func NewTask(
	id uuid.UUID,
	version int,
	title string,
	description *string,
	completed bool,
	completedAt *time.Time,
	createdAt time.Time,
	authorUserID uuid.UUID,
) Task {
	return Task{
		ID:           id,
		Version:      version,
		Title:        title,
		Description:  description,
		Completed:    completed,
		CompletedAt:  completedAt,
		CreatedAt:    createdAt,
		AuthorUserID: authorUserID,
	}
}

func NewTaskUninitialized(
	title string,
	description *string,
	authorUserID uuid.UUID,
) Task {
	return NewTask(
		UninitializedID,
		UninitializedVersion,
		title,
		description,
		false,
		nil,
		UninitializedCreatedAt,
		authorUserID,
	)
}

func (t Task) Initialized(createdAt time.Time) (Task, error) {
	if t.ID != UninitializedID {
		return Task{}, fmt.Errorf(
			"task is already initialized: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	t.ID = uuid.New()
	t.CreatedAt = createdAt
	t.Completed = false
	t.CompletedAt = nil

	if err := t.Validate(); err != nil {
		return Task{}, fmt.Errorf("validate initialized task: %w", err)
	}

	return t, nil
}

func (t *Task) Validate() error {
	titleLength := len([]rune(t.Title))

	if titleLength < 1 || titleLength > 100 {
		return fmt.Errorf(
			"invalid `Title` len: %d, %w",
			titleLength,
			core_errors.ErrInvalidArgument,
		)
	}

	if t.Description != nil {
		descriptionLength := len([]rune(*t.Description))

		if descriptionLength < 1 || descriptionLength > 1000 {
			return fmt.Errorf(
				"invalid `Description` len: %d, %w",
				descriptionLength,
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if t.Completed {
		if t.CompletedAt == nil {
			return fmt.Errorf(
				"`CompletedAt` can't be `nil` if `Completed`==`true`: %w",
				core_errors.ErrInvalidArgument,
			)
		}

		if t.CompletedAt.Before(t.CreatedAt) {
			return fmt.Errorf(
				"`CompletedAt` can't be before if `CreatedAt`: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	} else {
		if t.CompletedAt != nil {
			return fmt.Errorf(
				"`CompletedAt` must be `nil` if `Completed`==`false`: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	return nil
}
