# Dockerfile for API Gateway (Train Service)
FROM golang:1.22-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum (dependency files)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY main.go .

# Build the Go binary
RUN go build -o api-gateway main.go

# Use lightweight alpine for final image
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/api-gateway .

# Expose port 8080 (Train Service entrance)
EXPOSE 8080

# Command to start the Train Service
CMD ["./api-gateway"]
