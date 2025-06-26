package repository

import (
	"cart/internal/models"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type CartRepoInterface interface {
	GetCartIDByUserID(ctx context.Context, userID models.UserID, skuID models.SKUID) (models.CartID, error)
	UpdateItemByUserID(ctx context.Context, cart models.Cart) error
	AddItem(ctx context.Context, cart models.Cart) error
	DeleteItem(ctx context.Context, userID models.UserID, skuID models.SKUID) (int64, error)
	GetSKUIDsByUserID(ctx context.Context, userID models.UserID) ([]models.SKUID, error)
	ClearCartByUserID(ctx context.Context, userID models.UserID) (int64, error)
}

type CartRepo struct {
	tx pgx.Tx
}

func NewCartRepository(tx pgx.Tx) *CartRepo {
	return &CartRepo{tx: tx}
}

func (c *CartRepo) GetCartIDByUserID(ctx context.Context, userID models.UserID, skuID models.SKUID) (models.CartID, error) {
	query := `SELECT id FROM cart WHERE user_id = $1 AND sku_id = $2`

	var cartID models.CartID

	err := c.tx.QueryRow(ctx, query, userID, skuID).Scan(&cartID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return 0, err
	}

	return cartID, nil
}

func (c *CartRepo) UpdateItemByUserID(ctx context.Context, cart models.Cart) error {
	query := `UPDATE cart SET count = $1 WHERE user_id = $2 AND sku_id = $3`
	_, err := c.tx.Exec(ctx, query, cart.Count, cart.UserID, cart.SKUID)

	return err
}

func (c *CartRepo) AddItem(ctx context.Context, cart models.Cart) error {
	query := `INSERT INTO cart (user_id, sku_id, count) VALUES ($1, $2, $3)`

	_, err := c.tx.Exec(ctx, query, cart.UserID, cart.SKUID, cart.Count)

	return err
}

func (c *CartRepo) DeleteItem(ctx context.Context, userID models.UserID, skuID models.SKUID) (int64, error) {
	query := `DELETE FROM cart WHERE user_id = $1 AND sku_id = $2`
	tag, err := c.tx.Exec(ctx, query, userID, skuID)

	return tag.RowsAffected(), err
}

func (c *CartRepo) GetSKUIDsByUserID(ctx context.Context, userID models.UserID) ([]models.SKUID, error) {
	query := `SELECT sku_id FROM cart WHERE user_id = $1`

	rows, err := c.tx.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var skuIDs []models.SKUID

	for rows.Next() {
		var skuID models.SKUID
		if err := rows.Scan(&skuID); err != nil {
			return nil, err
		}

		skuIDs = append(skuIDs, skuID)
	}

	return skuIDs, nil
}

func (c *CartRepo) ClearCartByUserID(ctx context.Context, userID models.UserID) (int64, error) {
	query := `DELETE FROM cart WHERE user_id = $1`

	tag, err := c.tx.Exec(ctx, query, userID)
	if err != nil {
		return 0, err
	}

	return tag.RowsAffected(), nil
}
