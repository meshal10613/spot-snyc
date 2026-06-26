# 🚗 SpotSync API

> Smart Parking & EV Charging Reservation Platform

A centralized backend API for managing parking zones and handling high-demand reservation of limited EV charging spots at airports and malls. Built with Go, Echo, GORM, and PostgreSQL.

**Live URL:** `https://spotsync-api.onrender.com`

---

## ✨ Features

- **User Authentication** — Register & Login with JWT-based Bearer token auth (24h expiry)
- **Role-Based Access Control** — `driver` and `admin` roles with middleware-enforced permissions
- **Parking Zone Management** — Full CRUD for parking zones (admin) with real-time availability
- **Reservation System** — Concurrency-safe parking spot reservation with row-level locking (`SELECT ... FOR UPDATE`)
- **EV Charging Support** — Dedicated `ev_charging` zone type with capacity management
- **Query Builder** — Built-in pagination, sorting, and search across all list endpoints
- **Standardized Responses** — Consistent JSON response format using `httpresponse` package
- **Global Error Handler** — Centralized 404 Not Found, 405 Method Not Allowed, and internal error handling
- **Admin Seeder** — Auto-seeds an admin user on startup from environment variables (skips if already exists)

---

## 🛠️ Tech Stack

| Technology | Purpose |
|---|---|
| **Go 1.22+** | Backend language |
| **Echo v4** | HTTP web framework |
| **GORM** | ORM with PostgreSQL driver |
| **PostgreSQL** | Relational database |
| **JWT (golang-jwt/jwt/v5)** | Token-based authentication |
| **bcrypt** | Password hashing (cost 12) |
| **go-playground/validator** | Request validation |

---

## 🏛️ Architecture

This project follows **Clean Architecture** with strict separation of concerns:

```
┌─────────────┐     ┌─────────────┐     ┌──────────────┐     ┌──────────────┐
│   Handler    │ ──▶ │   Service   │ ──▶ │  Repository  │ ──▶ │   Database   │
│  (HTTP/JSON) │     │  (Business) │     │   (GORM)     │     │ (PostgreSQL) │
└─────────────┘     └─────────────┘     └──────────────┘     └──────────────┘
       │                    │
       ▼                    ▼
   ┌───────┐          ┌──────────┐
   │  DTO  │          │  Models  │
   │(Req/Res)│        │  (GORM)  │
   └───────┘          └──────────┘
```

**Dependency Injection** is done manually in `internal/server/http.go`:

```
Repository → Service → Handler → Routes
```

The project is organized into three top-level directories:

| Directory | Purpose |
|---|---|
| **`cmd/`** | Application entry point |
| **`internal/`** | Private application code (not importable by external projects) |
| **`pkg/`** | Reusable packages that can be imported by external projects |

---

## 📁 Project Structure

```
spot-sync/
├── cmd/
│   └── main.go                        # Entry point (config, DB, migrations, seed, server start)
│
├── internal/
│   ├── config/
│   │   ├── config.go                  # Environment variable loading
│   │   └── db.go                      # PostgreSQL connection via GORM
│   ├── dto/
│   │   ├── auth_dto.go                # Auth request/response structs
│   │   ├── reservation_dto.go         # Reservation request/response structs
│   │   └── zone_dto.go               # Zone request/response structs
│   ├── handler/
│   │   ├── auth_handler.go            # Register & Login endpoints
│   │   ├── reservation_handler.go     # Reservation CRUD endpoints
│   │   └── zone_handler.go           # Zone CRUD endpoints
│   ├── models/
│   │   ├── enum.go                    # Role, ZoneType, ReservationStatus enums
│   │   ├── migrate.go                 # Auto-migration runner
│   │   ├── parking_zone.go            # ParkingZone GORM model
│   │   ├── reservation.go             # Reservation GORM model
│   │   └── user.go                    # User GORM model with bcrypt
│   ├── repository/
│   │   ├── auth_repository.go         # User database operations
│   │   ├── reservation_repository.go  # Reservation DB ops with row locking
│   │   └── zone_repository.go         # Zone DB ops with query builder
│   ├── routes/
│   │   └── routes.go                  # All API route registration
│   ├── server/
│   │   └── http.go                    # Echo server setup, middleware, DI, global error handler
│   └── service/
│       ├── auth_service.go            # Auth business logic & JWT generation
│       ├── reservation_service.go     # Reservation business logic
│       └── zone_service.go            # Zone business logic
│
├── pkg/
│   ├── httpresponse/
│   │   └── response.go               # Standardized Success, Error, Meta structs
│   ├── middleware/
│   │   ├── jwt_auth.go                # JWT Bearer token validation
│   │   └── role_auth.go               # Role-based access control
│   ├── seed/
│   │   └── admin_seeder.go            # Auto-seeds admin user on startup
│   └── utils/
│       └── query_builder.go           # Pagination, sorting, search utility
│
├── .env.example                       # Environment template
├── .air.toml                          # Hot-reload config (Air)
├── Dockerfile                         # Multi-stage production build
└── docker-compose.yaml                # Docker setup
```

---

## 🚀 Getting Started

### Prerequisites

