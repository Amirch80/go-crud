package database

import (
	"context"
	"crud-task/pkg"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"sync"
)

var (
	dbPool *pgxpool.Pool
	once   sync.Once
)

func Connect() *pgxpool.Pool {
	once.Do(func() {
		dsn := pkg.MakeDSN()
		var err error
		dbPool, err = pgxpool.New(context.Background(), dsn)

		if err != nil {
			log.Fatal(err)
		}

		if err = dbPool.Ping(context.Background()); err != nil {
			log.Fatal(err)
		}
	})

	return dbPool
}

func Close() {
	if dbPool != nil {
		dbPool.Close()
	}
}
