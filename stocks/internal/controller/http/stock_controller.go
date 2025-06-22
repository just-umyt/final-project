package http

import (
	"encoding/json"
	"net/http"
	"stocks/internal/models"
	"stocks/internal/usecase"
	"stocks/pkg/logger"
	"stocks/pkg/utils"
)

type StockController struct {
	usecase usecase.StockUsecaseInterface
}

func NewStockController(stUsecase usecase.StockUsecaseInterface) *StockController {
	return &StockController{usecase: stUsecase}
}

const ErrBadRequest string = "Bad Request: Failed to decode request body"

func (c *StockController) AddStock(w http.ResponseWriter, r *http.Request) {
	var addItemReq AddStockRequest
	if err := json.NewDecoder(r.Body).Decode(&addItemReq); err != nil {
		logger.Log.Errorf("ADD | %s: %v", ErrBadRequest, err)
		utils.Error(w, err, http.StatusBadRequest)

		return
	}

	addItemDto := usecase.AddStockDto{
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

func (c *StockController) DeleteStockBySkuId(w http.ResponseWriter, r *http.Request) {
	var deleteStockReq DeleteStockRequest

	if err := json.NewDecoder(r.Body).Decode(&deleteStockReq); err != nil {
		logger.Log.Errorf("DELETE | %s: %v", ErrBadRequest, err)
		utils.Error(w, err, http.StatusBadRequest)

		return
	}

	deleteStockDto := usecase.DeleteStockDto{
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

func (c *StockController) GetSkusByLocation(w http.ResponseWriter, r *http.Request) {
	var paginationReq GetSkuByLocationParamsRequest
	if err := json.NewDecoder(r.Body).Decode(&paginationReq); err != nil {
		logger.Log.Errorf("GET BY LOCATION | %s: %v", ErrBadRequest, err)
		utils.Error(w, err, http.StatusBadRequest)

		return
	}

	paginationDto := usecase.GetSkuByLocationParamsDto{
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

func (c *StockController) GetSkuStocksBySkuId(w http.ResponseWriter, r *http.Request) {
	var skuIdReq GetSkuBySkuIdRequest
	if err := json.NewDecoder(r.Body).Decode(&skuIdReq); err != nil {
		logger.Log.Errorf("GET | %s: %v", ErrBadRequest, err)
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

	stockRes := StockResponse{
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
