package repository

type GetSKU struct {
	Name  string  `db:"name"`
	Price float64 `db:"price"`
	Count int     `db:"count"`
}
