package http

import (
	"cart/internal/dto"
	"cart/internal/models"
	"cart/internal/usecase"
	"cart/pkg/logger"
	"cart/pkg/utils"
	"encoding/json"
	"net/http"
)

type CartAddItemRequest struct {
	UserId int64  `json:"user_id"`
	SkuId  uint32 `json:"sku"`
	Count  uint16 `json:"count"`
}

func (c *CartController) CartAddItemController(w http.ResponseWriter, r *http.Request) {
	var req CartAddItemRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Errorf("ADD | %s: %v", ErrBadRequest, err)

		utils.Error(w, err, http.StatusBadRequest)
		return
	}

	cartAddDto := dto.CartAddItemDto{
		UserId: models.UserID(req.UserId),
		SkuId:  models.SKUID(req.SkuId),
		Count:  req.Count,
	}

	err := c.usecase.CartAddItemUsecase(r.Context(), cartAddDto)
	if err != nil {
		if err.Error() == usecase.NotEnoughStock {
			logger.Log.Errorf("ADD | Item %v not found: %v", cartAddDto, err)
			utils.Error(w, err, http.StatusPreconditionFailed)

			return
		} else {
			logger.Log.Errorf("ADD | Failed to add item to cart: %v", err)
			utils.Error(w, err, http.StatusInternalServerError)

			return
		}
	}

	utils.SuccessResponse(w, "", http.StatusOK)
}
