package repository

import (
	"context"
	"errors"
	"stocks/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	getItemSKUquery    = `SELECT * FROM sku l LEFT JOIN stock r ON r.sku_id = l.sku_id WHERE l.sku_id = $1`
	addStockquery      = `INSERT INTO stock (price, location, count, user_id, sku_id) VALUES ($1, $2, $3, $4, $5)`
	updateStockquery   = `UPDATE stock SET price = $1, location = $2, count = $3 WHERE sku_id = $4`
	deleteStockquery   = `DELETE FROM stock WHERE sku_id = $1 AND user_id = $2`
	getItemsByLocquery = `SELECT * FROM sku l INNER JOIN stock r ON r.sku_id = l.sku_id WHERE r.location = $1 AND r.user_id = $2 LIMIT $3 OFFSET $4`
)

type IDBQuery interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

//go:generate mkdir -p mock
//go:generate minimock -o ./mock/ -s .go
type IStockRepo interface {
	GetItemBySKU(ctx context.Context, skuID models.SKUID) (models.Item, error)
	AddStock(ctx context.Context, stock models.Stock) error
	UpdateStock(ctx context.Context, stock models.Stock) error
	DeleteStock(ctx context.Context, skuID models.SKUID, userID models.UserID) error
	GetItemsByLocation(ctx context.Context, param GetStockByLocation) ([]models.Item, error)
}

type StockRepo struct {
	db IDBQuery
}

var ErrNotFound error = errors.New("not found")

func NewStockRepository(db IDBQuery) *StockRepo {
	return &StockRepo{db: db}
}

func (r *StockRepo) GetItemBySKU(ctx context.Context, skuID models.SKUID) (models.Item, error) {
	var sku SKU
	var stock Stock

	var item models.Item

	err := r.db.QueryRow(ctx, getItemSKUquery, skuID).Scan(&sku.ID, &sku.Name, &sku.Type, &stock.ID, &stock.SKUID, &stock.Price, &stock.Location, &stock.Count, &stock.UserID)
	if err != nil {
		return item, err
	}

	item.SKU = models.SKU{
		ID:   sku.ID,
		Name: sku.Name,
		Type: sku.Type,
	}

	if !stock.ID.Valid {
		return item, nil
	} else {
		item.Stock.ID = models.StockID(stock.ID.Int64)
		item.Stock.SKUID = models.SKUID(stock.SKUID.Uint32)
		item.Stock.Count = uint16(float32(stock.Count.Uint32))
		item.Stock.Price = stock.Price.Uint32
		item.Stock.Location = stock.Location.String
		item.Stock.UserID = models.UserID(stock.UserID.Int64)

		return item, nil
	}
}

func (r *StockRepo) AddStock(ctx context.Context, stock models.Stock) error {
	tag, err := r.db.Exec(ctx, addStockquery, stock.Price, stock.Location, stock.Count, stock.UserID, stock.SKUID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() < 1 {
		return ErrNotFound
	}

	return nil
}

func (r *StockRepo) UpdateStock(ctx context.Context, stock models.Stock) error {
	tag, err := r.db.Exec(ctx, updateStockquery, stock.Price, stock.Location, stock.Count, stock.SKUID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() < 1 {
		return ErrNotFound
	}

	return nil
}

func (r *StockRepo) DeleteStock(ctx context.Context, skuID models.SKUID, userID models.UserID) error {
	tag, err := r.db.Exec(ctx, deleteStockquery, skuID, userID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() < 1 {
		return ErrNotFound
	}

	return err
}

func (r *StockRepo) GetItemsByLocation(ctx context.Context, param GetStockByLocation) ([]models.Item, error) {
	var items []models.Item

	rows, err := r.db.Query(ctx, getItemsByLocquery, param.Location, param.UserID, param.Limit, param.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var sku SKU
		var stock Stock

		err = rows.Scan(&sku.ID, &sku.Name, &sku.Type, &stock.ID, &stock.SKUID, &stock.Price, &stock.Location, &stock.Count, &stock.UserID)
		if err != nil {
			return nil, err
		}

		item := models.Item{
			SKU: models.SKU{
				ID:   sku.ID,
				Name: sku.Name,
				Type: sku.Type,
			},
			Stock: models.Stock{
				ID:       models.StockID(stock.ID.Int64),
				Price:    stock.Price.Uint32,
				Location: stock.Location.String,
				Count:    uint16(float32(stock.Count.Uint32)),
				UserID:   models.UserID(stock.UserID.Int64),
			},
		}

		items = append(items, item)
	}

	return items, nil
}
