package services

import (
	"bytes"
	"cart/internal/models"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type StockService interface {
	GetItemInfo(ctx context.Context, skuID models.SKUID) (ItemDTO, error)
}

type stockService struct {
	httpClient *http.Client
	baseUrl    string
}

func NewStockService(timeoutDur time.Duration, url string) StockService {
	httpClient := &http.Client{
		Timeout: timeoutDur,
	}

	return &stockService{httpClient: httpClient, baseUrl: url}
}

func (s *stockService) GetItemInfo(ctx context.Context, skuID models.SKUID) (ItemDTO, error) {
	reqData := getSKUIDRequest{
		SKUID: skuID,
	}

	body, err := json.Marshal(&reqData)
	if err != nil {
		return ItemDTO{}, err
	}

	reqBody := bytes.NewBuffer(body)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.baseUrl, reqBody)
	if err != nil {
		return ItemDTO{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return ItemDTO{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ItemDTO{}, err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return ItemDTO{}, err
	}

	var respData httpResponse

	err = json.Unmarshal(respBody, &respData)
	if err != nil {
		return ItemDTO{}, err
	}

	stock := respData.Message
	item := ItemDTO{
		SKUID:    models.SKUID(stock.SKUID),
		Name:     stock.Name,
		Type:     stock.Type,
		Count:    stock.Count,
		Price:    stock.Price,
		Location: stock.Location,
		UserID:   models.UserID(stock.UserID),
	}

	return item, nil
}
