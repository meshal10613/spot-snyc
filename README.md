# 🚗 SpotSync API

> **Smart Parking & EV Charging Reservation Platform**
>
> A robust, high-performance, and concurrency-safe backend REST API designed to manage parking zones and handle high-demand reservations for limited EV charging and standard spots at airports, malls, and public zones.

---

## ⚡ Key Highlights & Features

- **🛡️ Secure Authentication** — Custom JWT-based user registration and login with bcrypt password hashing (cost factor: 12) and 24-hour token expiration.
- **👥 Role-Based Access Control (RBAC)** — Distinct access controls for `driver` and `admin` roles, enforced via lightweight middleware.
- **📦 Clean Architecture & Modular Layout** — Separation of concerns through a professional Go layout dividing code into `internal` (private app domain) and `pkg` (sharable packages).
- **🔒 Pessimistic Concurrency Controls** — Prevents double-bookings and race conditions for high-demand parking slots using strict row-level locking (`SELECT ... FOR UPDATE`) in GORM transactions.
- **🔍 Dynamic Query Builder** — Modular pagination, sorting, and ILIKE search filters parsed directly from HTTP queries and applied to GORM database adapters.
- **🎯 Global Exception Handler** — Centrally handles unhandled execution panics, 404 (Not Found), 405 (Method Not Allowed), and struct validation errors (`go-playground/validator`).
- **🌱 Automated Admin Provisioning** — Bootstraps/seeds the super-admin account on startup using environment variables if it does not already exist.

---

## 🛠️ Technology Stack

| Technology | Purpose | Description |
| :--- | :--- | :--- |
| **Go 1.22+** | Language | High-concurrency systems programming language |
| **Echo v4** | HTTP Web Framework | High performance, extensible, and minimalist Go web framework |
| **GORM** | Database ORM | Developer-friendly ORM with clean transaction syntax |
| **PostgreSQL** | Database | Relational database ideal for transactional consistency and locking |
| **JWT** | Token Auth | Signed tokens using `golang-jwt/jwt/v5` for stateless security |
| **Bcrypt** | Encryption | Secure password hashing using `golang.org/x/crypto/bcrypt` |
| **Validator** | Data Validation | Struct field verification using `go-playground/validator/v10` |

---

## 🏛️ Architecture & Project Layout

The codebase implements **Clean Layered Architecture** principles, maintaining a strict unidirectional flow of dependencies:

```
┌─────────────────┐       ┌─────────────────┐       ┌─────────────────┐       ┌─────────────────┐
│     Handler     │ ────▶ │     Service     │ ────▶ │   Repository    │ ────▶ │   PostgreSQL    │
│  (HTTP & DTOs)  │       │ (Business Logic)│       │ (Data Access)   │       │  (Persistence)  │
└─────────────────┘       └─────────────────┘       └─────────────────┘       └─────────────────┘
         │                         │
         ▼                         ▼
┌─────────────────┐       ┌─────────────────┐
│       DTO       │       │     Models      │
│ (Request/Resp)  │       │  (GORM Entity)  │
└─────────────────┘       └─────────────────┘
```

### Dependency Injection Pipeline
All dependencies are manually wired in the server bootstrapping process located in `internal/server/http.go`:
```
Repository ────▶ Service ────▶ Handler ────▶ HTTP Routes (Echo)
```

---

## 📁 Directory Structure

The project conforms to modern Go structure guidelines by separating private API logic (`internal/`) from reuseable packages (`pkg/`):

