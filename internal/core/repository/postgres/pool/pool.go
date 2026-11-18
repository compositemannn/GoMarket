package core_postgres_pool

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row                     // одна строка
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)              // много строк
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) // результат не строки (UPDATE)
	Close()

	OpTimeout() time.Duration // возвращает продолжительность таймаута
}

type ConnectionPool struct {
	*pgxpool.Pool
	opTimeout time.Duration // время на 1 sql-запрос
}

func NewConnectionPool(
	ctx context.Context,
	config Config,
) (*ConnectionPool, error) {
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
		return nil, fmt.Errorf("parse pgxconfig: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("create pgxpool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil { // можно ли установить соединение - проверка
		return nil, fmt.Errorf("ping pgxpool: %w", err)
	}

	return &ConnectionPool{
		Pool:      pool,
		opTimeout: config.Timeout,
	}, nil
}

func (p *ConnectionPool) OpTimeout() time.Duration {
	return p.opTimeout
}
