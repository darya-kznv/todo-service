package domain

import "time"

type Statistics struct {
	TasksCreated               int
	TasksCompleted             int
	TasksCompletedRate         *float64
	TasksAverageCompletionTime *time.Duration
}

func NewStatistics(
	tasksCreated int,
	tasksCompleted int,
	tasksCompletedRate *float64,
	tasksAverageCompletionTime *time.Duration,
) Statistics {
	return Statistics{
		TasksCreated:               tasksCreated,
		TasksCompleted:             tasksCompleted,
		TasksCompletedRate:         tasksCompletedRate,
		TasksAverageCompletionTime: tasksAverageCompletionTime,
	}
}

func CreateStatistics(tasks []Task) Statistics {
	tasksCreated := len(tasks)
	if tasksCreated == 0 {
		return NewStatistics(0, 0, nil, nil)
	}

	tasksCompleted := 0
	var totalDuration time.Duration

	for _, task := range tasks {
		if task.Completed {
			tasksCompleted++
		}

		completionDuration := task.CompletionDuration()
		if completionDuration != nil {
			totalDuration += *completionDuration
		}
	}

	tasksCompletedRate := float64(tasksCompleted) / float64(tasksCreated) * 100

	var averageCompletionTime *time.Duration
	if tasksCompleted > 0 {
		average := totalDuration / time.Duration(tasksCompleted)
		averageCompletionTime = &average
	}

	return NewStatistics(
		tasksCreated,
		tasksCompleted,
		&tasksCompletedRate,
		averageCompletionTime,
	)
}

func (t *Task) CompletionDuration() *time.Duration {
	if !t.Completed {
		return nil
	}

	if t.CompletedAt == nil {
		return nil
	}

	duration := t.CompletedAt.Sub(t.CreatedAt)

	return &duration
}
