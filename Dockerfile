# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy all source code
COPY . .

# Build the binary (static linking for alpine)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o main .

# Final stage
FROM alpine:3.19

WORKDIR /app

# Install ca-certificates for HTTPS requests to external APIs
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN adduser -D -g '' appuser

# Copy binary from builder
COPY --from=builder /app/main .

# Ensure binary is executable
RUN chmod +x /app/main

# Use non-root user
USER appuser

# Expose port (Cloud Run uses PORT env var)
EXPOSE 8080

# Run the binary
ENTRYPOINT ["./main"]