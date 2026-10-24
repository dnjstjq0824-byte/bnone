# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.work go.work.sum ./
COPY services/api/go.mod services/api/go.sum ./services/api/
COPY services/worker/go.mod services/worker/go.sum ./services/worker/

WORKDIR /app/services/api
RUN go mod download

WORKDIR /app
COPY services/api ./services/api
COPY database ./database

WORKDIR /app/services/api
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /bin/api ./cmd/server

# Final stage
FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /bin/api /app/api
COPY database/migrations /app/migrations

EXPOSE 8080
ENTRYPOINT ["/app/api"]
