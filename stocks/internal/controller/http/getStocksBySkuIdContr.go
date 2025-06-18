package http

import (
	"encoding/json"
	"net/http"
	"stocks/internal/dto"
	"stocks/pkg/logger"
	"stocks/pkg/utils"
)

func (c *StockController) GetSkuStocksBySkuIdControlller(w http.ResponseWriter, r *http.Request) {
	var skuIdDto dto.GetSkuBySkuIdDto
	if err := json.NewDecoder(r.Body).Decode(&skuIdDto); err != nil {
		logger.Log.Errorf("Failed to decode request body: %v", err)
		utils.Error(w, err, http.StatusBadRequest)

		return
	}

	stock, err := c.usecase.GetSkuStocksBySkuIdUsecase(r.Context(), skuIdDto.SkuId)
	if err.Message != nil {
		logger.Log.Errorf("Failed to get stock: %v", err)
		utils.Error(w, err.Message, err.Code)

		return
	}

	logger.Log.Debugf("Stock retrieved successfully: %v", stock)
	utils.SuccessResponse(w, stock, http.StatusOK)
}
