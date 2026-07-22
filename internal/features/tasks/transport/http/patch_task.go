package tasks_transport_http

import (
	"fmt"
	"net/http"

	"github.com/KyoshiBlame/TodoKy/internal/core/domain"
	core_logger "github.com/KyoshiBlame/TodoKy/internal/core/logger"
	core_http_request "github.com/KyoshiBlame/TodoKy/internal/core/transport/http/request"
	core_http_response "github.com/KyoshiBlame/TodoKy/internal/core/transport/http/response"
	core_http_types "github.com/KyoshiBlame/TodoKy/internal/core/transport/http/types"
)

type PatchTaskRequest struct {
	title       core_http_types.Nullable[string] `json:"title"`
	description core_http_types.Nullable[string] `json:"description"`
	completed   core_http_types.Nullable[bool]   `json:"completed"`
}

type TaskPatchedResponse taskDTOResponse

func (r *PatchTaskRequest) Validate() error {
	if r.title.Set {

		if r.title.Value == nil {
			return fmt.Errorf("'Titile can't be NULL")
		}

		lenTitle := len([]rune(*r.title.Value))
		if lenTitle < 1 || lenTitle > 100 {
			return fmt.Errorf("'Title' must be between 1 and 100")
		}

	}

	if r.description.Set {
		descLen := len([]rune(*r.description.Value))
		if descLen < 1 || descLen > 1000 {
			return fmt.Errorf("'Description' must be between 1 and 1000")
		}
	}

	return nil
}

func (h *TaskHTTPHandler) PatchTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	taskID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get 'id' from path",
		)
	}

	var request PatchTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTp request",
		)
	}

	taskPatch := TaskPatchFromRequest(request)

	taskDomain, err := h.tasksService.PatchTask(ctx, taskID, taskPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to path task",
		)
	}

	response := TaskPatchedResponse(taskDTOFromDomain(taskDomain))

	responseHandler.JSONResponse(
		response,
		http.StatusOK,
	)

}

func TaskPatchFromRequest(request PatchTaskRequest) domain.TaskPatch {
	return domain.NewTaskPatch(
		request.title.ToDomain(),
		request.description.ToDomain(),
		request.completed.ToDomain(),
	)
}
