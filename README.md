# MacBook Marketplace

An online shop for used and refurbished MacBooks, replacement parts and accessories.
The backend is a Go REST API built only on the standard library plus a YAML parser.
A plain HTML/CSS/JavaScript website, served by the same Go program, uses that API.

## Run it

Requires Go 1.22 or newer (the project uses 1.27).

1. Create your local settings file and choose an admin password:
   ```
   copy .env.example .env        (Windows)
   cp .env.example .env          (macOS / Linux)
   ```
   Then edit `.env` and set `ADMIN_PASSWORD`. The `.env` file is git-ignored, so the password never
   goes into the repository. A real `ADMIN_PASSWORD` environment variable works too.
2. Start the server:
   ```
   go run .
   ```

| | URL |
|---|---|
| Shop | http://localhost:8080 |
| Admin | http://localhost:8080/admin.html (user `admin`, password from `.env`) |
| API health check | http://localhost:8080/api/health |

The catalog is read from [`data/catalog.yml`](data/catalog.yml) when the server starts.
Orders and admin changes are kept in memory, so they reset when the server restarts.

## Features

- Browse products by category, search by name or description, filter by condition and stock.
- Product details with specifications and compatible models.
- Cart and checkout. An order with several products is all-or-nothing: if one line can't
  be fulfilled, no stock changes.
- Cancel an order, which puts its stock back.
- Admin: add, edit and delete products, and restock them.
- Safe under concurrent requests: two customers can't buy the same last unit.

## Project structure

```
main.go                  Reads settings, loads the catalog, connects the packages, starts the server
.env.example             Settings template (copy to .env, which is git-ignored)
data/catalog.yml         Categories and products
internal/
  catalog/               Products, categories and stock rules (knows nothing about HTTP)
    models.go              Product, ProductDetails, Category, Condition
    catalog.go             Store: all product data behind a lock
    stock.go               Sell, Restock, Stock, clone
    filter.go              Filter for listing products
    load.go                Reads and validates the YAML file
    errors.go              Named errors (ErrProductNotFound, ErrOutOfStock, ...)
  order/                 Placing, listing and cancelling orders
    order.go               Order, Item, Customer, Address, Status
    service.go             Service: all orders behind a lock, plus validation
  api/                   HTTP: routes, handlers, JSON in and out
    routes.go              URL → handler table
    handlers.go            Category and product handlers, query filters
    admin.go               Add, edit, delete and restock handlers
    orders.go              Order handlers and their JSON shapes
    responses.go           Product and category JSON shapes
    errors.go              Maps named errors to HTTP status codes
    auth.go                AdminOnly middleware (admin login)
    json.go                writeJSON, writeError, decodeJSON, pathID
web/                     Shop page and admin page (HTML, CSS, JavaScript)
```

Dependencies point one way: `main → api → order → catalog`.
`catalog` imports none of the others, so the business rules can be used and tested without HTTP.

## How a request flows

```
Browser  ── GET /api/products/6 ──►  router (api/routes.go)
                                      └─► server.handleGetProduct (api/handlers.go)
                                            ├─ read and check the id     (400 if not a number)
                                            ├─ store.Product(6)          catalog: lock → map lookup → copy → unlock
                                            │                            (404 if missing)
                                            └─ newProductResponse → writeJSON (200 + JSON)
Browser  ◄── {"id":6,"name":"Battery for MacBook Air 13-inch M2",...}
```

## API

All requests and responses are JSON. Prices are whole numbers in **paise**
(`12999900` = ₹1,29,999), which avoids floating-point rounding errors.
Errors always look like `{"error": "message"}`.

### Catalog

| Method | Path | Description |
|---|---|---|
| GET | `/api/health` | `{"status":"ok"}` |
| GET | `/api/categories` | All categories |
| GET | `/api/categories/{id}/products` | Products in one category (accepts the filters below) |
| GET | `/api/products` | All products. Filters: `?search=air`, `?condition=new\|refurbished\|used`, `?inStock=true`, `?categoryId=2` |
| GET | `/api/products/{id}` | One product |

### Admin

These endpoints, the order endpoints except placing an order, and the admin page itself need the
admin login (HTTP Basic Auth). Without it they reply `401 Unauthorized`. With curl, add
`-u admin:<your password>`. In a browser, the login box appears when you open the admin page.

