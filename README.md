# e-commerce-api

E-Commerce REST API on Go: JWT auth, roles (`admin`, `seller`, `buyer`), catalog, cart, orders, payments (YooKassa or mock), Swagger UI.

## Stack

- Go, Chi, sqlx, PostgreSQL, Redis (optional cache layer)
- JWT access + refresh tokens
- Swagger UI at `/swagger/index.html`

## Quick start

```bash
docker compose up -d

# Migrations (from src/)
cd src
goose -dir migration postgres "host=localhost port=5432 user=postgres password=postgres dbname=ecommerce sslmode=disable" up

# Config
cp .env.example .env

# Run API
go run ./cmd/api
```

API: `http://localhost:8080`  
Swagger: `http://localhost:8080/swagger/index.html`  
Health: `http://localhost:8080/health`

## Roles

| Role | Permissions |
|------|-------------|
| **buyer** | cart, checkout, own orders, payments |
| **seller** | manage own products |
| **admin** | users, roles, order status |

Registration creates buyer only. Promote to seller/admin:

```sql
UPDATE users SET role = 'admin' WHERE username = 'your_username';
```

## Main endpoints

| Method | Path | Auth | Role |
|--------|------|------|------|
| POST | `/api/v1/auth/register` | — | — |
| POST | `/api/v1/auth/login` | — | — |
| POST | `/api/v1/auth/refresh` | — | — |
| GET | `/api/v1/users/me` | JWT | any |
| GET | `/api/v1/products` | — | — |
| POST | `/api/v1/products` | JWT | seller, admin |
| GET | `/api/v1/cart` | JWT | buyer |
| POST | `/api/v1/orders/checkout` | JWT | buyer |
| GET | `/api/v1/admin/users` | JWT | admin |

Price is stored in kopecks 199900 = 1999.00 RUB.

## Payments

Without `YOOKASSA_SHOP_ID` / `YOOKASSA_SECRET_KEY` the API uses mock payments.

## Regenerate Swagger

```bash
cd src
make swagger   # requires: go install github.com/swaggo/swag/cmd/swag@latest
```
