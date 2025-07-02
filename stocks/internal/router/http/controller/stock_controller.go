package controller

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"stocks/internal/models"
	"stocks/internal/usecase"
	"stocks/pkg/utils"
	"stocks/pkg/validation"
)

//go:generate mkdir -p mock
//go:generate minimock -o ./mock/ -s .go  -g
type IStockUsecase interface {
	AddStock(ctx context.Context, stock usecase.AddStockDTO) error
	DeleteStockBySKU(ctx context.Context, delStock usecase.DeleteStockDTO) error
	GetStocksByLocation(ctx context.Context, param usecase.GetItemByLocDTO) (usecase.ItemsByLocDTO, error)
	GetItemBySKU(ctx context.Context, sku models.SKUID) (usecase.StockDTO, error)
}

type StockController struct {
	stockUsecase IStockUsecase
}

func NewStockController(stockUsecase IStockUsecase) *StockController {
	return &StockController{stockUsecase: stockUsecase}
}

func (c *StockController) AddStock(w http.ResponseWriter, r *http.Request) {
	var req AddStockRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)

		return
	}

	if err := validation.IsValid(req); err != nil {
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

	if err := c.stockUsecase.AddStock(r.Context(), dto); err != nil {
		if errors.Is(err, usecase.ErrNotFound) || errors.Is(err, usecase.ErrUserID) {
			utils.ErrorResponse(w, err, http.StatusNotFound)

			return
		}

		utils.ErrorResponse(w, err, http.StatusInternalServerError)

		return
	}

	log.Println("ADD | succes")
	utils.SuccessResponse(w, "", http.StatusOK)
}

func (c *StockController) DeleteStockBySKU(w http.ResponseWriter, r *http.Request) {
	var req DeleteStockRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)

		return
	}

	if err := validation.IsValid(req); err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)

		return
	}

	dto := usecase.DeleteStockDTO{
		UserID: models.UserID(req.UserID),
		SKUID:  models.SKUID(req.SKUID),
	}

	if err := c.stockUsecase.DeleteStockBySKU(r.Context(), dto); err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			utils.ErrorResponse(w, err, http.StatusNotFound)

			return
		}

		utils.ErrorResponse(w, err, http.StatusInternalServerError)

		return
	}

	log.Println("DELETE | succes", dto)
	utils.SuccessResponse(w, "", http.StatusOK)
}

func (c *StockController) GetItemsByLocation(w http.ResponseWriter, r *http.Request) {
	var req GetItemsByLocRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)

		return
	}

	if err := validation.IsValid(req); err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)

		return
	}

	dto := usecase.GetItemByLocDTO{
		UserID:      models.UserID(req.UserID),
		Location:    req.Location,
		PageSize:    req.PageSize,
		CurrentPage: req.CurrentPage,
	}

	items, err := c.stockUsecase.GetStocksByLocation(r.Context(), dto)
	if err != nil {
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

	log.Println("GET ITEMS | succes")
	utils.SuccessResponse(w, resp, http.StatusOK)
}

func (c *StockController) GetItemBySKU(w http.ResponseWriter, r *http.Request) {
	var req GetItemBySKURequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)

		return
	}

	if err := validation.IsValid(req); err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)

		return
	}

	skuID := models.SKUID(req.SKU)

	item, err := c.stockUsecase.GetItemBySKU(r.Context(), skuID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			utils.ErrorResponse(w, err, http.StatusNotFound)

			return
		}

		utils.ErrorResponse(w, err, http.StatusInternalServerError)

		return
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

	log.Printf("GET | succes: %v", item)
	utils.SuccessResponse(w, resp, http.StatusOK)
}
