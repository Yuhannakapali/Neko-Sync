# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X main.version=${VERSION:-dev} -X main.buildTime=$(date -u '+%Y-%m-%d_%H:%M:%S')" \
    -o nekosync \
    ./cmd/nekosync

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create app user
RUN addgroup -g 1001 -S nekosync && \
    adduser -u 1001 -S nekosync -G nekosync

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /build/nekosync .

# Copy configuration files
COPY --from=builder /build/.env.example .env.example

# Create upload directory
RUN mkdir -p uploads && \
    chown -R nekosync:nekosync /app

# Switch to app user
USER nekosync

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./nekosync"]
