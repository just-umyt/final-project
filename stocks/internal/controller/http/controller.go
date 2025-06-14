package http

import (
	"stocks/internal/usecase"
)

type StockController struct {
	usecase usecase.StockUsecaseInterface
}

func NewStockController(stUsecase usecase.StockUsecaseInterface) *StockController {
	return &StockController{usecase: stUsecase}
}
