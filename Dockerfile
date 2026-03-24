# -------- Build --------
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Install git
RUN apk add --no-cache git

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN go build -o web-analyzer cmd/server/main.go

# -------- Run --------
FROM alpine:latest

WORKDIR /app

# Install CA certificates
RUN apk add --no-cache ca-certificates

# Copy binary from builder
COPY --from=builder /app/web-analyzer .

# Copy templates/static files if needed
COPY templates ./templates
COPY static ./static

# Expose app port
EXPOSE 8080

# Run the app
CMD ["./web-analyzer"]