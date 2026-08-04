package statistics_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/KyoshiBlame/TodoKy/internal/core/domain"
	core_logger "github.com/KyoshiBlame/TodoKy/internal/core/logger"
	core_http_request "github.com/KyoshiBlame/TodoKy/internal/core/transport/http/request"
	core_http_response "github.com/KyoshiBlame/TodoKy/internal/core/transport/http/response"
)

type GetStatisticsResponse struct {
	TasksCreated               int      `json:"tasks_created" example:"10"`
	TasksCompleted             int      `json:"tasks_completed" example:"5"`
	TasksCompletedRate         *float64 `json:"tasks_completed_rate" example:"50"`
	TasksAverageCompletionTime *string  `json:"tasks_avg_completion_time" example:"1m30s"`
}

// GetStatistics 	godoc
// @Summary 	Статистика задач
// @Description Получить статистику о существующий задачах в системе
// @Tags 		Statistics
// @Produce 	json
// Param		user_id query int 	 false "Фильтрация статистики по конкретному пользователю"
// Param		from 	query string false "Начало промежутка рассмотрения статистики (включительно), формат: YYYY-MM-DD"
// Param		to 		query string false "Конец промежутка рассмотрения статистики (не включительно), формат: YYYY-MM-DD"
// @Success 	200 	{object} GetStatisticsResponse 			  "Успешно полученная статистика"
// @Failure 	400 	{object} core_http_response.ErrorResponse "Bad request"
// @Failure 	500 	{object} core_http_response.ErrorResponse "Internal sercer error"
// @Router 		/statistics [get]
func (h *StatisticsHTTPHandler) GetStatistics(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseWriter := core_http_response.NewHTTPResponseHandler(log, rw)

	queryParams, err := getUserIDFromToQueryParams(r)
	if err != nil {
		responseWriter.ErrorResponse(
			err,
			"failed to get 'userID/from/to' from query param",
		)
		return
	}

	stat, err := h.statisticsService.GetStatistics(
		ctx,
		queryParams.userID,
		queryParams.from,
		queryParams.to,
	)
	if err != nil {
		responseWriter.ErrorResponse(
			err,
			"failed to get statistics",
		)
		return
	}

	response := toDTOFromDomain(stat)

	responseWriter.JSONResponse(
		response,
		http.StatusOK,
	)
}

func toDTOFromDomain(stat domain.Statistics) GetStatisticsResponse {
	var avgTime *string
	if stat.TasksAverageCompletionTime != nil {
		duration := stat.TasksAverageCompletionTime.String()
		avgTime = &duration
	}

	return GetStatisticsResponse{
		TasksCreated:               stat.TasksCreated,
		TasksCompleted:             stat.TasksCompleted,
		TasksCompletedRate:         stat.TasksCompletedRate,
		TasksAverageCompletionTime: avgTime,
	}
}

type queryParams struct {
	userID *int
	from   *time.Time
	to     *time.Time
}

func NewQueryParams(userID *int, from *time.Time, to *time.Time) queryParams {
	return queryParams{
		userID: userID,
		from:   from,
		to:     to,
	}
}

func getUserIDFromToQueryParams(r *http.Request) (queryParams, error) {

	const (
		userIDQueryParamKey = "user_id"
		FromQueryParamKey   = "from"
		ToQueryParamKey     = "to"
	)

	userID, err := core_http_request.GetIntQueryParam(r, userIDQueryParamKey)
	if err != nil {
		return queryParams{}, fmt.Errorf("get 'user_id' from query param")
	}

	from, err := core_http_request.GetDateQueryParam(r, FromQueryParamKey)
	if err != nil {
		return queryParams{}, fmt.Errorf("failed to get 'from' from query param")
	}

	to, err := core_http_request.GetDateQueryParam(r, ToQueryParamKey)
	if err != nil {
		return queryParams{}, fmt.Errorf("failed to get 'to' from query param")
	}

	return NewQueryParams(
		userID,
		from,
		to,
	), nil
}
