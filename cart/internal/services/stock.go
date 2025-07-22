package services

import (
	"cart/internal/models"
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"

	pb "stocks/pkg/api"
)

type StockService struct {
	client *grpc.ClientConn
}

func NewStockClient(cl *grpc.ClientConn) *StockService {
	return &StockService{
		client: cl,
	}
}

func (s *StockService) GetItemInfo(ctx context.Context, skuID models.SKUID) (ItemDTO, error) {
	client := pb.NewStockServiceClient(s.client)
	req := pb.StockGetItemRequest{Sku: uint32(skuID)}

	resp, err := client.GetItem(ctx, &req)
	if err != nil {
		log.Println(err)
		return ItemDTO{}, err
	}

	count, err := models.Uint32ToUint16(resp.Count)
	if err != nil {
		return ItemDTO{}, fmt.Errorf("failed to convert stock count: %w", err)
	}

	return ItemDTO{
		SKUID:    models.SKUID(resp.Sku),
		Name:     resp.Name,
		Type:     resp.Type,
		Count:    count,
		Price:    resp.Price,
		Location: resp.Location,
		UserID:   models.UserID(resp.UserId),
	}, nil
}
