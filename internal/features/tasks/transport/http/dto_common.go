package tasks_transport_http

import (
	"time"

	"github.com/KyoshiBlame/TodoKy/internal/core/domain"
)

type TaskDTOResponse struct {
	ID           int        `json:"id" example:"3"`
	Version      int        `json:"version" example:"1"`
	Titile       string     `json:"title" example:"Сходить в магазин"`
	Description  *string    `json:"description" example:"Купить рис, курицу и молоко"`
	Completed    bool       `json:"completed" example:"false"`
	CreatedAt    time.Time  `json:"created_at" example:"04.08.2026 17:01:00"`
	CompletedAt  *time.Time `json:"completed_at" example:""`
	AuthorUserID int        `json:"author_user_id" example:"3"`
}

func taskDTOFromDomain(task domain.Task) TaskDTOResponse {
	return TaskDTOResponse{
		ID:           task.ID,
		Version:      task.Version,
		Titile:       task.Titile,
		Description:  task.Description,
		CreatedAt:    task.CreatedAt,
		Completed:    task.Completed,
		CompletedAt:  task.CompletedAt,
		AuthorUserID: task.AuthorUserID,
	}
}

func taskDTosFromDomains(tasks []domain.Task) []TaskDTOResponse {
	dtos := make([]TaskDTOResponse, len(tasks))

	for i, task := range tasks {
		dtos[i] = taskDTOFromDomain(task)
	}

	return dtos
}
