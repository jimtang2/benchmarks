# ================================
# go/pq/Dockerfile
# ================================

# ---- Builder ----
FROM golang:1.25.1-alpine AS builder

RUN apk add --no-cache git ca-certificates gcc musl-dev

WORKDIR /app

# lib/pq needs CGO for native libpq
ENV CGO_ENABLED=1

COPY go/go.mod go/go.sum ./
RUN go mod download

COPY go/pq/main.go ./pq/

RUN CGO_ENABLED=1 GOOS=linux \
    go build -trimpath -ldflags="-s -w" -o /app/server ./pq

# ---- Runtime (includes libpq deps) ----
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata libpq

RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

COPY --from=builder --chown=appuser:appgroup /app/server .

USER appuser

EXPOSE 8080

CMD ["./server"]