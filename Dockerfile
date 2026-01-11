# --- Build Stage ---
FROM golang:1.25.5-bookworm AS builder

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# --- Final Stage ---
FROM debian:bookworm-slim

# Install CA certificates and create a non-root user
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/* \
    && groupadd -r appgroup \
    && useradd -r -g appgroup -d /home/appuser -m appuser

WORKDIR /app

COPY --from=builder --chown=appuser:appgroup /app/server .

USER appuser

EXPOSE 8888

CMD ["./server"]