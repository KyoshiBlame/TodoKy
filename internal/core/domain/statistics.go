package domain

import "time"

type Statistics struct {
	TasksCreated               int
	TasksCompleted             int
	TasksCompletedRate         *float64
	TasksAverageCompletionTime *time.Duration
}

func NewStatisctics(
	tasksCreated int,
	tasksCompleted int,
	tasksCompletedRate *float64,
	taskAvetageCompletionTime *time.Duration,
) Statistics {
	return Statistics{
		TasksCreated:               tasksCompleted,
		TasksCompleted:             tasksCompleted,
		TasksCompletedRate:         tasksCompletedRate,
		TasksAverageCompletionTime: taskAvetageCompletionTime,
	}
}
