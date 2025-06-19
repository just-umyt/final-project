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

type DeleteStockRequest struct {
	UserId int64  `json:"user_id"`
	SkuId  uint32 `json:"sku"`
}

func (c *StockController) DeleteStockBySkuIdController(w http.ResponseWriter, r *http.Request) {
	var deleteStockReq DeleteStockRequest

	if err := json.NewDecoder(r.Body).Decode(&deleteStockReq); err != nil {
		logger.Log.Errorf("DELETE | %s: %v", ErrBadRequest, err)
		utils.Error(w, err, http.StatusBadRequest)

		return
	}

	deleteStockDto := dto.DeleteStockDto{
		UserId: models.UserID(deleteStockReq.UserId),
		SkuId:  models.SKUID(deleteStockReq.SkuId),
	}

	if err := c.usecase.DeleteStockBySkuIdUsecase(r.Context(), deleteStockDto); err != nil {
		if err.Error() == usecase.NotFoundError {
			logger.Log.Errorf("DELETE | Sku %v not found: %v", deleteStockDto.SkuId, err)
			utils.Error(w, err, http.StatusNotFound)

			return
		} else {
			logger.Log.Errorf("DELETE | Failed to delete stock: %v", err)
			utils.Error(w, err, http.StatusInternalServerError)

			return
		}
	}

	logger.Log.Debug("DELETE | STOCK SUCCES", deleteStockDto)
	utils.SuccessResponse(w, "", http.StatusOK)
}
