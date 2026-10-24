# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.work go.work.sum ./
COPY services/api/go.mod services/api/go.sum ./services/api/
COPY services/worker/go.mod services/worker/go.sum ./services/worker/

WORKDIR /app/services/worker
RUN go mod download

WORKDIR /app
COPY services/worker ./services/worker

WORKDIR /app/services/worker
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /bin/worker ./cmd

# Final stage
FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /bin/worker /app/worker

ENTRYPOINT ["/app/worker"]
