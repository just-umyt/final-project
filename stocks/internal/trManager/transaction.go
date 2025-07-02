package trManager

import (
	"context"
	"log"
	"stocks/internal/models"
	"stocks/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate mkdir -p mock
//go:generate minimock -o ./mock/ -s .go  -g
type IStockRepo interface {
	GetItemBySKU(ctx context.Context, skuID models.SKUID) (models.Item, error)
	AddStock(ctx context.Context, stock models.Stock) error
	UpdateStock(ctx context.Context, stock models.Stock) error
	DeleteStock(ctx context.Context, skuID models.SKUID, userID models.UserID) error
	GetItemsByLocation(ctx context.Context, param repository.GetStockByLocation) ([]models.Item, error)
}

//go:generate mkdir -p mock
//go:generate minimock -o ./mock -s .go  -g
type IPgTxManager interface {
	WithTx(ctx context.Context, fn func(IStockRepo) error) error
}

type PgTxManager struct {
	pool *pgxpool.Pool
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

	factory := repository.NewStockRepository(tx)

	if err = fn(factory); err != nil {
		return err
	}

	err = tx.Commit(ctx)

	return err
}