| Method | Path | Body | Description |
|---|---|---|---|
| POST | `/api/products` | product details + `stock` | Add a product. Replies 201 with its new ID |
| PUT | `/api/products/{id}` | product details | Edit details. Stock can't be set here (see below) |
| DELETE | `/api/products/{id}` | | Remove a product. Replies 204 |
| POST | `/api/products/{id}/restock` | `{"quantity": 10}` | Add stock |

Product details: `name`, `description`, `categoryId`, `condition`, `price`, `specs` (object), `compatibleModels` (array).

### Orders

| Method | Path | Description |
|---|---|---|
| POST | `/api/orders` | Place an order. Replies 201. Public |
| GET | `/api/orders` | All orders. Admin |
| GET | `/api/orders/{id}` | One order. Admin |
| POST | `/api/orders/{id}/cancel` | Cancel and put the stock back. Admin |

Example order:

```json
{
  "customer": {
    "name": "Asha Rao",
    "phone": "9876543210",
    "address": { "line1": "12 MG Road", "city": "Bengaluru", "state": "Karnataka", "pincode": "560001" }
  },
  "items": [
    { "productId": 3, "quantity": 2 },
    { "productId": 7, "quantity": 1 }
  ]
}
```

With curl (in PowerShell, write `curl.exe`):

```
curl -X POST http://localhost:8080/api/orders -H "Content-Type: application/json" -d @order.json
```

### Status codes

| Code | Meaning | Examples |
|---|---|---|
| 200 | OK | Any successful read or update |
| 201 | Created | Product added, order placed |
| 204 | No Content | Product deleted |
| 400 | The request itself is wrong | Invalid JSON, unknown field, quantity 0, 9-digit phone |
| 401 | Admin login missing or wrong | Restocking without the password |
| 404 | It doesn't exist | Product #99, order #99 |
| 405 | Wrong method for this path | `DELETE /api/orders` |
| 409 | Valid, but clashes with the current state | Not enough stock, order already cancelled |

## Design decisions

**Private fields protect the business rules.** Other packages can only use the exported methods:

| Field | Why it is private |
|---|---|
| `Product.stock` | Changes only through `Sell` and `Restock`, so it can never go negative or oversell |
| `Store.products`, `Store.mu` | The data can only be reached through methods that hold the lock |
| `Store.nextID` | New product IDs are never reused |
| `Order.status` | Changes only through `Service.Cancel`, which also restocks. Cancelling twice is refused |
| `Order.total` | Calculated once from the items, so it always matches them |

**A lock (`sync.Mutex`) for concurrent requests.** Go's HTTP server handles each request in its own
goroutine, so `catalog.Store` and `order.Service` lock before reading or changing their data.

**Methods return copies.** `Store.Product` returns a copy, not a pointer into the store, so nobody
can change product data after the lock is released.

**Separate JSON types.** `api` has its own request and response structs (`productResponse`,
`placeOrderRequest`, ...). The business types stay free of JSON concerns, and the API's shape doesn't
change by accident when they do.

**Stock is never typed in.** `ProductDetails`, which is what an admin edits, has no stock field.
Stock changes only through real events (sale, cancellation, restock), so a stale edit can't
overwrite sales made in the meantime.

**Orders keep a snapshot.** Each order item copies the product's name and price at the time of
sale, so later price changes or deletions don't change past orders.

**Data is checked once, at startup.** `LoadFile` refuses unknown fields (typos), duplicate IDs,
unknown categories, invalid conditions and negative prices or stock, and the server won't start
with bad data. The same `ProductDetails.validate` rules are used by the admin API.

## Limitations and next steps

- **Persistence:** data lives in memory. Next step: PostgreSQL. A small interface in `catalog`
  would let the store load from either YAML or the database.
- **Authentication:** admin uses one password (from `.env`) with HTTP Basic Auth, which is fine
  for a demo on localhost. A real deployment needs HTTPS, user accounts with hashed passwords, and
  customer logins so customers can see their own orders.
- **Tests:** automated tests for `catalog`, `order` and the HTTP handlers (`net/http/httptest`).
