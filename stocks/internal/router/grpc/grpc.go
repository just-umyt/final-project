package grpc

import (
	"context"
	"errors"
	"fmt"
	"stocks/internal/models"
	"stocks/internal/usecase"
	pb "stocks/pkg/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type IStockUsecase interface {
	AddStock(ctx context.Context, stock usecase.AddStockDTO) error
	DeleteStockBySKU(ctx context.Context, delStock usecase.DeleteStockDTO) error
	GetStocksByLocation(ctx context.Context, param usecase.GetItemByLocDTO) (usecase.ItemsByLocDTO, error)
	GetItemBySKU(ctx context.Context, sku models.SKUID) (usecase.StockDTO, error)
}

type StockServer struct {
	stockUsecase IStockUsecase
	pb.UnimplementedStockServiceServer
}

func NewStockServer(us IStockUsecase) *StockServer {
	return &StockServer{stockUsecase: us}
}

func (s *StockServer) AddItem(ctx context.Context, req *pb.StockAddIemRequest) (*emptypb.Empty, error) {
	count, err := models.Uint32ToUint16(req.Count)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	dto := usecase.AddStockDTO{
		SKUID:    models.SKUID(req.Sku),
		UserID:   models.UserID(req.UserId),
		Count:    count,
		Price:    req.Price,
		Location: req.Location,
	}

	if err := s.stockUsecase.AddStock(ctx, dto); err != nil {
		if errors.Is(err, usecase.ErrNotFound) || errors.Is(err, usecase.ErrUserID) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *StockServer) DeleteItem(ctx context.Context, req *pb.StockDeleteItemRequest) (*emptypb.Empty, error) {
	dto := usecase.DeleteStockDTO{
		UserID: models.UserID(req.UserId),
		SKUID:  models.SKUID(req.Sku),
	}

	if err := s.stockUsecase.DeleteStockBySKU(ctx, dto); err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *StockServer) ListItem(ctx context.Context, req *pb.StockListItemRequest) (*pb.StockListItemResponse, error) {
	dto := usecase.GetItemByLocDTO{
		UserID:      models.UserID(req.UserId),
		Location:    req.Location,
		PageSize:    req.PageSize,
		CurrentPage: req.CurrentPage,
	}

	list, err := s.stockUsecase.GetStocksByLocation(ctx, dto)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	respList := make([]*pb.StockItemResponse, len(list.Stocks))

	for i, item := range list.Stocks {
		respItem := pb.StockItemResponse{}

		respItem.Sku = uint32(item.SKU.SKUID)
		respItem.Name = item.SKU.Name
		respItem.Type = item.SKU.Type
		respItem.Count = uint32(item.Count)
		respItem.Price = item.Price
		respItem.Location = item.Location
		respItem.UserId = int64(item.UserID)

		respList[i] = &respItem
	}

	totalCount, err := models.IntToInt32(list.TotalCount)
	if err != nil {
		return nil, fmt.Errorf("totalCount = %w", err)
	}

	response := &pb.StockListItemResponse{
		Items:      respList,
		TotalCount: totalCount,
		PageNumber: list.PageNumber,
	}

	return response, nil
}

func (s *StockServer) GetItem(ctx context.Context, req *pb.StockGetItemRequest) (*pb.StockItemResponse, error) {
	item, err := s.stockUsecase.GetItemBySKU(ctx, models.SKUID(req.Sku))
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &pb.StockItemResponse{
		Sku:      uint32(item.SKU.SKUID),
		Name:     item.SKU.Name,
		Type:     item.SKU.Type,
		Count:    uint32(item.Count),
		Price:    item.Price,
		Location: item.Location,
		UserId:   int64(item.UserID),
	}, nil
}
