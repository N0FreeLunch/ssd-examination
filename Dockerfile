# syntax=docker/dockerfile:1

# ------------------------------------------------------------------------------
# Stage 1: Builder
# ------------------------------------------------------------------------------
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Retrieve Litestream
# (We download the binary directly from the official release for the builder stage)
ADD https://github.com/benbjohnson/litestream/releases/download/v0.3.13/litestream-v0.3.13-linux-amd64.tar.gz /tmp/litestream.tar.gz
RUN tar -C /usr/local/bin -xzf /tmp/litestream.tar.gz

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build the application
# We copy the rest of the source code. 
# Note: .dockerignore MUST exclude sdd-examination-spec to avoid context errors.
COPY . .

# Build the binary
# -ldflags="-w -s" reduces binary size by stripping debug info
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/server ./cmd/server

# ------------------------------------------------------------------------------
# Stage 2: Final
# ------------------------------------------------------------------------------
FROM alpine:3.19

WORKDIR /app

# Install runtime dependencies
# ca-certificates is needed for HTTPS requests (e.g. S3)
RUN apk add --no-cache ca-certificates bash sqlite

# Copy Litestream binary
COPY --from=builder /usr/local/bin/litestream /usr/local/bin/litestream

# Copy application binary
COPY --from=builder /app/server /app/server

# Setup entrypoint script
# We will use this to possibly run litestream replicate -exec logic
COPY <<EOF /app/entrypoint.sh
#!/bin/bash
set -e

# If REPLICA_URL is set, run with litestream
if [ -n "\$REPLICA_URL" ]; then
    echo "Starting with Litestream replication to \$REPLICA_URL"
    exec litestream replicate -exec "/app/server"
else
    echo "Starting server directly (no replication)"
    exec /app/server
fi
EOF

RUN chmod +x /app/entrypoint.sh
ENTRYPOINT ["/app/entrypoint.sh"]
