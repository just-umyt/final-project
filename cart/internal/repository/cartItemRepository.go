package repository

import (
	"cart/internal/models"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type CartRepoInterface interface {
	GetItemIdByUserId(ctx context.Context, userId models.UserID, skuId models.SKUID) (models.CartID, error)
	UpdateItemByUserId(ctx context.Context, cart models.Cart) error
	AddItem(ctx context.Context, cart models.Cart) error
	DeleteItem(ctx context.Context, userId models.UserID, sku models.SKUID) (int64, error)
	GetItemsSkuIdByCartId(ctx context.Context, userId models.UserID) ([]models.SKUID, error)
	ClearCartByUserId(ctx context.Context, userId models.UserID) (int64, error)
}

type CartRepo struct {
	tx pgx.Tx
}

func NewCartRepository(tx pgx.Tx) *CartRepo {
	return &CartRepo{tx: tx}
}

func (c *CartRepo) GetItemIdByUserId(ctx context.Context, userId models.UserID, skuId models.SKUID) (models.CartID, error) {
	query := `SELECT id FROM cart WHERE user_id = $1 AND sku_id = $2 `

	var cartId models.CartID

	err := c.tx.QueryRow(ctx, query, userId, skuId).Scan(&cartId)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return 0, err
	}

	return cartId, nil
}

func (c *CartRepo) UpdateItemByUserId(ctx context.Context, cart models.Cart) error {
	query := `UPDATE cart SET count = $1 WHERE user_id = $2 AND sku_id = $3`
	_, err := c.tx.Exec(ctx, query, cart.Count, cart.UserId, cart.SKUId)

	return err
}

func (c *CartRepo) AddItem(ctx context.Context, cart models.Cart) error {
	query := `INSERT INTO cart (user_id, sku_id, count) VALUES ($1, $2, $3)`

	_, err := c.tx.Exec(ctx, query, cart.UserId, cart.SKUId, cart.Count)

	return err
}

func (c *CartRepo) DeleteItem(ctx context.Context, userId models.UserID, skuId models.SKUID) (int64, error) {
	query := `DELETE FROM cart WHERE user_id = $1 AND sku_id = $2`
	tag, err := c.tx.Exec(ctx, query, userId, skuId)

	return tag.RowsAffected(), err
}

func (c *CartRepo) GetItemsSkuIdByCartId(ctx context.Context, userId models.UserID) ([]models.SKUID, error) {
	query := `SELECT sku_id FROM cart WHERE user_id = $1`

	rows, err := c.tx.Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var skus []models.SKUID

	for rows.Next() {
		var skuId models.SKUID
		if err := rows.Scan(&skuId); err != nil {
			return nil, err
		}

		skus = append(skus, skuId)
	}

	return skus, nil
}

func (c *CartRepo) ClearCartByUserId(ctx context.Context, userId models.UserID) (int64, error) {
	query := `DELETE FROM cart WHERE user_id = $1`

	tag, err := c.tx.Exec(ctx, query, userId)
	if err != nil {
		return 0, err
	}

	return tag.RowsAffected(), nil
}
