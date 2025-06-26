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

func (c *CartController) CartAddItem(w http.ResponseWriter, r *http.Request) {
	var req CartAddItemRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Errorf("ADD | %s: %v", ErrBadRequest, err)

		utils.ErrorResponse(w, err, http.StatusBadRequest)

		return
	}

	cartAddDto := usecase.CartAddItemDto{
		UserId: models.UserID(req.UserId),
		SkuId:  models.SKUID(req.SkuId),
		Count:  req.Count,
	}

	err := c.usecase.CartAddItem(r.Context(), cartAddDto)
	if err != nil {
		if errors.Is(err, usecase.ErrNotEnoughStock) {
			logger.Log.Errorf("ADD | Item %v not found: %v", cartAddDto, err)
			utils.ErrorResponse(w, err, http.StatusPreconditionFailed)

			return
		} else {
			logger.Log.Errorf("ADD | Failed to add item to cart: %v", err)
			utils.ErrorResponse(w, err, http.StatusInternalServerError)

			return
		}
	}

	utils.SuccessResponse(w, "", http.StatusOK)
}

func (c *CartController) CartClear(w http.ResponseWriter, r *http.Request) {
	var userIdReq UserIdRequest

	err := json.NewDecoder(r.Body).Decode(&userIdReq)
	if err != nil {
		logger.Log.Errorf("CLEAR | %s: %v", ErrBadRequest, err)

		utils.ErrorResponse(w, err, http.StatusBadRequest)

		return
	}

	userIdDto := models.UserID(userIdReq.UserId)

	err = c.usecase.CartClearByUserId(r.Context(), userIdDto)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			logger.Log.Errorf("CLEAR | User %v not found: %v", userIdDto, err)
			utils.ErrorResponse(w, err, http.StatusNotFound)

			return
		} else {
			logger.Log.Errorf("CLEAR | Failed to clear cart: %v", err)
			utils.ErrorResponse(w, err, http.StatusInternalServerError)

			return
		}
	}

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

	deleteItemDto := usecase.DeleteItemDto{
		UserId: models.UserID(req.UserId),
		SkuId:  models.SKUID(req.SkuId),
	}

	err = c.usecase.CartDeleteItem(r.Context(), deleteItemDto)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			logger.Log.Errorf("DELETE | Item %v not found: %v", deleteItemDto, err)
			utils.ErrorResponse(w, err, http.StatusNotFound)

			return
		} else {
			logger.Log.Errorf("DELETE | Failed to delete item from cart: %v", err)
			utils.ErrorResponse(w, err, http.StatusInternalServerError)

			return
		}
	}

	utils.SuccessResponse(w, nil, http.StatusOK)
}

func (c *CartController) CartList(w http.ResponseWriter, r *http.Request) {
	var userIdReq UserIdRequest

	err := json.NewDecoder(r.Body).Decode(&userIdReq)
	if err != nil {
		logger.Log.Errorf("LIST | %s: %v", ErrBadRequest, err)
		utils.ErrorResponse(w, err, http.StatusBadRequest)

		return
	}

	userIdDto := models.UserID(userIdReq.UserId)

	list, err := c.usecase.CartListByUserId(r.Context(), userIdDto)
	if err != nil {
		logger.Log.Errorf("LIST | Failed to get list from cart: %v", err)
		utils.ErrorResponse(w, err, http.StatusInternalServerError)

		return
	}

	utils.SuccessResponse(w, list, http.StatusOK)
}
