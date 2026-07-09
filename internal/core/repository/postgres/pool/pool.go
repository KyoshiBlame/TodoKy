package core_postgres_pool

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ConnectionPool struct {
	*pgxpool.Pool
	opTimeout time.Duration
}

func NewConnectionPool(
	ctx context.Context,
	config Config,
) (*ConnectionPool, error) {

	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
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

	if err := pool.Ping(); err != nil {
		return nil, fmt.Errorf("pgxPool ping: %w", err)
	}

	return &ConnectionPool{Pool: pool, OpTimeout: config.Timeout}, nil
}

func (p *ConnectionPool) OpTimeout() time.Duration {
	return p.opTimeout
}
