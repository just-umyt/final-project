package http

import (
	"encoding/json"
	"net/http"
	"stocks/internal/dto"
	"stocks/pkg/utils"
)

func (c *StockController) DeleteSkuBySkuIdController(w http.ResponseWriter, r *http.Request) {
	var deleteSkuDto dto.DeleteSkuDto

	if err := json.NewDecoder(r.Body).Decode(&deleteSkuDto); err != nil {
		utils.Error(w, err, http.StatusBadRequest)
		return
	}

	if err := c.usecase.DeleteSkuBySkuIdUsecase(r.Context(), deleteSkuDto); err != nil {
		utils.Error(w, err, http.StatusNotFound)
		return
	}

	utils.SuccesResponse(w, "", http.StatusOK)
}
