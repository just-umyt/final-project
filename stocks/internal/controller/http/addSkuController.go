package http

import (
	"encoding/json"
	"net/http"
	"stocks/internal/dto"
	"stocks/pkg/logger"
	"stocks/pkg/utils"
)

func (c *StockController) AddStockController(w http.ResponseWriter, r *http.Request) {
	var addItemDto dto.AddStockDto
	if err := json.NewDecoder(r.Body).Decode(&addItemDto); err != nil {
		logger.Log.Errorf("Failed to decode request body: %v", err)
		utils.Error(w, err, http.StatusBadRequest)

		return
	}

	if err := c.usecase.AddStockUsecase(r.Context(), addItemDto); err.Message != nil {
		logger.Log.Errorf("Failed to add stock: %v", err)
		utils.Error(w, err.Message, err.Code)

		return
	}

	logger.Log.Debug("Stock added successfully")
	utils.SuccessResponse(w, "", http.StatusOK)
}
