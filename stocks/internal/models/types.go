package models

import (
	"fmt"
	"math"
)

// SKUID - type id of stock keeping unit.
type SKUID uint32

// UserID - type id of user.
type UserID int64

// StockID - type id of stock.
type StockID int64

func Uint32ToUint16(v uint32) (uint16, error) {
	if v > math.MaxUint16 {
		return 0, fmt.Errorf("%d out of uint16 range", v)
	}

	return uint16(v), nil
}

func IntToInt32(v int) (int32, error) {
	if v < math.MinInt32 || v > math.MaxInt32 {
		return 0, fmt.Errorf("%d out of int32 range", v)
	}

	return int32(v), nil
}
