# Stage 1: Build the Go binary
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server -ldflags="-w -s" ./cmd/api/main.go

# Stage 2: Create the final, minimal image
FROM alpine:latest

# Add a non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Create the application directory
RUN mkdir -p /app

# Explicitly copy the binary from the builder to the absolute path /app/server
COPY --from=builder /app/server /app/server

# Set ownership of the entire app directory
RUN chown -R appuser:appgroup /app

# Set the user
USER appuser

# Set the working directory
WORKDIR /app

# Expose the port
EXPOSE 8080

# Start the server using a relative path from the working directory
CMD ["/bin/sh", "-c", "set -a && . /secrets/.env && set +a && /app/server"]