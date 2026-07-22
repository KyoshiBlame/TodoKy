package task_service

import (
	"context"
	"fmt"

	"github.com/KyoshiBlame/TodoKy/internal/core/domain"
)

func (s *TasksService) PatchTask(
	ctx context.Context,
	taskID int,
	taskDomain domain.TaskPatch,
) (domain.Task, error) {
	task, err := s.tasksRepository.GetTask(ctx, taskID)
	if err != nil {
		return domain.Task{}, fmt.Errorf("failed to get task from repository: %w", err)
	}

	if err := task.ApplyPatch(taskDomain); err != nil {
		return domain.Task{}, fmt.Errorf("error to update patch: %w", err)
	}

	taskPatched, err := s.tasksRepository.PatchTask(ctx, taskID, task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("failed to update patch task: %w", err)
	}

	return taskPatched, err
}
