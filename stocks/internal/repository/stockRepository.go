package repository

import (
	"context"
	"fmt"
	"stocks/internal/models"

	"github.com/jackc/pgx/v5"
)

type StockRepoInterface interface {
	AddSkuRepo(ctx context.Context, item models.SKU) error
	GetSkuBySkuIdRepo(ctx context.Context, skuId models.SKUID) (*GetSKU, error)
	GetUserIdBySkuIdRepo(ctx context.Context, skuId models.SKUID) (*models.UserID, error)
	UpdateSkuBySkuIdRepo(ctx context.Context, newItem models.SKU) error
	DeleteSkuBySkuIdRepo(ctx context.Context, skuId models.SKUID, userId models.UserID) error
	GetSkusByLocationRepo(ctx context.Context, parameter GetSkuByLocationParameter) ([]GetSkuByLocation, error)
}

type StockRepo struct {
	tx pgx.Tx
}

func NewRepository(tx pgx.Tx) *StockRepo {
	return &StockRepo{tx: tx}
}

// Add new item to stock storage.
func (r *StockRepo) AddSkuRepo(ctx context.Context, item models.SKU) error {
	query := `INSERT INTO sku (sku_id, name, type, price, location, count, user_id) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.tx.Exec(ctx, query, item.SkuId, item.Name, item.Type, item.Price, item.Location, item.Count, item.UserId)

	return err
}

// Get user_id by sku(sku_id).
func (r *StockRepo) GetUserIdBySkuIdRepo(ctx context.Context, skuId models.SKUID) (*models.UserID, error) {
	query := `SELECT user_id FROM sku WHERE sku_id = $1`

	var userId models.UserID

	err := r.tx.QueryRow(ctx, query, skuId).Scan(&userId)
	if err != nil {
		return nil, err
	}

	return &userId, nil
}

// Get sku by sku(sku_id).
func (r *StockRepo) GetSkuBySkuIdRepo(ctx context.Context, skuId models.SKUID) (*GetSKU, error) {
	query := `SELECT sku_id, name, price, count, type, location FROM sku WHERE sku_id = $1`

	var item GetSKU

	err := r.tx.QueryRow(ctx, query, skuId).Scan(&item.Sku, &item.Name, &item.Price, &item.Count, &item.Type, &item.Location)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

// Update sku by sku(sku_id).
func (r *StockRepo) UpdateSkuBySkuIdRepo(ctx context.Context, newItem models.SKU) error {
	query := `UPDATE sku SET name = $1, type = $2, price = $3, location = $4, count = $5 WHERE sku_id = $6`

	_, err := r.tx.Exec(ctx, query, newItem.Name, newItem.Type, newItem.Price, newItem.Location, newItem.Count, newItem.SkuId)
	if err != nil {
		return err
	}

	return nil
}

// Delete sku by sku(sku_id).
func (r *StockRepo) DeleteSkuBySkuIdRepo(ctx context.Context, skuId models.SKUID, userId models.UserID) error {
	query := `DELETE FROM sku WHERE sku_id = $1 AND user_id = $2`

	row, err := r.tx.Exec(ctx, query, skuId, userId)
	if err != nil {
		return err
	}

	if row.RowsAffected() == 0 {
		return fmt.Errorf("there is no sku with sku_id = %d & user_id = %d", skuId, userId)
	}

	return nil
}

// Get skus by location.
func (r *StockRepo) GetSkusByLocationRepo(ctx context.Context, parameter GetSkuByLocationParameter) ([]GetSkuByLocation, error) {
	var items []GetSkuByLocation
	var err error

	query := `SELECT sku_id, name, type, price, location, count FROM sku WHERE location = $1 AND user_id = $2 ORDER BY sku_id LIMIT $3 OFFSET $4`

	rows, err := r.tx.Query(ctx, query, parameter.Location, parameter.User_id, parameter.Limit, parameter.Offset)
	if err != nil {
		return items, err
	}
	defer rows.Close()

	for rows.Next() {
		var item GetSkuByLocation
		err = rows.Scan(&item.Sku, &item.Name, &item.Type, &item.Price, &item.Location, &item.Count)
		items = append(items, item)
	}

	return items, err
}
