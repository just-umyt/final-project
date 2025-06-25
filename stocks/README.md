# 📦 Stock Service

## 📦 Docker Image

```bash
docker pull umyt/my-stock-app:hw7
```

## 🚀 Application Ports

Make sure the following ports are accessible:

- App port: `8081`
- PostgreSQL port: `5434`

---

## ⚙️ Required Environment Variables

First, create a `.env` file and add the following variables:

```env
SERVER_HOST=0.0.0.0
SERVER_PORT=8081
SERVER_READ_HEADER_TIMEOUT=10
SERVER_SHUTDOWN_TIMEOUT=3

DB_HOST= "stocks_db"
DB_PORT= 5432
DB_USER= "postgres"
DB_PASSWORD= "password"
DB_NAME= "stocks_service_db"
DB_SSLMODE= "disable"
```

### ‼️ Don't Modify

- SERVER_HOST
- DB_HOST

Other things you can modify on your own info.

---

## 🧪 How to Test the Service

### ✅ Docker Compose

1. You need create a new docker network

```bash
docker network create app-network
```

2. Run:

```bash
docker compose up
```

---

## 📬 Sample Requests (Endpoints)

> All requests use the `POST` method.

### ➕ Add Stock

`POST /stocks/item/add`

```json
{
  "sku": 7077,
  "user_id": 2,
  "count": 100,
  "price": 2100,
  "location": "json changed loc"
}
```

---

### 📃 Get Item from Stock

`POST /stocks/item/get`

```json
{
  "user_id": 1,
  "sku": 5055
}
```

---

### 📦 List Stock Items By Location

`POST /stocks/list/location`

```json
{
  "user_id": 123,
  "location": "json loc",
  "page_size": 2,
  "current_page": 1
}
```

---

### ➖ Stock Delete

`POST /stocks/item/delete`

```json
{
  "user_id": 123,
  "sku": 4044
}
```
