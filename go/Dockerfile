FROM golang:1.25.1-alpine AS builder
ARG SOURCE_PATH=go/pgx
RUN apk add --no-cache git ca-certificates
WORKDIR /app
COPY go/go.mod go/go.sum ./
RUN go mod download
COPY $SOURCE_PATH/main.go ./server/
RUN CGO_ENABLED=0 GOOS=linux \
    go build -trimpath -ldflags="-s -w" -o /app/server ./server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup
WORKDIR /app
COPY --from=builder --chown=appuser:appgroup /app/server .
USER appuser
EXPOSE 8080
CMD ["./server"]