package services

import (
	"cart/internal/models"
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "cart/pkg/api/stock"
)

const (
	ctxTimeout             = 5
	errorConvertStockCount = "failed to convert stock count: %w"
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

	grpcCtx, cancel := context.WithTimeout(ctx, ctxTimeout*time.Second)
	defer cancel()

	resp, err := client.GetItem(grpcCtx, &req)
	if err != nil {
		log.Println(err)
		return ItemDTO{}, err
	}

	count, err := models.Uint32ToUint16(resp.Count)
	if err != nil {
		return ItemDTO{}, fmt.Errorf(errorConvertStockCount, err)
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
