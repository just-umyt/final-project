package http

import (
	"encoding/json"
	"net/http"
	"stocks/internal/dto"
	"stocks/pkg/utils"
)

func (c *StockController) GetSkusByLocationController(w http.ResponseWriter, r *http.Request) {
	var paginationReq dto.GetSkuByLocationParamsDto

	if err := json.NewDecoder(r.Body).Decode(&paginationReq); err != nil {
		utils.Error(w, err, http.StatusBadRequest)
		return
	}

	items, err := c.usecase.GetSkuByLocationUsecase(r.Context(), paginationReq)
	if err != nil && items == nil {
		utils.Error(w, err, http.StatusNotFound)
		return
	}

	paginationResponse := dto.PaginationResonse{Items: items, Err: err}

	utils.SuccesResponse(w, paginationResponse, http.StatusOK)
}
