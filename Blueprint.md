# Bardista (Coffee Shop API) Blueprint

## Stack

- **Language**: Go 1.26
- **Router**: Gin
- **Database**: PostgreSQL 18
- **Auth**: JWT (`golang-jwt/jwt/v5`)
- **DB Driver**: `jackx/pgx/v5`
- **Password hashing**: `golang.org/x/crypto/bcrypt`
- **UUIDs**: `google/uuid`

---

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

---

## Layer Architecture

```text
HTTP Request
     │
     ▼
 Handler
     │
     ▼
 Service
     │
     ▼
 Repository
     │
     ▼
 Database
```

- `handler` — parse HTTP request, call service method, write response. No SQL.
- `service` — business rules, transactions, calculations. No HTTP, no SQL.
- `repository` — one function per SQL query. Returns domain structs or typed errors.

---

## Entities

### User
| Field | Type | Notes |
|---|---|---|
| id | UUID | PK |
| name | string | |
| email | string | unique, used when logging in |
| password_hash | string | bcrypt, never returned in responses |
| role | enum | `customer` or `admin` |
| created_at | timestampz | |
| updated_at | timestampz | |

### Product
| Field | Type | Notes |
|---|---|---|
| id | UUID | PK |
| name | string | |
| description | string | nullable |
| price | decimal | NUMERIC(10, 2), never float |
| is_available | bool | false - hidden |
| deleted_at | timestampz | nullable, soft delete marker |
| created_at | timestampz | |
| updated_at | timestampz | |

### Order
| Field | Type | Notes |
|---|---|---|
| id | UUID | PK |
| user_id | UUID | FK -> users |
| status | enum | `confirmed`, `pending`, `completed`, `cancelled` |
| total_price | decimal | frozen at order time |
| created_at | timestampz | |
| updated_at | timestampz | |

### OrderItem 
| Field | Type | Notes |
|---|---|---|
| id | UUID | PK |
| order_id | UUID | FK -> orders |
| product_it | UUID | FK -> products |
| quantity | int | ≥ 1 |
| unit_price | decimal | snapshotted from product at order time |

### Relationships
- User -> Order: one-to-many
- Order -> OrderItem: one-to-many
- Product -> OrderItem: one-to-many

---

## Business rules

- Customers can only view their own orders. Return `404`.
- An order must have at least one item.
- Only admins can change status.

---

## Database Schema

TODO

---

## API Endpoints

Base URL: `api/v1`

### Auth
| Method | Route | Auth | Description |
|---|---|---|---|
| POST | `/auth/register` | - | Create an account |
| POST | `/auth/login` | - | Returns `{token, user}` |

### Products
| Method | Route | Auth | Description |
|---|---|---|---|
| GET | `/products` | - | List of available products |
| GET | `/products/:id` | - | Get single product |
| POST | `/products` | admin | Create product |
| PUT | `/products/:id` | admin | Update product |
| DELETE | `/product/:id` | admin | Remove product (soft) |

### Orders
| Method | Route | Auth | Description |
|---|---|---|---|
| POST | `/orders` | customer | Place order |
| GET | `/orders` | - | customers sees own orders, admin sees all |
| GET | `/orders/:id` | - | `404` if customers tries to see others order |
| PATCH | `/orders/:id/status` | admin | Update order status |

---

## Error Codes

All errors use this format:
```json
{ "error": {"code": "ERROR_CODE", "message": "msg"}}
```

| CODE | HTTP status | When |
|---|---|---|
| `INVALID_INPUT` | 400 | |
| `UNAUTHORIZED` | 401 | |
| `FORBIDDEN` | 403 | |
| `NOT_FOUND` | 404 | |
| `CONFLICT` | 409 | |
| `UNPROCESSABLE` | 422 | |
| `INTERNAL_ERROR` | 500 | Unexpected server error |

---

## Request / Response examples

---



