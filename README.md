# 🚗 SpotSync API

> Smart Parking & EV Charging Reservation Platform

A centralized backend API for managing parking zones and handling high-demand reservation of limited EV charging spots at airports and malls.

**Live URL:** `https://spotsync-api.onrender.com`

---

## ✨ Features

- **User Authentication** — Register & Login with JWT-based Bearer token auth
- **Role-Based Access Control** — `driver` and `admin` roles with middleware-enforced permissions
- **Parking Zone Management** — Full CRUD for parking zones (admin) with real-time availability
- **Reservation System** — Concurrency-safe parking spot reservation with row-level locking
- **EV Charging Support** — Dedicated `ev_charging` zone type with capacity management

---

## 🛠️ Tech Stack

| Technology | Purpose |
|-----------|---------|
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

| Layer | Directory | Responsibility |
|-------|-----------|----------------|
| **DTO** | `dto/` | Request payloads and response structures |
| **Handler** | `handler/` | HTTP binding, validation, response formatting |
| **Service** | `service/` | Business logic, JWT generation, capacity checks |
| **Repository** | `repository/` | GORM database operations, transactions, row locks |
| **Models** | `models/` | GORM structs representing database tables |

**Dependency Injection** is done manually in `cmd/main.go`:
```
Repository → Service → Handler → Routes
```

---

## 📁 Project Structure

```
spot-sync/
├── cmd/main.go              # Entry point & DI wiring
├── config/                  # Environment & database config
├── dto/                     # Request/Response data transfer objects
├── handler/                 # HTTP handlers (Echo v4)
├── middleware/               # JWT auth & role authorization
├── models/                  # GORM database models
├── repository/              # Data access layer (GORM)
├── routes/                  # Route registration
├── service/                 # Business logic layer
├── .env.example             # Environment template
├── .air.toml                # Hot-reload config
├── Dockerfile               # Production build
└── docker-compose.yaml      # Docker setup
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
   git clone https://github.com/yourusername/spotsync-api.git
   cd spotsync-api
   ```

2. **Configure environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

3. **Required `.env` variables:**
   ```
   PORT=5000
   DSN="host=localhost user=postgres password=yourpassword dbname=spotsync port=5432 sslmode=disable"
   JWT_SECRET=your-secret-key-here
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

---

## 🌐 API Endpoints

### Authentication
| Method | Endpoint | Access | Description |
|--------|----------|--------|-------------|
| `POST` | `/api/v1/auth/register` | Public | Register a new user |
| `POST` | `/api/v1/auth/login` | Public | Login and receive JWT |

### Parking Zones
| Method | Endpoint | Access | Description |
|--------|----------|--------|-------------|
| `GET` | `/api/v1/zones` | Public | Get all parking zones |
| `GET` | `/api/v1/zones/:id` | Public | Get a single parking zone |
| `POST` | `/api/v1/zones` | Admin | Create a parking zone |
| `PUT` | `/api/v1/zones/:id` | Admin | Update a parking zone |
| `DELETE` | `/api/v1/zones/:id` | Admin | Delete a parking zone |

### Reservations
| Method | Endpoint | Access | Description |
|--------|----------|--------|-------------|
| `POST` | `/api/v1/reservations` | Auth | Reserve a parking spot |
| `GET` | `/api/v1/reservations/my-reservations` | Auth | View my reservations |
| `DELETE` | `/api/v1/reservations/:id` | Auth | Cancel a reservation |
| `GET` | `/api/v1/reservations` | Admin | View all reservations |

### Authentication Header
```
Authorization: Bearer <your-jwt-token>
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

## 📜 License

MIT
