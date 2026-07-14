package core_pgx_pool

import (
	"context"
	"fmt"
	"time"

	core_postgres_pool "github.com/KyoshiBlame/TodoKy/internal/core/repository/postgres/pool"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool struct {
	*pgxpool.Pool
	opTimeout time.Duration
}

func NewPool(
	ctx context.Context,
	config Config,
) (*Pool, error) {

	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	pgxConfig, err := pgxpool.ParseConfig(connectionString)

	if err != nil {
		return nil, fmt.Errorf("parse pgxConfig: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)

	if err != nil {
		return nil, fmt.Errorf("create new pgx pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("pgxPool ping: %w", err)
	}

	return &Pool{Pool: pool, opTimeout: config.Timeout}, nil
}

func (p *Pool) Query(
	ctx context.Context,
	sql string,
	args ...any,
) (core_postgres_pool.Rows, error) {
	return p.Pool.Query(ctx, sql, args...)
}

func (p *Pool) QueryRow(ctx context.Context,
	sql string,
	args ...any,
) core_postgres_pool.Row {
	return p.Pool.QueryRow(ctx, sql, args...)
}

func (p *Pool) Exec(
	ctx context.Context,
	sql string,
	args ...any,
) (core_postgres_pool.CommandTag, error) {
	return p.Pool.Exec(ctx, sql, args...)
}

func (p *Pool) OpTimeout() time.Duration {
	return p.opTimeout
}
