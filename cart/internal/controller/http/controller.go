package http

import "cart/internal/usecase"

type CartController struct {
	usecase usecase.CartUsecaseInterface
}

func NewCartController(cartUsecase usecase.CartUsecaseInterface) *CartController {
	return &CartController{usecase: cartUsecase}
}
