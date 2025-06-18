package http

import (
	"encoding/json"
	"net/http"
	"stocks/internal/dto"
	"stocks/pkg/logger"
	"stocks/pkg/utils"
)

func (c *StockController) GetSkusByLocationController(w http.ResponseWriter, r *http.Request) {
	var paginationReq dto.GetSkuByLocationParamsDto
	if err := json.NewDecoder(r.Body).Decode(&paginationReq); err != nil {
		logger.Log.Errorf("Failed to decode request body: %v", err)
		utils.Error(w, err, http.StatusBadRequest)

		return
	}

	stockByLoc, err := c.usecase.GetStocksByLocationUsecase(r.Context(), paginationReq)
	if err.Message != nil {
		logger.Log.Errorf("Failed to get stocks by location: %v", err)
		utils.Error(w, err.Message, err.Code)

		return
	}

	logger.Log.Debug("Stocks retrieved successfully")
	utils.SuccessResponse(w, stockByLoc, http.StatusOK)
}
