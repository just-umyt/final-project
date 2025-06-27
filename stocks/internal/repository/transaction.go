package repository

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgTxManager struct {
	pool *pgxpool.Pool
}

type PgTxManagerInterface interface {
	WithTx(ctx context.Context, fn func(IStockRepo) error) error
}

func NewPgTxManager(pool *pgxpool.Pool) *PgTxManager {
	return &PgTxManager{pool: pool}
}

func (tm *PgTxManager) WithTx(ctx context.Context, fn func(IStockRepo) error) error {
	tx, err := tm.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			rollBackErr := tx.Rollback(ctx)
			if rollBackErr != nil {
				log.Println(rollBackErr)
			}
		}
	}()

	factory := NewStockRepository(tx)

	if err = fn(factory); err != nil {
		return err
	}

	err = tx.Commit(ctx)

	return err
}
