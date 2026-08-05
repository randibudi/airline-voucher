FROM golang:1.26-alpine AS build

WORKDIR /src/backend

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api

FROM alpine:3.23

RUN addgroup -S app && adduser -S -G app app \
    && mkdir -p /data \
    && chown app:app /data

WORKDIR /app

COPY --from=build --chown=app:app /out/api ./api

USER app

EXPOSE 8080

ENTRYPOINT ["/app/api"]
