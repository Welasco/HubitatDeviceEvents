# Copilot Instructions

## Project Overview

HubitatDeviceEvents is a Go backend + SolidJS frontend application that records and displays device events from a [Hubitat Elevation](https://hubitat.com/) smart home hub. The Hubitat hub pushes events to this app via webhooks, and the app stores them in MySQL for querying and visualization.

## Build & Run

### Go backend

```sh
go build -o hubitatdeviceevents
./hubitatdeviceevents    # listens on :3000 (or $PORT)
```

### Frontend (SolidJS + Vite)

```sh
cd frontend
npm install
npm run dev    # dev server
npm run build  # production build → frontend/dist/
```

### Docker (full stack)

```sh
docker build -f Dockerfile -t hubitatdeviceevents .
```

The multi-stage Dockerfile builds the Go binary, then the SolidJS frontend, and combines them into an Alpine image. In production the frontend is served from `./public`; in dev mode (`Environment=Dev`) it's served from `./frontend/dist`.

### No tests exist yet

There are no `_test.go` files or frontend test scripts.

## Architecture

```
main.go                  → entry point: init logger, start Fiber server
transport/trasnport.go   → Fiber app setup, route registration, static file serving
device/device.go         → HTTP handlers (business logic between routes and DB)
database/
  database_interface.go  → Database interface all backends must implement
  database_mysql.go      → MySQL implementation (only active backend)
model/                   → Shared data structs: Device, DeviceEvent, Config, UrlQueries
common/logger/           → Custom leveled logger (Err/Warn/Inf/Deb) writing to file + stdout
common/http/             → Simple HTTP GET client used to fetch device info from the Hubitat hub
frontend/                → SolidJS + Vite + Tailwind CSS + AG Grid UI
```

### Request flow

1. Routes are defined in `transport.setupRoutes()` under `/api/v1/`.
2. Each route calls a handler in `device/device.go`.
3. Handlers use the `database.Database` interface (package-level `db` variable) for persistence.
4. On `RegisterDeviceEvent`, if the device doesn't exist in the DB, the app auto-fetches it from the Hubitat Maker API and inserts it before storing the event.

### Database interface pattern

`database/database_interface.go` defines the `Database` interface. New database backends should implement this interface and be wired up in `device.DBInit()`. Currently only MySQL is implemented.

## Configuration

All config is via environment variables, loaded automatically from `.env` by `godotenv/autoload`. See `.env.sample` for required values:

| Variable | Purpose |
|---|---|
| `ConnectionString` | MySQL DSN (e.g. `root:password@tcp(host:3306)/hubitatdeviceevents`) |
| `DatabaseType` | Must be `mysql` (only supported backend) |
| `HubitatGetDevicesUrl` | Hubitat Maker API URL for auto-registering devices |
| `logLevel` | `Err`, `Warn`, `Inf`, or `Deb` |
| `logPath` | Log file path prefix |
| `PORT` | Server port (default `3000`) |
| `Environment` | Set to `Dev` to serve frontend from `./frontend/dist` |

## Key Conventions

- **Logging format**: Always use `[package][function] message` — e.g. `logger.Info("[device][DBInit] Initializing database")`.
- **API versioning**: All endpoints live under `/api/v1/`.
- **Model JSON tags**: `Device.Id` uses uppercase `"Id"` (matching Hubitat API); other fields are lowercase.
- **Error responses**: Return `fiber.Map{"status": "error", "message": "...", "data": err.Error()}` with appropriate HTTP status codes.
- **Database auto-init**: `NewMySQLDB` creates the database and tables automatically on startup if they don't exist.
- **Frontend API base URL**: Set via `window.base_url` in the HTML host page, not hardcoded.
