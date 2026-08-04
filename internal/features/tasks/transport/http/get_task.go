package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/KyoshiBlame/TodoKy/internal/core/logger"
	core_http_request "github.com/KyoshiBlame/TodoKy/internal/core/transport/http/request"
	core_http_response "github.com/KyoshiBlame/TodoKy/internal/core/transport/http/response"
)

type taskDTOResponse TaskDTOResponse

// GetTask 	godoc
// @Summary 	Получить задачу
// @Description Получить задачу из системы
// @Tags 		Tasks
// @Produce 	json
// @Param 		id path int true "ID задачи"
// @Success 	200 {object} CreateTaskResponse "Успешное получение задачи"
// @Failure 	400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 	500 {object} core_http_response.ErrorResponse "Internal sercer error"
// @Router 		/tasks/{id} [get]
func (h *TaskHTTPHandler) GetTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	taskID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get 'id' from query param",
		)
		return
	}

	domainTask, err := h.tasksService.GetTask(ctx, taskID)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get task",
		)
		return
	}

	response := taskDTOResponse(taskDTOFromDomain(domainTask))

	responseHandler.JSONResponse(
		response,
		http.StatusOK,
	)

}
