package http

import (
	"encoding/json"
	"errors"
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
	var req AddStockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Errorf("ADD | %s: %v", ErrBadRequest, err)
		utils.ErrorResponse(w, err, http.StatusBadRequest)

		return
	}

	dto := usecase.AddStockDTO{
		SKUID:    models.SKUID(req.SKUID),
		UserID:   models.UserID(req.UserID),
		Count:    req.Count,
		Price:    req.Price,
		Location: req.Location,
	}

	if err := c.usecase.AddStock(r.Context(), dto); err != nil {
		switch {
		case errors.Is(err, usecase.ErrNotFound):
			logger.Log.Errorf("ADD | SKU %v not found: %v", dto.SKUID, err)
			utils.ErrorResponse(w, err, http.StatusNotFound)

			return
		case errors.Is(err, usecase.ErrUserID):
			logger.Log.Errorf("ADD | User %v not found: %v", dto.UserID, err)
			utils.ErrorResponse(w, err, http.StatusNotFound)

			return
		default:
			logger.Log.Errorf("ADD | Failed to add stock: %v", err)
			utils.ErrorResponse(w, err, http.StatusInternalServerError)

			return
		}
	}

	logger.Log.Debug("ADD | succes")
	utils.SuccessResponse(w, "", http.StatusOK)
}

func (c *StockController) DeleteStockBySKU(w http.ResponseWriter, r *http.Request) {
	var req DeleteStockRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Errorf("DELETE | %s: %v", ErrBadRequest, err)
		utils.ErrorResponse(w, err, http.StatusBadRequest)

		return
	}

	dto := usecase.DeleteStockDTO{
		UserID: models.UserID(req.UserID),
		SKUID:  models.SKUID(req.SKUID),
	}

	if err := c.usecase.DeleteStockBySKU(r.Context(), dto); err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			logger.Log.Errorf("DELETE | Sku %v not found: %v", dto.SKUID, err)
			utils.ErrorResponse(w, err, http.StatusNotFound)

			return
		} else {
			logger.Log.Errorf("DELETE | Failed to delete stock: %v", err)
			utils.ErrorResponse(w, err, http.StatusInternalServerError)

			return
		}
	}

	logger.Log.Debug("DELETE | succes", dto)
	utils.SuccessResponse(w, "", http.StatusOK)
}

func (c *StockController) GetItemsByLocation(w http.ResponseWriter, r *http.Request) {
	var req GetItemsByLocRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Errorf("GET ITEMS | %s: %v", ErrBadRequest, err)
		utils.ErrorResponse(w, err, http.StatusBadRequest)

		return
	}

	dto := usecase.GetItemByLocDTO{
		UserID:      models.UserID(req.UserID),
		Location:    req.Location,
		PageSize:    req.PageSize,
		CurrentPage: req.CurrentPage,
	}

	items, err := c.usecase.GetStocksByLocation(r.Context(), dto)
	if err != nil {
		logger.Log.Errorf("GET ITEMS | Failed to get stocks by location: %v", err)
		utils.ErrorResponse(w, err, http.StatusInternalServerError)

		return
	}

	var resp StockByLocResponse

	for _, stock := range items.Stocks {
		item := ItemResponse{
			SKU:      uint32(stock.SKU.SKUID),
			Name:     stock.SKU.Name,
			Type:     stock.SKU.Type,
			Count:    stock.Count,
			Price:    stock.Price,
			Location: stock.Location,
			UserID:   int64(stock.UserID),
		}
		resp.Items = append(resp.Items, item)
	}

	resp.TotalCount = items.TotalCount
	resp.PageNumber = items.PageNumber

	logger.Log.Debug("GET ITEMS | succes")
	utils.SuccessResponse(w, resp, http.StatusOK)
}

func (c *StockController) GetItemBySKU(w http.ResponseWriter, r *http.Request) {
	var req GetItemBySKURequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Errorf("GET | %s: %v", ErrBadRequest, err)
		utils.ErrorResponse(w, err, http.StatusBadRequest)

		return
	}

	skuID := models.SKUID(req.SKU)

	item, err := c.usecase.GetItemBySKU(r.Context(), skuID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			logger.Log.Errorf("GET | SKU %v not found: %v", skuID, err)
			utils.ErrorResponse(w, err, http.StatusNotFound)

			return
		} else {
			logger.Log.Errorf("GET | Failed to get item: %v", err)
			utils.ErrorResponse(w, err, http.StatusInternalServerError)

			return
		}
	}

	resp := ItemResponse{
		SKU:      uint32(item.SKU.SKUID),
		Name:     item.SKU.Name,
		Type:     item.SKU.Type,
		Count:    item.Count,
		Price:    item.Price,
		Location: item.Location,
		UserID:   int64(item.UserID),
	}

	logger.Log.Debugf("GET | succes: %v", item)
	utils.SuccessResponse(w, resp, http.StatusOK)
}
