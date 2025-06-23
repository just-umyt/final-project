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
	GetItemInfo(ctx context.Context, skuId models.SKUID) (Item, error)
}

type stockService struct {
	httpClient *http.Client
	baseUrl    string
}

func NewSkuGetService(timeoutDur time.Duration, url string) StockService {
	httpClient := &http.Client{
		Timeout: timeoutDur,
	}

	return &stockService{httpClient: httpClient, baseUrl: url}
}

func (s *stockService) GetItemInfo(ctx context.Context, skuId models.SKUID) (Item, error) {
	reqDto := GetSkuRequest{
		SkuId: skuId,
	}

	body, err := json.Marshal(&reqDto)
	if err != nil {
		return Item{}, err
	}

	responseBody := bytes.NewBuffer(body)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.baseUrl, responseBody)

	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return Item{}, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return Item{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Item{}, err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return Item{}, err
	}

	var response Response

	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return Item{}, err
	}

	stockRes := response.Message

	sku := Item{
		SkuId:    models.SKUID(stockRes.SkuId),
		Name:     stockRes.Name,
		Type:     stockRes.Type,
		Count:    stockRes.Count,
		Price:    stockRes.Price,
		Location: stockRes.Location,
		UserId:   models.UserID(stockRes.UserId),
	}

	return sku, nil
}