```
spot-sync/
├── cmd/
│   └── main.go                     # App entry point (loads env, connects DB, migrates, seeds, starts HTTP)
├── internal/
│   ├── config/
│   │   ├── config.go               # Config struct & environment parser
│   │   └── db.go                   # GORM PostgreSQL connector setup
│   ├── dto/
│   │   ├── auth_dto.go             # Registration/Login request/response payloads
│   │   ├── reservation_dto.go      # Reservation request/response payloads
│   │   └── zone_dto.go             # Zone request/response payloads
│   ├── handler/
│   │   ├── auth_handler.go         # Authentication controller
│   │   ├── reservation_handler.go  # Reservation controller
│   │   └── zone_handler.go         # Parking zone controller
│   ├── models/
│   │   ├── enum.go                 # Shared Enums (Role, ZoneType, ReservationStatus)
│   │   ├── migrate.go              # Auto-migration runner using GORM schema synchronization
│   │   ├── parking_zone.go         # ParkingZone database model
│   │   ├── reservation.go          # Reservation database model
│   │   └── user.go                 # User database model
│   ├── repository/
│   │   ├── auth_repository.go      # User storage interactions
│   │   ├── reservation_repository.go # Concurrency-safe reservation transactions
│   │   └── zone_repository.go      # Parking zone database operations
│   ├── routes/
│   │   └── routes.go               # Route groups, permissions, and handlers mapping
│   ├── server/
│   │   └── http.go                 # Echo server setup (Validator, global middleware, DI, Error Handler)
│   └── service/
│       ├── auth_service.go         # Authentication and token issuance business logic
│       ├── reservation_service.go  # Reservation business rules
│       └── zone_service.go         # Parking zone business rules
├── pkg/
│   ├── httpresponse/
│   │   └── response.go             # Standardized JSend-like Success and Error JSON payloads
│   ├── middleware/
│   │   ├── jwt_auth.go             # Bearer Token extraction & verification middleware
│   │   └── role_auth.go            # Role check (driver/admin) validation middleware
│   ├── seed/
│   │   └── admin_seeder.go         # Admin user database seeder
│   └── utils/
│       └── query_builder.go        # Dynamic pagination, sorting, and filter builder
├── .air.toml                       # Hot-reloading daemon configuration
├── .env.example                    # Sample environment template
├── Dockerfile                      # Multistage production Docker configuration
└── docker-compose.yaml             # Docker Compose orchestrator file
```

---

## 🚀 Getting Started

### Prerequisites
- **Go**: 1.22+
- **PostgreSQL**: 14+
- **Air** (Optional): `go install github.com/air-verse/air@latest` (for hot reloading)

### Setup & Installation

1. **Clone the Repository:**
   ```bash
   git clone https://github.com/your-username/spot-sync.git
   cd spot-sync
   ```

2. **Configure Environment Variables:**
   Copy the example environment file and configure the settings for your database:
   ```bash
   cp .env.example .env
   ```
   Edit `.env`:
   ```env
   PORT=5000
   DSN="host=localhost user=postgres password=yourpassword dbname=spotsync port=5432 sslmode=disable TimeZone=Asia/Dhaka"
   JWT_SECRET=your-secure-jwt-secret-key-here
   
   # Admin Seed Credentials
   ADMIN_NAME=Admin
   ADMIN_EMAIL=admin@gmail.com
   ADMIN_PASSWORD=Admin@123
   ```

3. **Install Dependencies:**
   ```bash
   go mod tidy
   ```

4. **Launch the Application:**
   * **Standard Go execution:**
     ```bash
     go run cmd/main.go
     ```
   * **Hot-reload execution (Development):**
     ```bash
     air
     ```
   * **Docker Compose execution:**
     ```bash
     docker-compose up --build
     ```

### App Startup Lifecycle
On starting, the application undergoes the following sequence:
```
[Parse Env Variables] ──▶ [Init DB Conn] ──▶ [Execute GORM Schema Migration] ──▶ [Auto-Seed Admin] ──▶ [Start HTTP Server]
```
> [!NOTE]
> The admin account is seeded conditionally. If a user matching `ADMIN_EMAIL` exists in the database, the step is skipped; otherwise, the account is created using the environment settings.

---

## 🌐 API Route Specification

### HTTP Header Requirements
For protected routes, include the JWT token in your request headers:
```http
Authorization: Bearer <your_jwt_token_here>
```

---

### Endpoints Table

| Category | HTTP Method | Route | Access | Description |
| :--- | :--- | :--- | :--- | :--- |
| **System** | `GET` | `/` | Public | Api Health check |
| **Auth** | `POST` | `/api/v1/auth/register` | Public | Register a driver account |
| **Auth** | `POST` | `/api/v1/auth/login` | Public | Sign in to acquire JWT token |
| **Zones** | `GET` | `/api/v1/zones` | Public | List all zones (Paginated, Searchable) |
| **Zones** | `GET` | `/api/v1/zones/:id` | Public | Fetch a specific zone by ID |
| **Zones** | `POST` | `/api/v1/zones` | Admin | Create a new zone |
| **Zones** | `PUT` | `/api/v1/zones/:id` | Admin | Update parking zone details |
| **Zones** | `DELETE`| `/api/v1/zones/:id` | Admin | Delete a parking zone |
| **Reservations** | `POST` | `/api/v1/reservations` | Auth (Driver/Admin) | Reserve a parking spot (Concurrency-safe) |
| **Reservations** | `GET` | `/api/v1/reservations/my-reservations` | Auth (Driver/Admin) | List my reservations (Paginated) |
| **Reservations** | `DELETE`| `/api/v1/reservations/:id` | Auth (Driver/Admin) | Cancel a reservation |
| **Reservations** | `GET` | `/api/v1/reservations` | Admin | List all reservations (Paginated, Searchable) |

