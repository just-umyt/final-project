package http

import (
	"encoding/json"
	"net/http"
	"stocks/internal/dto"
	"stocks/internal/models"
	"stocks/pkg/utils"
)

func (c *StockController) AddSkuController(w http.ResponseWriter, r *http.Request) {
	var addItemDto dto.AddSkuDto
	if err := json.NewDecoder(r.Body).Decode(&addItemDto); err != nil {
		utils.Error(w, err, http.StatusBadRequest)
	}

	newSku := models.SKU{
		SkuId:    addItemDto.SkuId,
		Name:     addItemDto.Name,
		Count:    addItemDto.Count,
		Type:     addItemDto.Type,
		Price:    addItemDto.Price,
		Location: addItemDto.Location,
		UserId:   addItemDto.UserId,
	}

	if err := c.usecase.AddSkuUsecase(r.Context(), newSku); err != nil {
		utils.Error(w, err, http.StatusInternalServerError)

		return
	}

	utils.SuccesResponse(w, "", http.StatusOK)
}
