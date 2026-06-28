package repository

import (
	"context"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var dialect = goqu.Dialect("postgres")

type DBTX interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type BaseRepository struct {
	DBPool *pgxpool.Pool
}

func NewBaseRepository(pool *pgxpool.Pool) *BaseRepository {
	return &BaseRepository{DBPool: pool}
}

func (base *BaseRepository) RunInTx(context context.Context, callback func(db DBTX) error) error {
	tx, err := base.DBPool.Begin(context)
	if err != nil {
		return MapError(err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(context)
			panic(p)
		} else if err != nil {
			_ = tx.Rollback(context)
		} else {
			err = tx.Commit(context)
		}
	}()

	err = callback(tx)
	return MapError(err)
}
