# Airline Voucher Seat Assignment

An airline voucher seat assignment application built with React, TypeScript, Vite, Mantine, Go, Fiber, and SQLite.

## Requirements

- Node.js 22
- Go 1.26

## Run Locally

Start the backend from the repository root:

```sh
cd backend
go run ./cmd/api
```

The backend runs at `http://localhost:8080`. SQLite uses `vouchers.db` by default and initializes its schema automatically. Set `DATABASE_PATH` to use another database location.

In another terminal, start the frontend from the repository root:

```sh
npm --prefix frontend run dev
```

The frontend uses the Vite default URL, normally `http://localhost:5173`. During development, Vite proxies relative `/api` requests to `http://localhost:8080`.

## Voucher Workflow

The frontend normalizes the flight number and sends the flight date as `YYYY-MM-DD`. It first calls `POST /api/check` with `flightNumber` and `flightDate`. If a voucher exists, it displays that voucher without generating another one. Otherwise, it calls `POST /api/generate` with `crewName`, `crewId`, `flightNumber`, `flightDate`, and `aircraft`.

Supported aircraft names are `ATR`, `Airbus 320`, and `Boeing 737 Max`. Voucher responses contain `crewName`, `crewId`, normalized `flightNumber`, `flightDate`, `aircraft`, and exactly three `seats`. Dates are displayed in the frontend as `DD-MM-YYYY`.

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
