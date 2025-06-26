package repository

import (
	"stocks/internal/models"

	"github.com/jackc/pgx/v5/pgtype"
)

type GetStockByLocation struct {
	UserID   models.UserID
	Location string
	Limit    int64
	Offset   int64
}

type Stock struct {
	ID       pgtype.Int8   `db:"id"`
	SKUID    pgtype.Uint32 `db:"sku_id"`
	Count    pgtype.Uint32 `db:"count"`
	Price    pgtype.Uint32 `db:"price"`
	Location pgtype.Text   `db:"location"`
	UserID   pgtype.Int8   `db:"user_id"`
}

type SKU struct {
	ID   models.SKUID `db:"sku_id"`
	Name string       `db:"name"`
	Type string       `db:"type"`
}
