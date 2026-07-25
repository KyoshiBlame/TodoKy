package statistics_service

import (
	"context"
	"time"

	"github.com/KyoshiBlame/TodoKy/internal/core/domain"
)

type StatisticsService struct {
	staticsRepository StaticsRepository
}

type StaticsRepository interface {
	GetStatistics(
		ctx context.Context,
		UserID *int,
		from *time.Time,
		to *time.Time,
	) ([]domain.Task, error)
}

func NewStatiscticsService(
	statisticRepository StaticsRepository,
) *StatisticsService {
	return &StatisticsService{
		staticsRepository: statisticRepository,
	}
}
