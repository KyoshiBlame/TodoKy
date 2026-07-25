package statistics_transport_http

import (
	"context"
	"net/http"
	"time"

	"github.com/KyoshiBlame/TodoKy/internal/core/domain"
	core_http_server "github.com/KyoshiBlame/TodoKy/internal/core/transport/http/server"
)

type StatisticsHTTPHandler struct {
	statisticsService StatisticsService
}

type StatisticsService interface {
	GetStatistics(
		ctx context.Context,
		userID *int,
		from *time.Time,
		to *time.Time,
	) (domain.Statistics, error)
}

func NewStatiscticsHTTPHandler(
	staticService StatisticsService,
) *StatisticsHTTPHandler {
	return &StatisticsHTTPHandler{
		statisticsService: staticService,
	}
}

func (h *StatisticsHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/statistics",
			Handler: h.GetStatistics,
		},
	}
}