- Go 1.22+
- PostgreSQL database
- (Optional) [Air](https://github.com/air-verse/air) for hot-reloading

### Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/spot-sync.git
   cd spot-sync
   ```

2. **Configure environment variables**
   ```bash
   cp .env.example .env
   ```

3. **Required `.env` variables:**
   ```env
   PORT=5000
   DSN="host=localhost user=postgres password=yourpassword dbname=spotsync port=5432 sslmode=disable"
   JWT_SECRET=your-secret-key-here

   # Admin Seed Credentials
   ADMIN_NAME=Admin
   ADMIN_EMAIL=admin@gmail.com
   ADMIN_PASSWORD=Admin@123
   ```

4. **Install dependencies**
   ```bash
   go mod tidy
   ```

5. **Run the server**
   ```bash
   go run ./cmd
   ```

   Or with hot-reloading:
   ```bash
   air
   ```

### Docker

```bash
docker-compose up --build
```

### Startup Flow

```
Load Config → Connect DB → Run Migrations → Seed Admin → Start Server
```

On startup, the server automatically seeds an admin user using the `ADMIN_EMAIL` and `ADMIN_PASSWORD` from your `.env`. If the admin already exists (matched by email), the seed step is silently skipped.

---

## 🌐 API Endpoints

### Health Check

| Method | Endpoint | Access | Description |
|---|---|---|---|
| `GET` | `/` | Public | API health check |

### Authentication

| Method | Endpoint | Access | Description |
|---|---|---|---|
| `POST` | `/api/v1/auth/register` | Public | Register a new user |
| `POST` | `/api/v1/auth/login` | Public | Login and receive JWT |

### Parking Zones

| Method | Endpoint | Access | Description |
|---|---|---|---|
| `GET` | `/api/v1/zones` | Public | Get all zones (paginated) |
| `GET` | `/api/v1/zones/:id` | Public | Get a single zone |
| `POST` | `/api/v1/zones` | Admin | Create a zone |
| `PUT` | `/api/v1/zones/:id` | Admin | Update a zone |
| `DELETE` | `/api/v1/zones/:id` | Admin | Delete a zone |

### Reservations

| Method | Endpoint | Access | Description |
|---|---|---|---|
| `POST` | `/api/v1/reservations` | Auth | Reserve a parking spot |
| `GET` | `/api/v1/reservations/my-reservations` | Auth | View my reservations (paginated) |
| `DELETE` | `/api/v1/reservations/:id` | Auth | Cancel a reservation |
| `GET` | `/api/v1/reservations` | Admin | View all reservations (paginated) |

### Authentication Header

```
Authorization: Bearer <your-jwt-token>
```

---

## 📄 Query Parameters

All paginated (`GET` list) endpoints support these query parameters:

| Param | Default | Example | Description |
|---|---|---|---|
| `page` | `1` | `?page=2` | Page number (1-indexed) |
| `limit` | `10` | `?limit=20` | Items per page (max 100) |
| `sort` | `created_at` | `?sort=name` | Column to sort by |
| `order` | `desc` | `?order=asc` | Sort direction (`asc` / `desc`) |
| `search` | — | `?search=ev` | Search term (ILIKE across fields) |

**Search fields per endpoint:**

| Endpoint | Searchable Fields |
|---|---|
| `GET /api/v1/zones` | `name`, `type` |
| `GET /api/v1/reservations` | `license_plate`, `status` |
| `GET /api/v1/reservations/my-reservations` | `license_plate`, `status` |

**Example:**

```
GET /api/v1/zones?page=1&limit=5&search=ev&sort=name&order=asc
```

---

## 📦 Response Format

All API responses follow a standardized JSON structure.

### Success Response

```json
{
  "success": true,
  "message": "Parking zones retrieved successfully",
  "data": [ ... ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 25,
    "total_page": 3
  }
}
```

> `meta` is included only on paginated list endpoints. Single-resource responses omit it.

### Error Response

```json
{
  "success": false,
  "message": "Validation failed",
  "details": "Key: 'CreateZoneRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag"
}
```

### Not Found (404)

```json
{
  "success": false,
  "message": "The requested resource was not found"
}
```

### Method Not Allowed (405)

```json
{
  "success": false,
  "message": "Method not allowed"
}
```

---

## 🔒 Concurrency Safety

The reservation system uses **GORM database transactions** with **row-level locking** (`SELECT ... FOR UPDATE`) to prevent the "EV Spot Bottleneck" race condition:

```go
db.Transaction(func(tx *gorm.DB) error {
    // 1. Lock the parking zone row
    tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&zone, zoneID)
    // 2. Count active reservations
    // 3. Check capacity
    // 4. Create reservation or reject
})
```

This ensures that even if two drivers attempt to reserve the last available spot simultaneously, only one will succeed.

---

## 🌱 Environment Variables

| Variable | Required | Default | Description |
|---|---|---|---|
| `PORT` | No | `8080` | Server port |
| `DSN` | Yes | — | PostgreSQL connection string |
| `JWT_SECRET` | Yes | — | Secret key for signing JWT tokens |
| `ADMIN_NAME` | No | `Admin` | Seeded admin user's name |
| `ADMIN_EMAIL` | Yes | — | Seeded admin user's email |
| `ADMIN_PASSWORD` | Yes | — | Seeded admin user's password |

---

## 📜 License

MIT
