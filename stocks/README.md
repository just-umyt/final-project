# Stocks service

## Dokcer image name is

```bash
docker pull umyt/my-stock-app:hw7
```

## App port is

These port should be able to reach

app port = 8081
postgres port = 5434

## Required environment variables

In the docker-compose up you should pass my image

```docker
services:
  stocks_service:
    images: umyt/my-stock-app:hw7
```

## Sample Requests or endpoints

POST /stocks/item/add

```json
{
  "sku": 7077,
  "user_id": 2,
  "count": 100,
  "price": 2100,
  "location": "json changed loc"
}
```

POST /stocks/item/get

```json
{
  "user_id": 2,
  "sku": 7077
}
```

POST /stocks/item/delete

```json
{
  "user_id": 123,
  "sku": 4044
}
```

POST /stocks/list/location

```json
{
  "user_id": 123,
  "location": "json loc",
  "page_size": 2,
  "current_page": 1
}
```
