package statistics_postgres_repository

import core_postgres_pool "github.com/KyoshiBlame/TodoKy/internal/core/repository/postgres/pool"

type StaticsRepository struct {
	pool core_postgres_pool.Pool
}

func NewStatiscticsRepository(
	pool core_postgres_pool.Pool,
) *StaticsRepository {
	return &StaticsRepository{
		pool: pool,
	}
}
