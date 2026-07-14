# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Install git and ca-certificates for fetching dependencies
RUN apk add --no-cache git ca-certificates

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o melodihub .

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN adduser -D -g '' melodihub

# Copy binary from builder
COPY --from=builder /app/melodihub .

# Change ownership to non-root user
RUN chown -R melodihub:melodihub /app

USER melodihub

EXPOSE 3000

CMD ["./melodihub"]
