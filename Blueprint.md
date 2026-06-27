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
│   │   ├── inventory_item.go
│   │   ├── product_ingredient.go
│   │   ├── order.go
│   │   └── cart.go
│   │
│   ├── dto/
│   │   ├── auth.go
│   │   ├── product.go
│   │   ├── inventory.go
│   │   └── order.go
│   │
│   ├── handler/
│   │   ├── auth.go
│   │   ├── product.go
│   │   ├── inventory.go
│   │   └── order.go
│   │
│   ├── service/
│   │   ├── auth.go
│   │   ├── product.go
│   │   ├── inventory.go
│   │   └── order.go
│   │
│   ├── repository/
│   │   ├── user.go
│   │   ├── product.go
│   │   ├── inventory.go
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
│   ├── 002_create_inventory_items.sql
│   ├── 003_create_products.sql
│   ├── 004_create_product_ingredients.sql
│   ├── 005_create_orders.sql
│   └── 006_create_order_items.sql
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

### InventoryItem
Raw stock tracked in the shop (beans, milk, cups, syrup, etc.).

| Field | Type | Notes |
|---|---|---|
| id | UUID | PK |
| name | string | unique among active items |
| unit | string | e.g. `g`, `ml`, `pcs` |
| quantity | decimal | NUMERIC(12, 3), current stock on hand |
| low_stock_threshold | decimal | nullable, for admin alerts |
| deleted_at | timestampz | nullable, soft delete marker |
| created_at | timestampz | |
| updated_at | timestampz | |

### Product
A menu item sold to customers. Each product is made from one or more inventory items via a recipe.

| Field | Type | Notes |
|---|---|---|
| id | UUID | PK |
| name | string | |
| description | string | nullable |
| price | decimal | NUMERIC(10, 2), never float |
| is_available | bool | false — hidden from menu |
| deleted_at | timestampz | nullable, soft delete marker |
| created_at | timestampz | |
| updated_at | timestampz | |

### ProductIngredient
Recipe line: how much of each inventory item goes into one unit of a product.

| Field | Type | Notes |
|---|---|---|
| id | UUID | PK |
| product_id | UUID | FK -> products |
| inventory_item_id | UUID | FK -> inventory_items |
| quantity | decimal | NUMERIC(12, 3), amount consumed per 1 product sold |
| created_at | timestampz | |
| updated_at | timestampz | |

Unique constraint on `(product_id, inventory_item_id)`.

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
| product_id | UUID | FK -> products |
| quantity | int | ≥ 1 |
| unit_price | decimal | snapshotted from product at order time |

### Relationships
- User -> Order: one-to-many
- Order -> OrderItem: one-to-many
- Product -> OrderItem: one-to-many
- Product -> ProductIngredient: one-to-many (recipe)
- InventoryItem -> ProductIngredient: one-to-many
- ProductIngredient: many-to-one Product, many-to-one InventoryItem

```text
InventoryItem ──< ProductIngredient >── Product ──< OrderItem
                                              │
                                              └──< Order (via OrderItem)
```

---

## Business rules

### Orders
- Customers can only view their own orders. Return `404`.
- An order must have at least one item.
- Only admins can change status.

### Products & recipes
- A sellable product must have at least one ingredient in its recipe.
- Ingredient `quantity` is always per **one** unit of the product (e.g. 18 `g` espresso beans + 250 `ml` milk for one latte).
- When creating or updating a product, admins send the full recipe as a list of `{inventory_item_id, quantity}` pairs.
- Replacing a product recipe deletes old `ProductIngredient` rows and inserts the new set in one transaction.

### Inventory
- `quantity` on an inventory item is the current stock level, not a recipe amount.
- Only admins can create, update, restock, or soft-delete inventory items.
- Restocking adds to `quantity` (e.g. receive a delivery); admins can also set `quantity` directly.
- An inventory item cannot be hard-deleted while referenced by any product recipe.

### Stock checks & deduction
- Before an order is confirmed, the service checks that every product in the cart has enough inventory:
  - required = `order_item.quantity × product_ingredient.quantity` for each ingredient
