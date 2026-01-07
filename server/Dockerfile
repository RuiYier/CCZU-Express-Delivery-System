FROM golang:1.25-alpine AS builder

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /src

# Copy only the server module files first for better cache
COPY server/go.mod server/go.sum ./server/
WORKDIR /src/server
RUN go mod download

# Copy rest of the repository
WORKDIR /src
COPY . .

# Build the server binary from the server module
WORKDIR /src/server
RUN go build -o /app/server .

FROM alpine:3.18

# runtime dependencies if any
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server /app/server

# Prepare config directory (will be mounted by docker-compose in dev)
RUN mkdir -p /app/config

ENV CONFIG_PATH=/app/config/config.toml

EXPOSE 8088

ENTRYPOINT ["/app/server"]
