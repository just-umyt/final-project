package repository

import (
	"cart/internal/models"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type IDBQuery interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

type ICartRepo interface {
	GetCartInfoByUserID(ctx context.Context, userID models.UserID, skuID models.SKUID) (models.CartID, uint16, error)
	UpdateItemByUserID(ctx context.Context, cart models.Cart) error
	AddItem(ctx context.Context, cart models.Cart) error
	DeleteItem(ctx context.Context, userID models.UserID, skuID models.SKUID) error
	GetCartByUserID(ctx context.Context, userID models.UserID) ([]models.CartItem, error)
	ClearCartByUserID(ctx context.Context, userID models.UserID) error
}

type CartRepo struct {
	db IDBQuery
}

var ErrNotFound error = errors.New("not found")

func NewCartRepository(db IDBQuery) *CartRepo {
	return &CartRepo{db: db}
}

func (c *CartRepo) GetCartInfoByUserID(ctx context.Context, userID models.UserID, skuID models.SKUID) (models.CartID, uint16, error) {
	query := `SELECT id, count FROM cart WHERE user_id = $1 AND sku_id = $2`

	var cartID models.CartID
	var count uint16

	err := c.db.QueryRow(ctx, query, userID, skuID).Scan(&cartID, &count)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return 0, 0, err
	}

	return cartID, count, nil
}

func (c *CartRepo) UpdateItemByUserID(ctx context.Context, cart models.Cart) error {
	query := `UPDATE cart SET count = $1 WHERE user_id = $2 AND sku_id = $3`

	tag, err := c.db.Exec(ctx, query, cart.Count, cart.UserID, cart.SKUID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() < 1 {
		return ErrNotFound
	}

	return nil
}

func (c *CartRepo) AddItem(ctx context.Context, cart models.Cart) error {
	query := `INSERT INTO cart (user_id, sku_id, count) VALUES ($1, $2, $3)`

	tag, err := c.db.Exec(ctx, query, cart.UserID, cart.SKUID, cart.Count)
	if err != nil {
		return err
	}

	if tag.RowsAffected() < 1 {
		return ErrNotFound
	}

	return nil
}

func (c *CartRepo) DeleteItem(ctx context.Context, userID models.UserID, skuID models.SKUID) error {
	query := `DELETE FROM cart WHERE user_id = $1 AND sku_id = $2`

	tag, err := c.db.Exec(ctx, query, userID, skuID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() < 1 {
		return ErrNotFound
	}

	return nil
}

func (c *CartRepo) GetCartByUserID(ctx context.Context, userID models.UserID) ([]models.CartItem, error) {
	query := `SELECT sku_id, count FROM cart WHERE user_id = $1`

	rows, err := c.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var items []models.CartItem

	for rows.Next() {
		var dbItem cartItemDB
		if err := rows.Scan(&dbItem.SKUID, &dbItem.Count); err != nil {
			return nil, err
		}

		skuID, err := models.Int64ToSKUID(dbItem.SKUID)
		if err != nil {
			return nil, err
		}

		items = append(items, models.CartItem{
			SKUID: skuID,
			Count: dbItem.Count,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (c *CartRepo) ClearCartByUserID(ctx context.Context, userID models.UserID) error {
	query := `DELETE FROM cart WHERE user_id = $1`

	tag, err := c.db.Exec(ctx, query, userID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() < 1 {
		return ErrNotFound
	}

	return nil
}
