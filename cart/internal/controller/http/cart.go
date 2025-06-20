package http

import (
	"cart/internal/dto"
	"cart/internal/models"
	"cart/internal/usecase"
	"cart/pkg/logger"
	"cart/pkg/utils"
	"encoding/json"
	"net/http"
)

func (c *CartController) CartAddItemController(w http.ResponseWriter, r *http.Request) {
	var req CartAddItemRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Errorf("ADD | %s: %v", ErrBadRequest, err)

		utils.Error(w, err, http.StatusBadRequest)

		return
	}

	cartAddDto := dto.CartAddItemDto{
		UserId: models.UserID(req.UserId),
		SkuId:  models.SKUID(req.SkuId),
		Count:  req.Count,
	}

	err := c.usecase.CartAddItemUsecase(r.Context(), cartAddDto)
	if err != nil {
		if err.Error() == usecase.NotEnoughStock {
			logger.Log.Errorf("ADD | Item %v not found: %v", cartAddDto, err)
			utils.Error(w, err, http.StatusPreconditionFailed)

			return
		} else {
			logger.Log.Errorf("ADD | Failed to add item to cart: %v", err)
			utils.Error(w, err, http.StatusInternalServerError)

			return
		}
	}

	utils.SuccessResponse(w, "", http.StatusOK)
}

func (c *CartController) CartClearController(w http.ResponseWriter, r *http.Request) {
	var userIdReq UserIdRequest

	err := json.NewDecoder(r.Body).Decode(&userIdReq)
	if err != nil {
		logger.Log.Errorf("CLEAR | %s: %v", ErrBadRequest, err)

		utils.Error(w, err, http.StatusBadRequest)

		return
	}

	userIdDto := models.UserID(userIdReq.UserId)

	err = c.usecase.CartClearByUserIdUsecase(r.Context(), userIdDto)
	if err != nil {
		if err.Error() == usecase.NotFoundError {
			logger.Log.Errorf("CLEAR | User %v not found: %v", userIdDto, err)
			utils.Error(w, err, http.StatusNotFound)

			return
		} else {
			logger.Log.Errorf("CLEAR | Failed to clear cart: %v", err)
			utils.Error(w, err, http.StatusInternalServerError)

			return
		}
	}

	utils.SuccessResponse(w, "", http.StatusOK)
}

func (c *CartController) DeleteItemController(w http.ResponseWriter, r *http.Request) {
	var req DeleteItemRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Log.Errorf("DELETE | %s: %v", ErrBadRequest, err)
		utils.Error(w, err, http.StatusBadRequest)

		return
	}

	deleteItemDto := dto.DeleteItemDto{
		UserId: models.UserID(req.UserId),
		SkuId:  models.SKUID(req.SkuId),
	}

	err = c.usecase.CartDeleteItemUsecase(r.Context(), deleteItemDto)
	if err != nil {
		if err.Error() == usecase.NotFoundError {
			logger.Log.Errorf("DELETE | Item %v not found: %v", deleteItemDto, err)
			utils.Error(w, err, http.StatusNotFound)

			return
		} else {
			logger.Log.Errorf("DELETE | Failed to delete item from cart: %v", err)
			utils.Error(w, err, http.StatusInternalServerError)

			return
		}
	}

	utils.SuccessResponse(w, nil, http.StatusOK)
}

func (c *CartController) CartListController(w http.ResponseWriter, r *http.Request) {
	var userIdReq UserIdRequest

	err := json.NewDecoder(r.Body).Decode(&userIdReq)
	if err != nil {
		logger.Log.Errorf("LIST | %s: %v", ErrBadRequest, err)
		utils.Error(w, err, http.StatusBadRequest)

		return
	}

	userIdDto := models.UserID(userIdReq.UserId)

	list, err := c.usecase.CartListByUserIdUsecase(r.Context(), userIdDto)
	if err != nil {
		logger.Log.Errorf("LIST | Failed to get list from cart: %v", err)
		utils.Error(w, err, http.StatusInternalServerError)

		return
	}

	utils.SuccessResponse(w, list, http.StatusOK)
}
