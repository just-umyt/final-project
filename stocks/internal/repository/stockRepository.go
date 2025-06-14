package repository

import (
	"context"
	"stocks/internal/models"

	"github.com/jackc/pgx/v5"
)

type StockRepoInterface interface {
	AddSkuRepo(ctx context.Context, item models.SKU) error
	GetSkuBySkuIdRepo(ctx context.Context, skuId models.SKUID) (*GetSKU, error)
	UpdateSkuBySkuIdRepo(ctx context.Context, newItem models.SKU) error
}

type StockRepo struct {
	tx pgx.Tx
}

func NewRepository(tx pgx.Tx) *StockRepo {
	return &StockRepo{tx: tx}
}

func (r *StockRepo) AddSkuRepo(ctx context.Context, item models.SKU) error {
	query := `INSERT INTO sku (sku_id, name, type, price, location, count, user_id) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.tx.Exec(ctx, query, item.SkuId, item.Name, item.Type, item.Price, item.Location, item.Count, item.UserId)

	return err
}

func (r *StockRepo) GetSkuBySkuIdRepo(ctx context.Context, skuId models.SKUID) (*GetSKU, error) {
	query := `SELECT name, price, count, type, user_id FROM sku WHERE sku_id = $1`

	var item GetSKU
	err := r.tx.QueryRow(ctx, query, skuId).Scan(&item.Name, &item.Price, &item.Count, &item.Type, &item.UserId)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *StockRepo) UpdateSkuBySkuIdRepo(ctx context.Context, newItem models.SKU) error {
	query := `UPDATE sku SET name = $1, type = $2, price = $3, location = $4, count = $5 WHERE sku_id = $6`
	// UPDATE sku SET name = 'hello', type = 'newType' WHERE sku_id = 1;
	_, err := r.tx.Exec(ctx, query, newItem.Name, newItem.Type, newItem.Price, newItem.Location, newItem.Count, newItem.SkuId)
	if err != nil {
		return err
	}

	return nil
}
