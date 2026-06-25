# Bardista (Coffee Shop API) Blueprint

## Stack

- **Language**: Go 1.26
- **Router**: Gin
- **Database**: PostgreSQL 18
- **Auth**: JWT (`golang-jwt/jwt/v5`)
- **DB Driver**: `jackx/pgx/v5`
- **Password hashing**: `golang.org/x/crypto/bcrypt`
- **UUIDs**: `google/uuid`

## Project Stucture

```
bardista/
├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   ├── config/
│   │   └── config.go
│   │
│   ├── database/
│   │   └── postgres.go
│   │
│   ├── domain/
│   │   ├── user.go
│   │   ├── product.go
│   │   ├── order.go
│   │   └── cart.go
│   │
│   ├── dto/
│   │   ├── auth.go
│   │   ├── product.go
│   │   └── order.go
│   │
│   ├── handler/
│   │   ├── auth.go
│   │   ├── product.go
│   │   └── order.go
│   │
│   ├── service/
│   │   ├── auth.go
│   │   ├── product.go
│   │   └── order.go
│   │
│   ├── repository/
│   │   ├── user.go
│   │   ├── product.go
│   │   └── order.go
│   │
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── cors.go
│   │   └── logger.go
│   │
│   ├── router/
│   │   └── router.go
│   │
│   └── utils/
│       ├── jwt.go
│       ├── password.go
│       └── response.go
│
├── migrations/
│   ├── 001_create_users.sql
│   ├── 002_create_products.sql
│   ├── 003_create_orders.sql
│   └── 004_create_order_items.sql
│
├── .env
├── .env.example
├── go.mod
└── go.sum
```
