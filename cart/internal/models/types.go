package models

import (
	"fmt"
	"math"
)

// SKUID - type id of stock keeping unit.
type SKUID uint32

// CartID - type id of cart.
type CartID uint32

// UserID - type id of user.
type UserID int64

func Int64ToSKUID(v int64) (SKUID, error) {
	if v < 0 || v > math.MaxUint32 {
		return 0, fmt.Errorf("sku_id %d out of uint32 range", v)
	}

	return SKUID(v), nil
}
