package task_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/KyoshiBlame/TodoKy/internal/core/domain"
	core_errors "github.com/KyoshiBlame/TodoKy/internal/core/errors"
	core_postgres_pool "github.com/KyoshiBlame/TodoKy/internal/core/repository/postgres/pool"
)

func (r *TasksRepository) PatchTask(
	ctx context.Context,
	taskID int,
	task domain.Task,
) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE todoky.task
	SET
		title=$1,
		description=$2,
		completed=$3,
		completed_at=$4,
		version=version+1
	WHERE id=$5 AND version=$6
	RETURNING
		id,
		version,
		title,
		description;
		completed;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		task.Titile,
		task.Description,
		task.Completed,
		task.CompletedAt,
	)

	var taskModel TaskModel

	err := row.Scan(
		&taskModel.ID,
		&taskModel.Version,
		&taskModel.Title,
		&taskModel.Description,
		&taskModel.Completed,
		&taskModel.CompletedAt,
	)

	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Task{}, fmt.Errorf(
				"task with id='%d' concurrently accessed: %w",
				taskID,
				core_errors.ErrConflict,
			)
		}
		return domain.Task{}, fmt.Errorf("scan error: %w", err)
	}

	taskDomain := domain.NewTask(
		taskModel.ID,
		taskModel.Version,
		taskModel.Title,
		taskModel.Description,
		taskModel.Completed,
		taskModel.CreateAt,
		taskModel.CompletedAt,
		taskModel.AuthorUserID,
	)

	return taskDomain, nil

}
