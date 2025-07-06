# 📦 Stock Service

## 📦 Docker Image

```bash
docker pull umyt/my-stock-app:hw7
```

## 🚀 Application Ports

Make sure the following ports are accessible, or you can change it on your own free ports:

- App port: `SERVER_PORT=8081`
- PostgreSQL port: `DB_PORT=5434`

---

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
