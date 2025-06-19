package http

import (
	"cart/internal/models"
	"cart/internal/usecase"
	"cart/pkg/logger"
	"cart/pkg/utils"
	"encoding/json"
	"net/http"
)

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
