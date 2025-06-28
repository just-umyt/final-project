package http

import (
	"cart/internal/models"
	"cart/internal/usecase"
	"cart/pkg/utils"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type ICartController interface {
	AddItem(w http.ResponseWriter, r *http.Request)
	CartClear(w http.ResponseWriter, r *http.Request)
	DeleteItem(w http.ResponseWriter, r *http.Request)
	CartList(w http.ResponseWriter, r *http.Request)
}

type CartController struct {
	usecase usecase.ICartUsecase
}

func NewCartController(cartUsecase usecase.ICartUsecase) *CartController {
	return &CartController{usecase: cartUsecase}
}

const ErrBadRequest string = "Bad Request: Failed to decode request body"

func (c *CartController) AddItem(w http.ResponseWriter, r *http.Request) {
	var req AddItemRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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
			utils.ErrorResponse(w, err, http.StatusPreconditionFailed)

			return
		}

		utils.ErrorResponse(w, err, http.StatusInternalServerError)

		return
	}

	log.Println("ADD | succes")
	utils.SuccessResponse(w, "", http.StatusOK)
}

func (c *CartController) CartClear(w http.ResponseWriter, r *http.Request) {
	var req UserIDRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)

		return
	}

	userID := models.UserID(req.UserID)

	err = c.usecase.ClearCartByUserID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			utils.ErrorResponse(w, err, http.StatusNotFound)

			return
		}

		utils.ErrorResponse(w, err, http.StatusInternalServerError)

		return
	}

	log.Println("ClEAR | succes")
	utils.SuccessResponse(w, "", http.StatusOK)
}

func (c *CartController) DeleteItem(w http.ResponseWriter, r *http.Request) {
	var req DeleteItemRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
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
			utils.ErrorResponse(w, err, http.StatusNotFound)

			return
		}

		utils.ErrorResponse(w, err, http.StatusInternalServerError)

		return
	}

	log.Println("DELETE | succes")
	utils.SuccessResponse(w, nil, http.StatusOK)
}

func (c *CartController) CartList(w http.ResponseWriter, r *http.Request) {
	var req UserIDRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)

		return
	}

	userID := models.UserID(req.UserID)

	resp, err := c.usecase.GetItemsByUserID(r.Context(), userID)
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusInternalServerError)

		return
	}

	log.Println("LIST | succes")
	utils.SuccessResponse(w, resp, http.StatusOK)
}
