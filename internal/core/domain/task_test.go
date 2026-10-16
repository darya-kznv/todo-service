package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestTaskApplyPatchCompletesTask(t *testing.T) {
	createdAt := time.Date(2026, time.October, 16, 10, 0, 0, 0, time.UTC)
	completedAt := createdAt.Add(time.Hour)
	task := NewTask(uuid.New(), 1, "task", nil, false, nil, createdAt, uuid.New())
	completed := true

	err := task.ApplyPatch(
		NewTaskPatch(Nullable[string]{}, Nullable[string]{}, Nullable[bool]{Value: &completed, Set: true}),
		completedAt,
	)
	if err != nil {
		t.Fatalf("ApplyPatch() error = %v", err)
	}

	if !task.Completed {
		t.Fatal("ApplyPatch() did not complete task")
	}

	if task.CompletedAt == nil || !task.CompletedAt.Equal(completedAt) {
		t.Fatalf("ApplyPatch() CompletedAt = %v, want %v", task.CompletedAt, completedAt)
	}
}

func TestTaskApplyPatchReopensTask(t *testing.T) {
	createdAt := time.Date(2026, time.October, 16, 10, 0, 0, 0, time.UTC)
	completedAt := createdAt.Add(time.Hour)
	task := NewTask(uuid.New(), 1, "task", nil, true, &completedAt, createdAt, uuid.New())
	completed := false

	err := task.ApplyPatch(
		NewTaskPatch(Nullable[string]{}, Nullable[string]{}, Nullable[bool]{Value: &completed, Set: true}),
		completedAt.Add(time.Hour),
	)
	if err != nil {
		t.Fatalf("ApplyPatch() error = %v", err)
	}

	if task.Completed {
		t.Fatal("ApplyPatch() did not reopen task")
	}

	if task.CompletedAt != nil {
		t.Fatalf("ApplyPatch() CompletedAt = %v, want nil", task.CompletedAt)
	}
}

func TestTaskApplyPatchDoesNotMutateTaskOnValidationError(t *testing.T) {
	createdAt := time.Date(2026, time.October, 16, 10, 0, 0, 0, time.UTC)
	task := NewTask(uuid.New(), 1, "task", nil, false, nil, createdAt, uuid.New())
	original := task
	invalidTitle := ""

	err := task.ApplyPatch(
		NewTaskPatch(Nullable[string]{Value: &invalidTitle, Set: true}, Nullable[string]{}, Nullable[bool]{}),
		createdAt.Add(time.Hour),
	)
	if err == nil {
		t.Fatal("ApplyPatch() error = nil, want validation error")
	}

	if task != original {
		t.Fatalf("ApplyPatch() mutated task on error: got %+v, want %+v", task, original)
	}
}
