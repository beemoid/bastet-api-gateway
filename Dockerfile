# =============================================================================
# Multi-stage Dockerfile — API Gateway
# =============================================================================

# ── Stage 1: Build ────────────────────────────────────────────────────────────
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Download dependencies first (layer cache)
COPY go.mod go.sum ./
RUN go mod download

# Build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o api-gateway main.go

# ── Stage 2: Runtime ──────────────────────────────────────────────────────────
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

# Non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

# Binary
COPY --from=builder /app/api-gateway ./api-gateway

# Runtime assets
COPY --from=builder /app/templates   ./templates
COPY --from=builder /app/docs        ./docs

RUN chown -R appuser:appuser /app

USER appuser

# Port is read from SERVER_PORT env var — default 8080
EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:${SERVER_PORT:-8080}/health || exit 1

CMD ["./api-gateway"]
