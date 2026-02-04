# Build stage
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /bin/iot-gateway ./cmd/iot-gateway

# Runtime stage
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /bin/iot-gateway /app/iot-gateway

# Copy config files
COPY configs/ /app/configs/

# Create non-root user
RUN adduser -D -g '' appuser
USER appuser

EXPOSE 8084

ENTRYPOINT ["/app/iot-gateway"]
