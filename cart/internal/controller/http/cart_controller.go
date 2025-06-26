package http

import (
	"cart/internal/models"
	"cart/internal/usecase"
	"cart/pkg/logger"
	"cart/pkg/utils"
	"encoding/json"
	"errors"
	"net/http"
)

type CartController struct {
	usecase usecase.CartUsecaseInterface
}

func NewCartController(cartUsecase usecase.CartUsecaseInterface) *CartController {
	return &CartController{usecase: cartUsecase}
}

const ErrBadRequest string = "Bad Request: Failed to decode request body"

func (c *CartController) AddItem(w http.ResponseWriter, r *http.Request) {
	var req AddItemRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Errorf("ADD | %s: %v", ErrBadRequest, err)

		utils.ErrorResponse(w, err, http.StatusBadRequest)

		return
	}

	dto := usecase.AddItemDTO{
		UserID: models.UserID(req.UserID),
		SKUID:  models.SKUID(req.SKUID),
		Count:  req.Count,
	}

	err := c.usecase.AddItem(r.Context(), dto)
	if err != nil {
		if errors.Is(err, usecase.ErrNotEnoughStock) {
			logger.Log.Errorf("ADD | Item %v not found: %v", dto, err)
			utils.ErrorResponse(w, err, http.StatusPreconditionFailed)

			return
		} else {
			logger.Log.Errorf("ADD | Failed to add item to cart: %v", err)
			utils.ErrorResponse(w, err, http.StatusInternalServerError)

			return
		}
	}

	logger.Log.Debug("ADD | succes")
	utils.SuccessResponse(w, "", http.StatusOK)
}

func (c *CartController) CartClear(w http.ResponseWriter, r *http.Request) {
	var req UserIDRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Log.Errorf("CLEAR | %s: %v", ErrBadRequest, err)

		utils.ErrorResponse(w, err, http.StatusBadRequest)

		return
	}

	userID := models.UserID(req.UserID)

	err = c.usecase.ClearCartByUserID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			logger.Log.Errorf("CLEAR | User %v not found: %v", userID, err)
			utils.ErrorResponse(w, err, http.StatusNotFound)

			return
		} else {
			logger.Log.Errorf("CLEAR | Failed to clear cart: %v", err)
			utils.ErrorResponse(w, err, http.StatusInternalServerError)

			return
		}
	}

	logger.Log.Debug("ClEAR | succes")
	utils.SuccessResponse(w, "", http.StatusOK)
}

func (c *CartController) DeleteItem(w http.ResponseWriter, r *http.Request) {
	var req DeleteItemRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Log.Errorf("DELETE | %s: %v", ErrBadRequest, err)
		utils.ErrorResponse(w, err, http.StatusBadRequest)

		return
	}

	dto := usecase.DeleteItemDTO{
		UserID: models.UserID(req.UserID),
		SKUID:  models.SKUID(req.SKUID),
	}

	err = c.usecase.DeleteItem(r.Context(), dto)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			logger.Log.Errorf("DELETE | Item %v not found: %v", dto, err)
			utils.ErrorResponse(w, err, http.StatusNotFound)

			return
		} else {
			logger.Log.Errorf("DELETE | Failed to delete item from cart: %v", err)
			utils.ErrorResponse(w, err, http.StatusInternalServerError)

			return
		}
	}

	logger.Log.Debug("DELETE | succes")
	utils.SuccessResponse(w, nil, http.StatusOK)
}

func (c *CartController) CartList(w http.ResponseWriter, r *http.Request) {
	var req UserIDRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Log.Errorf("LIST | %s: %v", ErrBadRequest, err)
		utils.ErrorResponse(w, err, http.StatusBadRequest)

		return
	}

	userID := models.UserID(req.UserID)

	resp, err := c.usecase.GetItemsByUserID(r.Context(), userID)
	if err != nil {
		logger.Log.Errorf("LIST | Failed to get list from cart: %v", err)
		utils.ErrorResponse(w, err, http.StatusInternalServerError)

		return
	}

	logger.Log.Debug("LIST | succes")
	utils.SuccessResponse(w, resp, http.StatusOK)
}
