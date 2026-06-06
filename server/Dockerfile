# Build stage
FROM golang:1.26-alpine AS builder

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary (CGO not needed for pure Go PostgreSQL driver)
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# Runtime stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Set timezone
ENV TZ=Asia/Shanghai

# Create non-root user
RUN addgroup -g 1000 -S app && \
    adduser -u 1000 -S app -G app

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/server .

# Copy data files (tarot, hexagrams, horoscope data)
COPY --from=builder /build/data ./data

# Create data directory
RUN mkdir -p /app/data && chown -R app:app /app

# Switch to non-root user
USER app

# Expose Render's default web-service port. Local docker-compose overrides the
# application port with ZHANBU_SERVER_PORT=8080 and maps it to localhost:18080.
EXPOSE 10000

# Run the server
CMD ["./server"]
