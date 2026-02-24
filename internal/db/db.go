package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool() (*pgxpool.Pool, error) {
	dsn := os.Getenv("GOOSE_DBSTRING")
	if dsn == "" {
		return nil, fmt.Errorf("postgres | GOOSE_DBSTRING is empty")
	}

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("postgres | failed to parse dsn: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("postgres | failed to create pool: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("postgres | database unreachable: %w", err)
	}

	return pool, nil
}
