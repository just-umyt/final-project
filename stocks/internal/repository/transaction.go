package repository

import (
	"context"
	"errors"
	"stocks/pkg/logger"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgTxManager struct {
	pool *pgxpool.Pool
}

type PgTxManagerInterface interface {
	WithTx(ctx context.Context, fn func(StockRepoInterface) error) error
}

func NewPgTxManager(pool *pgxpool.Pool) *PgTxManager {
	return &PgTxManager{pool: pool}
}

func (tm *PgTxManager) WithTx(ctx context.Context, fn func(StockRepoInterface) error) error {
	tx, err := tm.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			logger.Log.Error(err)
		}
	}()

	factory := NewRepository(tx)

	if err := fn(factory); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
