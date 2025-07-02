package usecase

import (
	"stocks/internal/usecase/mock"
	"testing"
)

func TestAddStock(t *testing.T) {
	//when i declare these error appears
	repo := mock.NewIStockRepoMock(t)
	trx := mock.NewIPgTxManagerMock(t)
}
