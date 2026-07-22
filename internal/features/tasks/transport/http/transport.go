package tasks_transport_http

import (
	"context"
	"net/http"

	"github.com/KyoshiBlame/TodoKy/internal/core/domain"
	core_http_server "github.com/KyoshiBlame/TodoKy/internal/core/transport/http/server"
)

type TaskHTTPHandler struct {
	tasksService TaskService
}

type TaskService interface {
	CreateTask(
		ctx context.Context,
		task domain.Task,
	) (domain.Task, error)

	GetTasks(
		ctx context.Context,
		userID *int,
		limmit *int,
		offset *int,
	) ([]domain.Task, error)

	DeleteTask(
		ctx context.Context,
		id int,
	) error

	GetTask(
		ctx context.Context,
		taskID int,
	) (domain.Task, error)

	PatchTask(
		ctx context.Context,
		taskID int,
		patchTask domain.TaskPatch,
	) (domain.Task, error)
}

func NewTasksHTTPHandler(
	taskService TaskService,
) *TaskHTTPHandler {
	return &TaskHTTPHandler{
		tasksService: taskService,
	}
}

func (h *TaskHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/tasks",
			Handler: h.CreateTask,
		},
		{
			Method:  http.MethodGet,
			Path:    "/tasks",
			Handler: h.GetTasks,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/tasks/{id}",
			Handler: h.DeleteTask,
		},
		{
			Method:  http.MethodGet,
			Path:    "/tasks/{id}",
			Handler: h.GetTask,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/tasks/{id}",
			Handler: h.PatchTask,
		},
	}
}
