package task_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/KyoshiBlame/TodoKy/internal/core/domain"
	core_errors "github.com/KyoshiBlame/TodoKy/internal/core/errors"
	core_postgres_pool "github.com/KyoshiBlame/TodoKy/internal/core/repository/postgres/pool"
)

func (r *TasksRepository) GetTask(
	ctx context.Context,
	taskID int,
) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, title, version, description, completed, created_at, completed_at, author_user_id
	FROM todoky.tasks
	WHERE id=$1;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		taskID,
	)

	var taskModel TaskModel

	err := row.Scan(
		&taskModel.ID,
		&taskModel.Title,
		&taskModel.Version,
		&taskModel.Description,
		&taskModel.Completed,
		&taskModel.CreateAt,
		&taskModel.CompletedAt,
		&taskModel.AuthorUserID,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Task{}, fmt.Errorf("task with id='%d':%w", taskID, core_errors.ErrNotFound)
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