- If any ingredient would go negative, reject the order with `INSUFFICIENT_STOCK` (`422`).
- On successful order placement, deduct inventory inside the same DB transaction as order creation.
- Cancelling an order restores the deducted inventory (same transaction as status update).

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

### Inventory
| Method | Route | Auth | Description |
|---|---|---|---|
| GET | `/inventory` | admin | List all inventory items (include low-stock flag) |
| GET | `/inventory/:id` | admin | Get single inventory item |
| POST | `/inventory` | admin | Create inventory item |
| PUT | `/inventory/:id` | admin | Update name, unit, threshold |
| PATCH | `/inventory/:id/stock` | admin | Adjust stock (`delta` or absolute `quantity`) |
| DELETE | `/inventory/:id` | admin | Soft delete |

### Products
| Method | Route | Auth | Description |
|---|---|---|---|
| GET | `/products` | - | List available products |
| GET | `/products/:id` | - | Get single product with recipe (`ingredients`) |
| POST | `/products` | admin | Create product with recipe |
| PUT | `/products/:id` | admin | Update product and/or replace recipe |
| DELETE | `/product/:id` | admin | Remove product (soft) |

### Orders
| Method | Route | Auth | Description |
|---|---|---|---|
| POST | `/orders` | customer | Place order (checks stock, deducts inventory) |
| GET | `/orders` | - | customers sees own orders, admin sees all |
| GET | `/orders/:id` | - | `404` if customers tries to see others order |
| PATCH | `/orders/:id/status` | admin | Update order status (restores stock on cancel) |

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
| `INSUFFICIENT_STOCK` | 422 | Not enough inventory to fulfill order |
| `INTERNAL_ERROR` | 500 | Unexpected server error |

---

## Request / Response examples

### Create inventory item
`POST /inventory` (admin)

Request:
```json
{
  "name": "Espresso beans",
  "unit": "g",
  "quantity": "5000.000",
  "low_stock_threshold": "500.000"
}
```

Response `201`:
```json
{
  "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "name": "Espresso beans",
  "unit": "g",
  "quantity": "5000.000",
  "low_stock_threshold": "500.000",
  "is_low_stock": false,
  "created_at": "2026-06-27T10:00:00Z",
  "updated_at": "2026-06-27T10:00:00Z"
}
```

### Create product with recipe
`POST /products` (admin)

Request:
```json
{
  "name": "Latte",
  "description": "Double shot with steamed milk",
  "price": "4.50",
  "is_available": true,
  "ingredients": [
    { "inventory_item_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890", "quantity": "18.000" },
    { "inventory_item_id": "b2c3d4e5-f6a7-8901-bcde-f12345678901", "quantity": "250.000" }
  ]
}
```

Response `201`:
```json
{
  "id": "c3d4e5f6-a7b8-9012-cdef-123456789012",
  "name": "Latte",
  "description": "Double shot with steamed milk",
  "price": "4.50",
  "is_available": true,
  "ingredients": [
    {
      "inventory_item_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "name": "Espresso beans",
      "unit": "g",
      "quantity": "18.000"
    },
    {
      "inventory_item_id": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
      "name": "Whole milk",
      "unit": "ml",
      "quantity": "250.000"
    }
  ],
  "created_at": "2026-06-27T10:05:00Z",
  "updated_at": "2026-06-27T10:05:00Z"
}
```

### Place order (insufficient stock)
`POST /orders` (customer)

Request:
```json
{
  "items": [
    { "product_id": "c3d4e5f6-a7b8-9012-cdef-123456789012", "quantity": 100 }
  ]
}
```

Response `422`:
```json
{
  "error": {
    "code": "INSUFFICIENT_STOCK",
    "message": "Not enough Whole milk to fulfill order",
    "details": [
      {
        "inventory_item_id": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
        "name": "Whole milk",
        "required": "25000.000",
        "available": "1200.000"
      }
    ]
  }
}
```

---
