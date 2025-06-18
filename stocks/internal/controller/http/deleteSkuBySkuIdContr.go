package http

import (
	"encoding/json"
	"net/http"
	"stocks/internal/dto"
	"stocks/pkg/logger"
	"stocks/pkg/utils"
)

func (c *StockController) DeleteStockBySkuIdController(w http.ResponseWriter, r *http.Request) {
	var deleteSkuDto dto.DeleteStockDto

	if err := json.NewDecoder(r.Body).Decode(&deleteSkuDto); err != nil {
		logger.Log.Errorf("Failed to decode request body: %v", err)
		utils.Error(w, err, http.StatusBadRequest)

		return
	}

	if err := c.usecase.DeleteStockBySkuIdUsecase(r.Context(), deleteSkuDto); err.Message != nil {
		logger.Log.Errorf("Failed to delete stock: %v", err)
		utils.Error(w, err.Message, err.Code)

		return
	}

	logger.Log.Debugf("Stock deleted successfully: %v", deleteSkuDto)
	utils.SuccessResponse(w, "", http.StatusOK)
}
