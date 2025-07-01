package http

import (
	"cart/internal/models"
	"cart/internal/usecase"
	"cart/pkg/utils"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type ICartUsecase interface {
	AddItem(ctx context.Context, addItem usecase.AddItemDTO) error
	DeleteItem(ctx context.Context, delItem usecase.DeleteItemDTO) error
	GetItemsByUserID(ctx context.Context, userID models.UserID) (usecase.ListItemsDTO, error)
	ClearCartByUserID(ctx context.Context, userID models.UserID) error
}

type CartController struct {
	cartUsecase ICartUsecase
}

func NewCartController(us ICartUsecase) *CartController {
	return &CartController{cartUsecase: us}
}

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

	err := c.cartUsecase.AddItem(r.Context(), dto)
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

	err = c.cartUsecase.ClearCartByUserID(r.Context(), userID)
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

	err = c.cartUsecase.DeleteItem(r.Context(), dto)
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

	resp, err := c.cartUsecase.GetItemsByUserID(r.Context(), userID)
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusInternalServerError)

		return
	}

	log.Println("LIST | succes")
	utils.SuccessResponse(w, resp, http.StatusOK)
}
