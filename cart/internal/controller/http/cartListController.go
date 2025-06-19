package http

import (
	"cart/internal/models"
	"cart/pkg/utils"
	"encoding/json"
	"net/http"
)

type UserIdRequest struct {
	UserId int64 `json:"user_id"`
}

func (c *CartController) CartListController(w http.ResponseWriter, r *http.Request) {
	var userIdReq UserIdRequest
	err := json.NewDecoder(r.Body).Decode(&userIdReq)
	if err != nil {
		utils.Error(w, err, http.StatusBadRequest)
		return
	}

	userIdDto := models.UserID(userIdReq.UserId)

	list, err := c.usecase.CartListByUserIdUsecase(r.Context(), userIdDto)
	if err != nil {
		utils.Error(w, err, http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, list, http.StatusOK)
}
