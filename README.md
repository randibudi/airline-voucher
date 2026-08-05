# Airline Voucher Seat Assignment

A foundation for an airline voucher seat assignment application, built with React, TypeScript, Vite, Mantine, Go, and Fiber.

## Status

Phase 0 — Checkpoint B: Backend Foundation.

The current application provides a minimal frontend shell and a backend health endpoint. Voucher forms, voucher APIs, persistence, and seat assignment business logic are not included in this checkpoint.

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

The server listens on `http://localhost:8080` by default and provides:

```text
GET /api/health
```

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
