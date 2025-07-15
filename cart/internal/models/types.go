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

func Int64ToUint32(v int64) (uint32, error) {
	if v < 0 || v > math.MaxUint32 {
		return 0, fmt.Errorf("%d out of uint32 range", v)
	}

	return uint32(v), nil
}
