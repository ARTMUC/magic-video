# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy migrations
COPY --from=builder /app/internal/migrations ./internal/migrations

# Install goose for migrations
RUN apk add --no-cache curl && \
    curl -L https://github.com/pressly/goose/releases/download/v3.15.0/goose_linux_x86_64 -o /usr/local/bin/goose && \
    chmod +x /usr/local/bin/goose

EXPOSE 8080

CMD ["./main"]
