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
	Title       core_http_types.Nullable[string] `json:"title" swaggertype:"string" example:"Сходить в магазин"`
	Description core_http_types.Nullable[string] `json:"description" swaggertype:"string" example:"Погулять с пушком в новом парке"`
	Completed   core_http_types.Nullable[bool]   `json:"completed" swaggertype:"boolean"`
}

type TaskPatchedResponse taskDTOResponse

func (r *PatchTaskRequest) Validate() error {
	if r.Title.Set {

		if r.Title.Value == nil {
			return fmt.Errorf("'Titile can't be NULL")
		}

		lenTitle := len([]rune(*r.Title.Value))
		if lenTitle < 1 || lenTitle > 100 {
			return fmt.Errorf("'Title' must be between 1 and 100")
		}

	}

	if r.Description.Set {
		if r.Description.Value != nil {
			descLen := len([]rune(*r.Description.Value))
			if descLen < 1 || descLen > 1000 {
				return fmt.Errorf("'Description' must be between 1 and 1000")
			}
		}
	}

	if r.Completed.Set {
		if r.Completed.Value == nil {
			return fmt.Errorf("completed cant be nill")
		}
	}

	return nil
}

// PatchTasks 	godoc
// @Summary 	Изменнение задачи
// @Description Изменить существующую задачу в системе
// @Description ### Логика обновления полей (Three-state logic):
// @Tags 		Tasks
// @Description 1. **Поле не передано**: `description` игнорируется, значение в БД не меняется
// @Description 2. **Явно передано значение**: `"description":"помыть лапки собаке"` - устанавливает новый номер в БД
// @Description 3. **Передано NULL**: `"description":null` - очистит поле в БД (set NULL)
// @Description Ограниечение: `title` не может быть выставлен как NULL
// @Accept 		json
// @Produce 	json
// @Param 		id path int true "ID изменяемой задачи"
// @Param 		request body PatchTaskRequest true "PatchTaskRequest тело запроса"
// @Success 	200 {object} TaskPatchedResponse "Успешное изменение задачи"
// @Failure 	400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 	409 {object} core_http_response.ErrorResponse "Conflict"
// @Failure 	404 {object} core_http_response.ErrorResponse "Not found"
// @Failure 	500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router 		/tasks/{id} [patch]
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
		return
	}

	var request PatchTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)
		return
	}

	taskPatch := TaskPatchFromRequest(request)

	taskDomain, err := h.tasksService.PatchTask(ctx, taskID, taskPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to path task",
		)
		return
	}

	response := TaskPatchedResponse(taskDTOFromDomain(taskDomain))

	responseHandler.JSONResponse(
		response,
		http.StatusOK,
	)

}

func TaskPatchFromRequest(request PatchTaskRequest) domain.TaskPatch {
	return domain.NewTaskPatch(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		request.Completed.ToDomain(),
	)
}