---

## 📄 Pagination, Sorting, & Filter Parameters

All listing routes (`GET /api/v1/zones`, `GET /api/v1/reservations`, and `GET /api/v1/reservations/my-reservations`) integrate with `utils.QueryBuilder` to support standard query filters:

| Parameter | Default | Sample | Description |
| :--- | :--- | :--- | :--- |
| `page` | `1` | `?page=2` | Target page index (1-based) |
| `limit` | `10` | `?limit=15` | Records limit per page (capped at `100`) |
| `sort` | `created_at` | `?sort=name` | Database column sorting index |
| `order` | `desc` | `?order=asc` | Ordering direction (`asc` or `desc`) |
| `search` | `""` | `?search=Electric` | Substring search matching text fields (case-insensitive `ILIKE`) |

### Search Field Configurations
- **Parking Zones Search:** Queries against columns `name` and `type` (e.g. `ev_charging`, `standard`).
- **Reservations Search:** Queries against columns `license_plate` and `status` (e.g. `active`, `cancelled`).

---

## 📦 Standardized API Responses

The API uses a standard JSON response format to ensure consistency.

### 1. Success Response (Paginated Lists)
```json
{
  "success": true,
  "message": "Parking zones retrieved successfully",
  "data": [
    {
      "id": 1,
      "name": "Airport Term-1 EV Charge",
      "type": "ev_charging",
      "total_capacity": 5,
      "created_at": "2026-06-26T10:00:00Z",
      "updated_at": "2026-06-26T10:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 1,
    "total_page": 1
  }
}
```

### 2. Success Response (Single Resource)
```json
{
  "success": true,
  "message": "Parking zone created successfully",
  "data": {
    "id": 2,
    "name": "Mall of Asia Regular B1",
    "type": "standard",
    "total_capacity": 50
  }
}
```

### 3. Validation or Client Errors (400 Bad Request)
```json
{
  "success": false,
  "message": "Validation failed",
  "details": "Key: 'CreateZoneRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag"
}
```

### 4. Not Found Exception (404 Not Found)
```json
{
  "success": false,
  "message": "The requested resource was not found"
}
```

### 5. Method Not Allowed (405 Method Not Allowed)
```json
{
  "success": false,
  "message": "Method not allowed"
}
```

### 6. Internal Server Error (500 Server Error)
```json
{
  "success": false,
  "message": "Internal server error",
  "details": "connection refused by database server"
}
```

---

## 🔒 Concurrency Control & Safety

To prevent double bookings and race conditions when multiple users attempt to reserve the last remaining parking spot simultaneously, SpotSync implements **Pessimistic Locking** (`SELECT ... FOR UPDATE`) within database transactions:

```go
func (r *reservationRepository) CreateWithLock(reservation *models.Reservation) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Lock the parking zone row to block concurrent reservations on the same zone
		var zone models.ParkingZone
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&zone, reservation.ZoneID).Error; err != nil {
			return err
		}

		// 2. Count active reservations for this zone
		var activeCount int64
		if err := tx.Model(&models.Reservation{}).
			Where("zone_id = ? AND status = ?", reservation.ZoneID, models.ReservationStatusActive).
			Count(&activeCount).Error; err != nil {
			return err
		}

		// 3. Verify slot availability
		if int(activeCount) >= zone.TotalCapacity {
			return ErrZoneFull // Returns rollback automatically
		}

		// 4. Register the reservation atomically
		return tx.Create(reservation).Error
	})
}
```

This locking mechanism forces concurrent requests targeting the same `zone_id` to execute sequentially, guaranteeing that the database never exceeds a zone's designated capacity.

---

## 🌱 Environment Configuration Variables Reference

| Variable | Required | Default | Description |
| :--- | :--- | :--- | :--- |
| `PORT` | No | `8080` | Port on which the application web server listens |
| `DSN` | **Yes** | — | PostgreSQL Database Source Name connection details |
| `JWT_SECRET` | **Yes** | — | Security signature seed used to encode and decode user JWTs |
| `ADMIN_NAME` | No | `Admin` | Default name given to the seeded admin user |
| `ADMIN_EMAIL`| **Yes** | — | Email credentials used to seed/log in to the admin account |
| `ADMIN_PASSWORD`| **Yes** | — | Password assigned to the seeded admin user (hashed before save) |

---

## 📜 License

Distributed under the MIT License. See `LICENSE` for more information.
