# ================================
# go/pgx/Dockerfile
# ================================

# ---- Builder ----
FROM golang:1.25.1-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy go.mod and go.sum first (cache dependencies)
COPY go/go.mod go/go.sum ./
RUN go mod download

# Copy source
COPY go/pgx/main.go ./pgx/

# Build (static binary, no CGO)
RUN CGO_ENABLED=0 GOOS=linux \
    go build -trimpath -ldflags="-s -w" -o /app/server ./pgx

# ---- Runtime ----
FROM alpine:3.20

# For HTTPS/TLS certs
RUN apk add --no-cache ca-certificates tzdata

# Non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

COPY --from=builder --chown=appuser:appgroup /app/server .

USER appuser

EXPOSE 8080

CMD ["./server"]