package database

import (
	"context"
	"crud-task/pkg"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"sync"
)

var (
	dbPool *pgxpool.Pool
	once   sync.Once
)

func Connect() (*pgxpool.Pool, error) {
	var initError error
	once.Do(func() {
		dsn := pkg.MakeDSN()
		var err error
		pool, err := pgxpool.New(context.Background(), dsn)

		if err != nil {
			initError = fmt.Errorf("create db pool: %w", err)
			return
		}

		if err = pool.Ping(context.Background()); err != nil {
			initError = fmt.Errorf("ping db: %w", err)
			pool.Close()
			return
		}
		dbPool = pool
	})

	return dbPool, initError
}

func Close() {
	if dbPool != nil {
		dbPool.Close()
	}
}
