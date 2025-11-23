# Build Stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git make

# Copy dependency definitions
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binaries
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/server ./cmd/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/worker ./cmd/worker/main.go

# Runtime Stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for external API calls (SSL)
RUN apk --no-cache add ca-certificates tzdata

# Copy binaries from builder
COPY --from=builder /bin/server /app/server
COPY --from=builder /bin/worker /app/worker
# Copy default config (optional, usually mounted or env vars in prod)
COPY --from=builder /app/configs /app/configs

# Expose API port
EXPOSE 8080

# Default entrypoint (can be overridden to /app/worker)
CMD ["/app/server"]
