# 🛒 Cart Service

## 📦 Docker Image

```bash
docker pull umyt/my-cart-app:hw7
```

## 🚀 Application Ports

Make sure the following ports are accessible, or you can change it on your own free ports:

- App port: `SERVER_PORT=8080`
- PostgreSQL port: `DB_PORT=5433`

---

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
