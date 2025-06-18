package http

import (
	"encoding/json"
	"net/http"
	"stocks/internal/models"
	"stocks/internal/usecase"
	"stocks/pkg/logger"
	"stocks/pkg/utils"
)

type GetSkuBySkuIdRequest struct {
	SkuId models.SKUID `json:"sku"`
}

type Stocks struct {
	SkuId    uint32 `json:"sku"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Count    uint16 `json:"count,omitempty"`
	Price    uint32 `json:"price,omitempty"`
	Location string `json:"location,omitempty"`
	UserId   int64  `json:"user_id,omitempty"`
}

func (c *StockController) GetSkuStocksBySkuIdControlller(w http.ResponseWriter, r *http.Request) {
	var skuIdReq GetSkuBySkuIdRequest
	if err := json.NewDecoder(r.Body).Decode(&skuIdReq); err != nil {
		logger.Log.Errorf("Failed to decode request body: %v", err)
		utils.Error(w, err, http.StatusBadRequest)

		return
	}

	skuId := models.SKUID(skuIdReq.SkuId)

	stock, err := c.usecase.GetSkuStocksBySkuIdUsecase(r.Context(), skuId)
	if err != nil {
		if err.Error() == usecase.NotFoundError {
			logger.Log.Errorf("GET | Sku %v not found: %v", skuId, err)
			utils.Error(w, err, http.StatusNotFound)

			return
		} else {
			logger.Log.Errorf("GET | Failed to get stock: %v", err)
			utils.Error(w, err, http.StatusInternalServerError)

			return
		}
	}

	stockRes := Stocks{
		SkuId:    uint32(stock.SkuId),
		Name:     stock.Name,
		Type:     stock.Type,
		Count:    stock.Count,
		Price:    stock.Price,
		Location: stock.Location,
		UserId:   int64(stock.UserId),
	}

	logger.Log.Debugf("GET | Stock retrieved successfully: %v", stock)
	utils.SuccessResponse(w, stockRes, http.StatusOK)
}
