# Cinema Seat Booking

A high-concurrency movie theater seat booking application built with **Go** and **Redis**, featuring a responsive vanilla HTML/JS frontend. 

This project explores strict transaction concurrency, preventing race conditions (multiple users trying to book the exact same seat) through a well-designed session-based locking model backed by Redis.

## Features

- **Distributed Concurrency Control**: Uses Redis `SETNX` (Set if Not exists) mechanisms to guarantee atomic seat locking, so two users can never overlap a booking.
- **TTL-Based Holds (Shopping Cart)**: When a user selects a seat, it locks it for 2 minutes (`defaultHoldTTL`). After 2 minutes, Redis' native Time-To-Live expiration automatically frees the seat for others if the user didn't finalize their purchase.
- **Reverse Session Lookup**: Built-in secondary indexing on Redis to quickly map `sessionID` back to a `seatKey` when it's time to process payment and confirm a booking.
- **Real-time Frontend Polling**: The frontend aggressively polls the backend every 2 seconds, displaying color-coded visual changes to users instantly when someone else locks a seat in real-time.
- **Clean Go Architecture**: Decoupled layered architecture (`Handler -> Service -> Interface -> Redis`).

---

## Tech Stack

- **Backend:** Go 1.22+ (`net/http`)
- **Database/Cache:** Redis (`go-redis/v9`)
- **Frontend:** Vanilla HTML, CSS (Custom Properties), JavaScript

---

## 📁 Project Structure

```text
.
├── cmd/
│   └── main.go                  # Main application entrypoint, sets up HTTP router & wires dependencies
├── internal/
│   ├── adapters/redis/          # Redis connection pool & client configuration
│   ├── booking/                 # Core booking domain
│   │   ├── domain.go            # Models (Booking) and Interfaces (BookingStore)
│   │   ├── handler.go           # HTTP Handlers (HoldSeat, ConfirmSession, etc.)
│   │   ├── service.go           # Business logic layer (delegates to the Store)
│   │   ├── redis_store.go       # Persistent store using Redis Session Holds/Locks
│   │   └── memory_store.go      # In-memory generic store fallback
│   └── utils/
│       └── utils.go             # Generic JSON Response Formatting
├── static/
│   └── index.html               # Frontend UI
├── docker-compose.yaml          # Deploy instances of Redis and Redis-Commander
├── go.mod
└── go.sum
```

---

## API Endpoints

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/movies` | Fetch available movies and their theater row/seat counts. |
| `GET` | `/movies/{id}/seats` | Fetch all currently locked or confirmed seats for the movie. |
| `POST` | `/movies/{id}/seats/{seatId}/hold` | Acquires a 2-minute lock on a specific seat. Returns a `sessionID`. |
| `PUT` | `/sessions/{id}/confirm` | Submits payment logic. Strips the TTL off the Redis key, confirming it. |
| `DELETE` | `/sessions/{id}` | Releases a cart hold instantly, making the seat available to others. |

---

## Getting Started (Local Development)

### Prerequisites
- Go 1.22+
- Docker & Docker Compose

### 1. Launch Redis
Start up the local Redis instance and Redis-Commander backend using Docker Compose.
```bash
docker compose up -d
```

### 2. Run the Go Backend
```bash
go run cmd/main.go
```

### 3. Open the App!
Navigate to [http://localhost:8080/](http://localhost:8080/) in your web browser. 

You can try simulating a race condition by opening the URL in two side-by-side incognito browser windows, picking a movie, and simultaneously fighting over the exact same seat!
