# Airline Voucher Seat Assignment

An airline voucher seat assignment application built with React, TypeScript, Vite, Mantine, Go, Fiber, and SQLite.

## Status

The backend core supports voucher lookup, idempotent voucher generation, aircraft-specific seat assignment, and SQLite persistence. The frontend remains a foundation shell.

## Requirements

- Node.js 22
- Go 1.26

## Frontend Development

Install dependencies and start the development server from the repository root:

```sh
npm --prefix frontend install
npm --prefix frontend run dev
```

## Backend Development

From the `backend` directory, start the API server:

```sh
go run ./cmd/api
```

The server listens on `http://localhost:8080`. SQLite uses `vouchers.db` by default. Set `DATABASE_PATH` to use another location:

```sh
DATABASE_PATH=/path/to/vouchers.db go run ./cmd/api
```

The schema is initialized automatically at startup.

### API Endpoints

- `GET /api/health` returns backend health as `{ "status": "ok" }`.
- `POST /api/check` accepts `flightNumber` and `flightDate` (`YYYY-MM-DD`), then returns `{ "exists": boolean, "voucher": object | null }`.
- `POST /api/generate` accepts `crewName`, `crewId`, `flightNumber`, `flightDate` (`YYYY-MM-DD`), and `aircraft`, then returns the existing or newly generated voucher.

Supported aircraft names are `ATR`, `Airbus 320`, and `Boeing 737 Max`. Voucher responses contain `crewName`, `crewId`, normalized `flightNumber`, `flightDate`, `aircraft`, and exactly three `seats`.

## Quality Checks

Run the frontend checks from the repository root:

```sh
npm --prefix frontend run format:check
npm --prefix frontend run lint
npm --prefix frontend run typecheck
npm --prefix frontend run test
npm --prefix frontend run build
```

Run the backend checks from the `backend` directory:

```sh
test -z "$(gofmt -l .)"
go vet ./...
go test ./... -count=1
go build ./...
```
