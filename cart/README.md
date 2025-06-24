# Cart service

## Dokcer image name is

```bash
docker pull umyt/my-cart-app:hw7
```

## App port is

These port should be able to reach

app port = 8080
postgres port = 5433

## Required environment variables

In the docker-compose up you should pass my image

```docker
services:
  cart_service:
    images: umyt/my-cart-app:hw7
```

## Sample Requests or endpoints

POST - /cart/item/add

```json
{
  "user_id": 10,
  "sku": 7077,
  "count": 1
}
```

POST - /cart/item/delete

```json
{
  "user_id": 2,
  "sku": 7077
}
```

POST /cart/list

```json
{
  "user_id": 2
}
```

POST - /cart/clear

```json
{
  "user_id": 2
}
```
