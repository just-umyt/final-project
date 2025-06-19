package http

import (
	"cart/internal/models"
	"cart/pkg/utils"
	"encoding/json"
	"net/http"
)

func (c *CartController) CartClearController(w http.ResponseWriter, r *http.Request) {
	var userIdReq UserIdRequest
	err := json.NewDecoder(r.Body).Decode(&userIdReq)
	if err != nil {
		utils.Error(w, err, http.StatusBadRequest)
		return
	}

	userIdDto := models.UserID(userIdReq.UserId)

	err = c.usecase.CartClearByUserIdUsecase(r.Context(), userIdDto)
	if err != nil {
		utils.Error(w, err, http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, "", http.StatusOK)

}
