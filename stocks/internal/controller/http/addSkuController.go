package http

import (
	"encoding/json"
	"net/http"
	"stocks/internal/dto"
	"stocks/internal/models"
	"stocks/internal/usecase"
	"stocks/pkg/logger"
	"stocks/pkg/utils"
)

type AddStockRequest struct {
	SkuId    uint32 `json:"sku"`
	UserId   int64  `json:"user_id"`
	Count    uint16 `json:"count"`
	Price    uint32 `json:"price"`
	Location string `json:"location"`
}

func (c *StockController) AddStockController(w http.ResponseWriter, r *http.Request) {
	var addItemReq AddStockRequest
	if err := json.NewDecoder(r.Body).Decode(&addItemReq); err != nil {
		logger.Log.Errorf("ADD | %s: %v", ErrBadRequest, err)
		utils.Error(w, err, http.StatusBadRequest)

		return
	}

	addItemDto := dto.AddStockDto{
		SkuId:    models.SKUID(addItemReq.SkuId),
		UserId:   models.UserID(addItemReq.UserId),
		Count:    addItemReq.Count,
		Price:    addItemReq.Price,
		Location: addItemReq.Location,
	}

	if err := c.usecase.AddStockUsecase(r.Context(), addItemDto); err != nil {
		switch err.Error() {
		case usecase.NotFoundError:
			logger.Log.Errorf("ADD | Sku %v not found: %v", addItemDto.SkuId, err)
			utils.Error(w, err, http.StatusNotFound)

			return
		case usecase.UserIdError:
			logger.Log.Errorf("ADD | User %v not found: %v", addItemDto.UserId, err)
			utils.Error(w, err, http.StatusNotFound)

			return
		default:
			logger.Log.Errorf("ADD | Failed to add stock: %v", err)
			utils.Error(w, err, http.StatusInternalServerError)

			return
		}
	}

	logger.Log.Debug("ADD STOCK SUCCES")
	utils.SuccessResponse(w, "", http.StatusOK)
}
