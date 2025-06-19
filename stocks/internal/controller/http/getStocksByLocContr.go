package http

import (
	"encoding/json"
	"net/http"
	"stocks/internal/dto"
	"stocks/internal/models"
	"stocks/pkg/logger"
	"stocks/pkg/utils"
)

type GetSkuByLocationParamsRequest struct {
	User_id     int64  `json:"user_id"`
	Location    string `json:"location"`
	PageSize    int64  `json:"page_size"`
	CurrentPage int64  `json:"current_page"`
}

type StockByLocResponse struct {
	Stocks     []StockResponse `json:"stocks"`
	TotalCount int             `json:"total_count"`
	PageNumber int64           `json:"page_number"`
}

func (c *StockController) GetSkusByLocationController(w http.ResponseWriter, r *http.Request) {
	var paginationReq GetSkuByLocationParamsRequest
	if err := json.NewDecoder(r.Body).Decode(&paginationReq); err != nil {
		logger.Log.Errorf("GET BY LOCATION | %s: %v", ErrBadRequest, err)
		utils.Error(w, err, http.StatusBadRequest)

		return
	}

	paginationDto := dto.GetSkuByLocationParamsDto{
		User_id:     models.UserID(paginationReq.User_id),
		Location:    paginationReq.Location,
		PageSize:    paginationReq.PageSize,
		CurrentPage: paginationReq.CurrentPage,
	}

	stockByLoc, err := c.usecase.GetStocksByLocationUsecase(r.Context(), paginationDto)
	if err != nil {
		logger.Log.Errorf("GET BY LOCATION | Failed to get stocks by location: %v", err)
		utils.Error(w, err, http.StatusInternalServerError)

		return
	}

	var stocksRes StockByLocResponse

	for _, stock := range stockByLoc.Stocks {
		st := StockResponse{
			SkuId:    uint32(stock.SkuId),
			Name:     stock.Name,
			Type:     stock.Type,
			Count:    stock.Count,
			Price:    stock.Price,
			Location: stock.Location,
			UserId:   int64(stock.UserId),
		}
		stocksRes.Stocks = append(stocksRes.Stocks, st)
	}

	stocksRes.TotalCount = stockByLoc.TotalCount
	stocksRes.PageNumber = stockByLoc.PageNumber

	logger.Log.Debug("GET | STOCKS BY LOCATION SUCCES")
	utils.SuccessResponse(w, stocksRes, http.StatusOK)
}
