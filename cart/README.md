# 🛒 Cart Service

## 📦 Docker Image

```bash
docker pull umyt/my-cart-app:hw7
```

## 🚀 Application Ports

Make sure the following ports are accessible:

- App port: `8080`
- PostgreSQL port: `5433`

---

## ⚙️ Required Environment Variables

First, create a `.env` file and add the following variables:

```env
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
SERVER_READ_HEADER_TIMEOUT=10
SERVER_SHUTDOWN_TIMEOUT=3

DB_HOST= "cart_db"
DB_PORT= 5432
DB_USER= "postgres"
DB_NAME= "cart_service_db"
DB_SSLMODE= "disable"
DB_PASSWORD= "password"

CLIENT_URL=http://stocks_service:8081/stocks/item/get
CLIENT_TIMEOUT=2
```

### ‼️ Don't Modify

- SERVER_HOST
- DB_HOST
- CLIENT_URL

Other things you can modify on your own info.

---

## 🧪 How to Test the Service

### ✅ Docker Compose

1. You need create a new docker network. If you already created, no need to create again.

```bash
docker network create public-net
```

2. Run:

```bash
docker compose up
```

---

## 📬 Sample Requests (Endpoints)

> All requests use the `POST` method.

### ➕ Add Item to Cart

`POST /cart/item/add`

```json
{
  "userId": 1,
  "sku": 1001,
  "count": 1
}
```

---

### ➖ Remove Item from Cart

`POST /cart/item/delete`

```json
{
  "userId": 1,
  "sku": 1001
}
```

---

### 📦 List Cart Items

`POST /cart/list`

```json
{
  "userId": 1
}
```

---

### 🧹 Clear Cart

`POST /cart/clear`

```json
{
  "userId": 1
}
```
