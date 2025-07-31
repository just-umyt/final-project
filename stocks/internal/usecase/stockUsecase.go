package usecase

import (
	"context"
	"errors"
	"stocks/internal/models"
	"stocks/internal/producer"
	"stocks/internal/repository"
	"time"

	myLog "stocks/internal/observability/log"

	"go.opentelemetry.io/otel"
)

const (
	eventSKUCreateType   = "sku_created"
	eventStockChangeType = "stock_changed"

	eventService = "stock"

	topic = "metrics"

	tracingServiceName = "stock-service"
	addSpanName        = "stock-add-usecase"
	delSpanName        = "stock-del-usecase"
	listSpanName       = "stock-list-usecase"
	getSpanName        = "stock-get-usecase"
)

var (
	ErrNotFound error = errors.New("not found")
	ErrUserID   error = errors.New("user id is not matched")
)

//go:generate mkdir -p mock
//go:generate minimock -o ./mock -s .go  -g
type IPgTxManager interface {
	WithTx(ctx context.Context, fn func(repository.IStockRepo) error) error
}

type IProducer interface {
	Produce(messsageDTO producer.ProducerMessageDTO, topic string, t time.Time) error
}

type StockUsecase struct {
	stockRepo     repository.IStockRepo
	trManager     IPgTxManager
	kafkaProducer IProducer
	logger        myLog.Logger
}

func NewStockUsecase(repo repository.IStockRepo, trManager IPgTxManager, kafkaPr IProducer, logg myLog.Logger) *StockUsecase {
	return &StockUsecase{stockRepo: repo, trManager: trManager, kafkaProducer: kafkaPr, logger: logg}
}

func (u *StockUsecase) AddStock(ctx context.Context, stock AddStockDTO) error {
	ctx, span := otel.Tracer(tracingServiceName).Start(ctx, addSpanName)
	defer span.End()

	messageDTO := producer.ProducerMessageDTO{
		Service:   eventService,
		Timestamp: time.Now(),
	}

	if err := u.trManager.WithTx(ctx, func(repo repository.IStockRepo) error {
		item, err := repo.GetItemBySKU(ctx, stock.SKUID)
		if err != nil {
			if item.Stock.ID == 0 {
				return ErrNotFound
			}

			return err
		}

		newItem := models.Stock{
			Count:    item.Stock.Count + stock.Count,
			Price:    stock.Price,
			Location: stock.Location,
			UserID:   stock.UserID,
			SKUID:    stock.SKUID,
		}

		switch item.Stock.UserID {
		case 0:
			err := repo.AddStock(ctx, newItem)
			if errors.Is(err, repository.ErrNotFound) {
				return ErrNotFound
			}

			messageDTO.Type = eventSKUCreateType
			messageDTO.SKU = newItem.SKUID
			messageDTO.Count = newItem.Count
			messageDTO.Price = newItem.Price

			return err
		case stock.UserID:
			err := repo.UpdateStock(ctx, newItem)
			if errors.Is(err, repository.ErrNotFound) {
				return ErrNotFound
			}

			messageDTO.Type = eventStockChangeType
			messageDTO.SKU = newItem.SKUID
			messageDTO.Count = newItem.Count
			messageDTO.Price = newItem.Price

			return err
		default:
			return ErrUserID
		}
	}); err != nil {
		return err
	}

	u.logger.Info("kafka", myLog.Error(u.kafkaProducer.Produce(messageDTO, topic, time.Now())))

	return nil
}

func (u *StockUsecase) DeleteStockBySKU(ctx context.Context, delStock DeleteStockDTO) error {
	ctx, span := otel.Tracer(tracingServiceName).Start(ctx, delSpanName)
	defer span.End()

	err := u.stockRepo.DeleteStock(ctx, delStock.SKUID, delStock.UserID)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrNotFound
	}

	return err
}

func (u *StockUsecase) GetStocksByLocation(ctx context.Context, param GetItemByLocDTO) (ItemsByLocDTO, error) {
	ctx, span := otel.Tracer(tracingServiceName).Start(ctx, listSpanName)
	defer span.End()

	var items ItemsByLocDTO

	limit := param.PageSize
	offset := limit * (param.CurrentPage - 1)

	params := repository.GetStockByLocation{
		UserID:   param.UserID,
		Location: param.Location,
		Limit:    limit,
		Offset:   offset,
	}

	err := u.trManager.WithTx(ctx, func(repo repository.IStockRepo) error {
		stocksFromRepo, err := repo.GetItemsByLocation(ctx, params)
		if err != nil {
			return err
		}

		for _, s := range stocksFromRepo {
			item := StockDTO{
				SKU: SKUDTO{
					SKUID: s.SKU.ID,
					Name:  s.SKU.Name,
					Type:  s.SKU.Type,
				},
				Price:    s.Stock.Price,
				Count:    s.Stock.Count,
				Location: s.Stock.Location,
				UserID:   s.Stock.UserID,
			}

			items.Stocks = append(items.Stocks, item)
		}

		return nil
	})

	items.TotalCount = len(items.Stocks)
	items.PageNumber = param.CurrentPage

	return items, err
}

func (u *StockUsecase) GetItemBySKU(ctx context.Context, sku models.SKUID) (StockDTO, error) {
	ctx, span := otel.Tracer(tracingServiceName).Start(ctx, getSpanName)
	defer span.End()

	var stockDTO StockDTO
	err := u.trManager.WithTx(ctx, func(repo repository.IStockRepo) error {
		item, err := repo.GetItemBySKU(ctx, sku)
		if err != nil {
			if item.SKU.ID == 0 {
				return ErrNotFound
			} else {
				return err
			}
		}

		stockDTO = StockDTO{
			SKU: SKUDTO{
				SKUID: item.SKU.ID,
				Name:  item.SKU.Name,
				Type:  item.SKU.Type,
			},
			Price:    item.Stock.Price,
			Count:    item.Stock.Count,
			Location: item.Stock.Location,
			UserID:   item.Stock.UserID,
		}

		return nil
	})

	return stockDTO, err
}
