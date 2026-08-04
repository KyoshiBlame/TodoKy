package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/KyoshiBlame/TodoKy/internal/core/logger"
	core_http_request "github.com/KyoshiBlame/TodoKy/internal/core/transport/http/request"
	core_http_response "github.com/KyoshiBlame/TodoKy/internal/core/transport/http/response"
)

// DeleteTask 	godoc
// @Summary 	Удалить задачу
// @Description Удалить существующую задачу в системе
// @Tags 		Tasks
// @Param 		id path int true "ID задачи"
// @Success 	204 "Успешно удалена задача"
// @Failure 	404 {object} core_http_response.ErrorResponse "Not found"
// @Failure 	500 {object} core_http_response.ErrorResponse "Internal sercer error"
// @Failure 	404 {object} core_http_response.ErrorResponse "Bad request"
// @Router 		/tasks/{id} [delete]
func (h *TaskHTTPHandler) DeleteTask(rw http.ResponseWriter, r *http.Request) {
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

	if err := h.tasksService.DeleteTask(ctx, taskID); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to delete task",
		)
		return
	}

	responseHandler.NoContentResponse()
}
