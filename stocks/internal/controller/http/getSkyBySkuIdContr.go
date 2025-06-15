package http

import (
	"encoding/json"
	"net/http"
	"stocks/internal/dto"
	"stocks/pkg/utils"
)

func (c *StockController) GetSkuBySkuIdControlller(w http.ResponseWriter, r *http.Request) {
	var getSkuBySkuIdDto dto.GetSkuBySkuIdDto
	if err := json.NewDecoder(r.Body).Decode(&getSkuBySkuIdDto); err != nil {
		utils.Error(w, err, http.StatusBadRequest)
		return
	}

	newId := getSkuBySkuIdDto.SkuId

	item, err := c.usecase.GetSkuBySkuIdUsecase(r.Context(), newId)
	if err != nil {
		utils.Error(w, err, http.StatusNotFound)
		return
	}

	utils.SuccesResponse(w, item, http.StatusOK)
}
