package repository

import (
	"context"
	"fmt"
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

func (b *BaseRepository) WithTx(ctx context.Context) (pgx.Tx, error) {
	tx, err := b.DBPool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	return tx, nil
}
