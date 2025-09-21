# Multi-stage build for optimized image size
FROM golang:1.21-alpine AS builder

# Install required packages
RUN apk add --no-cache git ca-certificates tzdata

# Create working directory
WORKDIR /app

# Copy go.mod and go.sum for dependency caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o hepic-app-server-v2 .

# Final image
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Create user for security
RUN addgroup -g 1001 -S hepic && \
    adduser -u 1001 -S hepic -G hepic

# Create working directory
WORKDIR /app

# Copy binary file from builder stage
COPY --from=builder /app/hepic-app-server-v2 .

# Copy configuration file
COPY --from=builder /app/config.env .

# Set permissions
RUN chown -R hepic:hepic /app
RUN chmod +x hepic-app-server-v2

# Switch to non-privileged user
USER hepic

# Expose port
EXPOSE 8080

# Environment variables
ENV GIN_MODE=release
ENV PORT=8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/health || exit 1

# Start application
CMD ["./hepic-app-server-v2"]
