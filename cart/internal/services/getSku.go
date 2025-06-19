package services

import (
	"bytes"
	"cart/internal/dto"
	"cart/internal/models"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type SkuGetServiceInterface interface {
	GetSku(ctx context.Context, skuId models.SKUID) (dto.SKU, error)
}

type SkuGetService struct {
	httpClient *http.Client
	baseUrl    string
}

func NewSkuGetService(client *http.Client, url string) *SkuGetService {
	return &SkuGetService{httpClient: client, baseUrl: url}
}

type GetSkuRequest struct {
	SkuId models.SKUID `json:"sku"`
}

type Response struct {
	Message StockResponse `json:"message"`
}

type StockResponse struct {
	SkuId    uint32 `json:"sku"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Count    uint16 `json:"count,omitempty"`
	Price    uint32 `json:"price,omitempty"`
	Location string `json:"location,omitempty"`
	UserId   int64  `json:"user_id,omitempty"`
}

func (s *SkuGetService) GetSku(ctx context.Context, skuId models.SKUID) (dto.SKU, error) {
	reqDto := GetSkuRequest{
		SkuId: skuId,
	}

	body, err := json.Marshal(&reqDto)
	if err != nil {
		return dto.SKU{}, err
	}

	responseBody := bytes.NewBuffer(body)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.baseUrl, responseBody)

	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return dto.SKU{}, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return dto.SKU{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return dto.SKU{}, err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return dto.SKU{}, err
	}

	var response Response

	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return dto.SKU{}, err
	}

	stockRes := response.Message

	sku := dto.SKU{
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
