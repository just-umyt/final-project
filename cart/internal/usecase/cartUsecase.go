package usecase

import (
	"cart/internal/models"
	"cart/internal/producer"
	"cart/internal/repository"
	"cart/internal/services"
	"context"
	"errors"
	"log"
	"time"
)

//go:generate mkdir -p mock
//go:generate minimock -o ./mock/ -s .go  -g
type IPgTxManager interface {
	WithTx(ctx context.Context, fn func(repository.ICartRepo) error) error
}

type IStockService interface {
	GetItemInfo(ctx context.Context, skuID models.SKUID) (services.ItemDTO, error)
}

type IProducer interface {
	Produce(messsageDTO producer.ProducerMessageDTO, topic string, t time.Time) error
}

type CartUsecase struct {
	skuService    IStockService
	cartRepo      repository.ICartRepo
	trManager     IPgTxManager
	kafkaProducer IProducer
}

const (
	eventSuccessType = "cart_item_added"
	eventFailedType  = "cart_item_failed"

	eventStatusOk     = "success"
	eventStatusFailed = "failed"

	eventService = "cart"

	topic = "metrics"
)

var (
	ErrNotFound       error = errors.New("not found")
	ErrNotEnoughStock error = errors.New("not enough stock")
)

func NewCartUsecase(cartRepo repository.ICartRepo, trManager IPgTxManager, service IStockService, kafkaPr IProducer) *CartUsecase {
	return &CartUsecase{cartRepo: cartRepo, trManager: trManager, skuService: service, kafkaProducer: kafkaPr}
}

func (u *CartUsecase) AddItem(ctx context.Context, addItem AddItemDTO) error {
	item, err := u.skuService.GetItemInfo(ctx, addItem.SKUID)
	if err != nil {
		return err
	}

	id, err := u.cartRepo.GetCartID(ctx, addItem.UserID, addItem.SKUID)
	if err != nil {
		return err
	}

	messageDTO := producer.ProducerMessageDTO{
		Type:      eventSuccessType,
		Service:   eventService,
		Timestamp: time.Now(),
		CartID:    id,
		SKU:       addItem.SKUID,
		Count:     addItem.Count,
		Status:    eventStatusOk,
	}

	if item.Count < addItem.Count {
		messageDTO.Type = eventFailedType
		messageDTO.Status = eventStatusFailed
		messageDTO.Reason = ErrNotEnoughStock.Error()

		log.Println(u.kafkaProducer.Produce(messageDTO, topic, time.Now()))

		return ErrNotEnoughStock
	}

	if err = u.trManager.WithTx(ctx, func(repo repository.ICartRepo) error {
		cart := models.Cart{
			ID:     id,
			UserID: addItem.UserID,
			SKUID:  addItem.SKUID,
			Count:  addItem.Count,
		}

		if id > 0 {
			err = repo.UpdateItemByUserID(ctx, cart)
			if errors.Is(err, repository.ErrNotFound) {
				return ErrNotFound
			}

			return err
		}

		return repo.AddItem(ctx, cart)
	}); err != nil {
		return err
	}

	log.Println(u.kafkaProducer.Produce(messageDTO, topic, time.Now()))

	return nil
}

func (u *CartUsecase) DeleteItem(ctx context.Context, delItem DeleteItemDTO) error {
	err := u.cartRepo.DeleteItem(ctx, delItem.UserID, delItem.SKUID)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrNotFound
	}

	return err
}

func (u *CartUsecase) GetItemsByUserID(ctx context.Context, userID models.UserID) (ListItemsDTO, error) {
	var list ListItemsDTO

	carts, err := u.cartRepo.GetCartByUserID(ctx, userID)
	if err != nil {
		return list, err
	}

	for _, cart := range carts {
		sku, err := u.skuService.GetItemInfo(ctx, cart.SKUID)
		if err != nil {
			return ListItemsDTO{}, err
		}

		realCount := cart.Count

		if cart.Count > sku.Count {
			log.Printf("Warning: user requested %d of SKU %d, but only %d in stock. Adjusting.",
				cart.Count, cart.SKUID, sku.Count)
			realCount = sku.Count
		}
		sku.Count = realCount
		list.Items = append(list.Items, sku)
		list.TotalPrice += uint32(realCount) * sku.Price
	}

	return list, nil
}

func (u *CartUsecase) ClearCartByUserID(ctx context.Context, userID models.UserID) error {
	return u.trManager.WithTx(ctx, func(repo repository.ICartRepo) error {
		err := repo.ClearCartByUserID(ctx, userID)
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}

		return err
	})
}
