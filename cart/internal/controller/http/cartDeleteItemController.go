package http

import (
	"cart/internal/dto"
	"cart/internal/models"
	"cart/pkg/utils"
	"encoding/json"
	"net/http"
)

type DeleteItemRequest struct {
	UserId int64  `json:"user_id"`
	SkuId  uint32 `json:"sku"`
}

func (c *CartController) DeleteItemController(w http.ResponseWriter, r *http.Request) {
	var req DeleteItemRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.Error(w, err, http.StatusBadRequest)
		return
	}

	deleteItemDto := dto.DeleteItemDto{
		UserId: models.UserID(req.UserId),
		SkuId:  models.SKUID(req.SkuId),
	}

	err = c.usecase.CartDeleteItemUsecase(r.Context(), deleteItemDto)
	if err != nil {
		utils.Error(w, err, http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, nil, http.StatusOK)
}
