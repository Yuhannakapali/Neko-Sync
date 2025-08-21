# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /build

# Copy go mod files
COPY apps/backend/go.mod apps/backend/go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY apps/backend/ ./

# Build arguments
ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.gitCommit=${GIT_COMMIT}" \
    -o nekosync \
    ./cmd/nekosync

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S nekosync && \
    adduser -S nekosync -u 1001 -G nekosync

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /build/nekosync /app/nekosync

# Change ownership
RUN chown nekosync:nekosync /app/nekosync

# Switch to non-root user
USER nekosync

# Expose port (adjust as needed)
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./nekosync"]
