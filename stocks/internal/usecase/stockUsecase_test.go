package usecase

import (
	"log"
	"stocks/internal/trManager/mock"
	"testing"
)

func TestAddStock(t *testing.T) {
	//when i declare these error appears
	repo := mock.NewIStockRepoMock(t)
	trx := mock.NewIPgTxManagerMock(t)
	log.Println(repo)
	log.Println(trx)
}
