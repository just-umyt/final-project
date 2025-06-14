package repository

import (
	"context"
	"stocks/internal/models"

	"github.com/jackc/pgx/v5"
)

type StockRepoInterface interface {
	AddItem(ctx context.Context, item models.SKU) error
	GetBySku(ctx context.Context, skuId models.SKUID) (*GetSKU, error)
}

type StockRepo struct {
	tx pgx.Tx
}

func NewRepository(tx pgx.Tx) *StockRepo {
	return &StockRepo{tx: tx}
}

func (r *StockRepo) AddItem(ctx context.Context, item models.SKU) error {
	query := `INSERT INTO sku (sku_id, name, type, price, location, count, user_id) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.tx.Exec(ctx, query, item.ID, item.Name, item.Type, item.Price, item.Location, item.Count, item.UserId)

	return err
}

func (r *StockRepo) GetBySku(ctx context.Context, skuId models.SKUID) (*GetSKU, error) {
	query := `SELECT name, price, count FROM sku WHERE sku_id = $1`

	var item GetSKU
	err := r.tx.QueryRow(ctx, query, skuId).Scan(&item.Name, &item.Price, &item.Count)
	if err != nil {
		return nil, err
	}

	return &item, nil
}
