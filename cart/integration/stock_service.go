package integration

import (
	spb "cart/pkg/api/stock"
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
)

const (
	stockSku      = 1001
	stockCount    = 10
	stockTestNums = 1
)

type StockServer struct {
	spb.UnimplementedStockServiceServer
}

func NewStockServer() *StockServer {
	return &StockServer{}
}

func (s *StockServer) GetItem(ctx context.Context, req *spb.StockGetItemRequest) (*spb.StockItemResponse, error) {
	resp := &spb.StockItemResponse{
		Sku:      stockSku,
		Name:     "test name",
		Type:     "test type",
		Count:    stockCount,
		Price:    stockTestNums,
		Location: "test loc",
		UserId:   stockTestNums,
	}

	return resp, nil
}

func startFakeStockService(addr string) (*grpc.Server, error) {
	stockListener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	grpcServer := grpc.NewServer()

	stockServer := NewStockServer()

	spb.RegisterStockServiceServer(grpcServer, stockServer)

	go func() {
		if err := grpcServer.Serve(stockListener); err != nil {
			log.Fatalf("failed to serve stock:% v", err)
		}
	}()

	return grpcServer, nil
}
