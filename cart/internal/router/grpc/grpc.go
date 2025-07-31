package grpc

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"cart/internal/models"
	"cart/internal/usecase"
	pb "cart/pkg/api/cart"
)

type ICartUsecase interface {
	AddItem(ctx context.Context, addItem usecase.AddItemDTO) error
	DeleteItem(ctx context.Context, delItem usecase.DeleteItemDTO) error
	GetItemsByUserID(ctx context.Context, userID models.UserID) (usecase.ListItemsDTO, error)
	ClearCartByUserID(ctx context.Context, userID models.UserID) error
}

type CartServer struct {
	cartUsecase ICartUsecase
	tracer      trace.Tracer
	pb.UnimplementedCartServiceServer
}

func NewCartServer(us ICartUsecase, tracer trace.Tracer) *CartServer {
	return &CartServer{cartUsecase: us, tracer: tracer}
}

func (c *CartServer) AddItem(ctx context.Context, req *pb.CartAddItemRequest) (*emptypb.Empty, error) {
	count, err := models.Uint32ToUint16(req.Count)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	addItemDTO := usecase.AddItemDTO{
		UserID: models.UserID(req.UserId),
		SKUID:  models.SKUID(req.Sku),
		Count:  count,
	}

	if err = c.cartUsecase.AddItem(ctx, addItemDTO); err != nil {
		if errors.Is(err, usecase.ErrNotEnoughStock) {
			return nil, status.Error(codes.Aborted, err.Error())
		}

		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (c *CartServer) DeleteItem(ctx context.Context, req *pb.CartDeleteItemRequest) (*emptypb.Empty, error) {
	deleteItemDTO := usecase.DeleteItemDTO{
		UserID: models.UserID(req.UserId),
		SKUID:  models.SKUID(req.Sku),
	}

	if err := c.cartUsecase.DeleteItem(ctx, deleteItemDTO); err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (c *CartServer) ListItem(ctx context.Context, req *pb.CartUserIDRequest) (*pb.CartListItemResponse, error) {
	listDTO, err := c.cartUsecase.GetItemsByUserID(ctx, models.UserID(req.UserId))
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	respList := make([]*pb.CartItem, len(listDTO.Items))

	for i, item := range listDTO.Items {
		respItem := pb.CartItem{}

		respItem.Sku = uint32(item.SKUID)
		respItem.Name = item.Name
		respItem.Count = uint32(item.Count)
		respItem.Price = item.Price

		respList[i] = &respItem
	}

	return &pb.CartListItemResponse{Items: respList, TotalPrice: listDTO.TotalPrice}, nil
}

func (c *CartServer) ClearCart(ctx context.Context, req *pb.CartUserIDRequest) (*emptypb.Empty, error) {
	if err := c.cartUsecase.ClearCartByUserID(ctx, models.UserID(req.UserId)); err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &emptypb.Empty{}, nil
}
