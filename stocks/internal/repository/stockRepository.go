package repository

import (
	"context"
	"stocks/internal/models"

	"github.com/jackc/pgx/v5"
)

type StockRepoInterface interface {
	GetSkuStockBySkuId(ctx context.Context, skuId models.SKUID) (models.SKU, models.Stock, error)
	AddStock(ctx context.Context, newStock models.Stock) error
	UpdateStock(ctx context.Context, item models.Stock) error
	DeleteStock(ctx context.Context, skuId models.SKUID, userId models.UserID) (int64, error)
	GetStocksByLocation(ctx context.Context, parameter GetSkuByLocationParameter) ([]models.FullStock, error)
}

type StockRepo struct {
	tx pgx.Tx
}

func NewStockRepository(tx pgx.Tx) *StockRepo {
	return &StockRepo{tx: tx}
}

func (r *StockRepo) GetSkuStockBySkuId(ctx context.Context, skuId models.SKUID) (models.SKU, models.Stock, error) {
	query := `SELECT * FROM sku l LEFT JOIN stock r ON r.sku_id = l.sku_id WHERE l.sku_id = $1`

	var repoSku SKU
	var repoStock Stock

	var sku models.SKU
	var stock models.Stock

	err := r.tx.QueryRow(ctx, query, skuId).Scan(&repoSku.SkuId, &repoSku.Name, &repoSku.Type, &repoStock.Id, &repoStock.SkuId, &repoStock.Price, &repoStock.Location, &repoStock.Count, &repoStock.UserId)
	if err != nil {
		return sku, stock, err
	}

	sku = models.SKU{
		SkuID: repoSku.SkuId,
		Name:  repoSku.Name,
		Type:  repoSku.Type,
	}

	if repoStock.UserId == nil {
		return sku, stock, nil
	}

	stock.Id = *repoStock.Id
	stock.SkuId = *repoStock.SkuId
	stock.Count = *repoStock.Count
	stock.Price = *repoStock.Price
	stock.Count = *repoStock.Count
	stock.Location = *repoStock.Location
	stock.UserId = *repoStock.UserId

	return sku, stock, nil
}

func (r *StockRepo) AddStock(ctx context.Context, newStock models.Stock) error {
	query := `INSERT INTO stock (price, location, count, user_id, sku_id) VALUES ($1, $2, $3, $4, $5)`

	_, err := r.tx.Exec(ctx, query, newStock.Price, newStock.Location, newStock.Count, newStock.UserId, newStock.SkuId)
	if err != nil {
		return err
	}

	return nil
}

func (r *StockRepo) UpdateStock(ctx context.Context, item models.Stock) error {
	query := `UPDATE stock SET price = $1, location = $2, count = $3 WHERE sku_id = $4`

	_, err := r.tx.Exec(ctx, query, item.Price, item.Location, item.Count, item.SkuId)
	if err != nil {
		return err
	}

	return nil
}

func (r *StockRepo) DeleteStock(ctx context.Context, skuId models.SKUID, userId models.UserID) (int64, error) {
	query := `DELETE FROM stock WHERE sku_id = $1 AND user_id = $2`

	row, err := r.tx.Exec(ctx, query, skuId, userId)

	return row.RowsAffected(), err
}

func (r *StockRepo) GetStocksByLocation(ctx context.Context, parameter GetSkuByLocationParameter) ([]models.FullStock, error) {
	var items []models.FullStock
	var err error

	query := `SELECT * FROM sku l INNER JOIN stock r ON r.sku_id = l.sku_id WHERE r.location = $1 AND r.user_id = $2 LIMIT $3 OFFSET $4`

	rows, err := r.tx.Query(ctx, query, parameter.Location, parameter.User_id, parameter.Limit, parameter.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var repoSku SKU
		var repoStock Stock

		err = rows.Scan(&repoSku.SkuId, &repoSku.Name, &repoSku.Type, &repoStock.Id, &repoStock.SkuId, &repoStock.Price, &repoStock.Location, &repoStock.Count, &repoStock.UserId)
		if err != nil {
			return nil, err
		}

		fullStock := models.FullStock{
			SKU: models.SKU{
				SkuID: repoSku.SkuId,
				Name:  repoSku.Name,
				Type:  repoSku.Type,
			},
			Stock: models.Stock{
				Id:       *repoStock.Id,
				Price:    *repoStock.Price,
				Location: *repoStock.Location,
				Count:    *repoStock.Count,
				UserId:   *repoStock.UserId,
			},
		}

		items = append(items, fullStock)
	}

	return items, nil
}
