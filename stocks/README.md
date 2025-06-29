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

1. You need create a new docker network. If you already created no need to create again.

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

### ➕ Add Stock

`POST /stocks/item/add`

```json
{
  "sku": 1001,
  "userId": 1,
  "count": 10,
  "price": 100,
  "location": "AG"
}
```

---

### 📃 Get Item from Stock

`POST /stocks/item/get`

```json
{
  "userId": 1,
  "sku": 1001
}
```

---

### 📦 List Stock Items By Location

`POST /stocks/list/location`

```json
{
  "userId": 1,
  "location": "AG",
  "pageSize": 1,
  "currentPage": 1
}
```

---

### ➖ Stock Delete

`POST /stocks/item/delete`

```json
{
  "userId": 1,
  "sku": 1001
}
```
