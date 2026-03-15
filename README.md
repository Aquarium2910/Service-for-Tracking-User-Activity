# Activity Tracker Service

A Go-based REST API service with **PostgreSQL** storage, designed to record user activity events and produce aggregated statistics. The project includes a background worker for data aggregation and a minimal React client for monitoring.

## 🚀 Features

- **Event Recording**: Create user activity events with flexible metadata.
- **Filtering**: Retrieve events filtered by user ID and date range.
- **Background Worker**: Automatic aggregation of user activity every 4 hours.
- **Dockerized**: Fully containerized environment for easy setup and scaling.

## 🛠 Tech Stack

- **Backend**: Go 1.24, [Echo Framework](https://echo.labstack.com/), [pgx](https://github.com/jackc/pgx)
- **Frontend**: React (Vite), CSS3
- **Database**: PostgreSQL
- **Infrastructure**: Docker, Docker Compose

---

## 🏃 How to Run

### Option 1: Using Docker (Recommended)
The simplest way to run everything (Database, Backend, and Frontend) is using Docker Compose.

1. Ensure you have **Docker** (which includes Compose v2) installed.
2. Clone the repository and navigate to the root directory.
3. Run the following command:
   ```bash
   docker compose up --build
   ```
4. Access the services:
   - **Frontend Dashboard**: [http://localhost:3000](http://localhost:3000)
   - **Backend API**: [http://localhost:8080](http://localhost:8080)

### Option 2: Local Development

1. **Database**: Start a PostgreSQL instance. You can use the one in Docker:
   ```bash
   docker compose up activity_db
   ```
2. **Backend**:
   - Copy `.env.example` to `.env` and adjust the database credentials.
   - Run the server:
     ```bash
     go run cmd/server/main.go
     ```
3. **Frontend**:
   - Navigate to the `frontend` directory.
   - Install dependencies: `npm install`
   - Start the dev server: `npm run dev`
   - Open: [http://localhost:5173](http://localhost:5173)

---

## 📡 Sample API Requests

### 1. Create an Activity Event
**POST** `/api/v1/events`
```bash
curl -X POST http://localhost:8080/api/v1/events \
-H "Content-Type: application/json" \
-d '{
  "user_id": 42,
  "action": "page_view",
  "metadata": {"page": "/home", "ip": "127.0.0.1"}
}'
```

### 2. Get Events with Filters
**GET** `/api/v1/events?user_id=42&start_date=2024-01-01T00:00:00Z&end_date=2024-12-31T23:59:59Z`
```bash
curl "http://localhost:8080/api/v1/events?user_id=42"
```

---

## 🕒 Background Job Description

The service includes a background worker (located in `internal/worker/cron.go`) that runs **every 4 hours**.
- **Task**: Aggregates the total number of events created by each user during the last 4-hour interval.
- **Storage**: Results are saved into the `activity_stats` table for historical reporting.
- **Automatic Setup**: The worker starts automatically with the Go server.

---

## 📝 Notes & Optional Parts

- **Database Initialization**: The system uses `init.sql` mounted via Docker to automatically create the necessary tables (`events`, `activity_stats`) on the first run.
- **Structured Logging**: Implemented using the standard `log/slog` package for JSON-formatted logs.
- **CORS Handling**: Configured to allow cross-origin requests from both development (port 5173) and Docker (port 3000) environments.
- **Metadata**: Meta-data is stored as `JSONB` in PostgreSQL, allowing for complex and searchable event data.

---

## 📂 Project Structure
- `cmd/server/`: Main entry point for the API server.
- `internal/api/`: HTTP handlers and routing.
- `internal/database/`: Repository layer for PostgreSQL interaction.
- `internal/models/`: Data structures and entity definitions.
- `internal/service/`: Business logic.
- `internal/worker/`: Background aggregation worker.
- `frontend/`: React application (UI).
- `init.sql`: Initial database schema.
