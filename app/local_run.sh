#!/bin/bash
set -e

# Path to Go API binary
GO_BINARY=./cmd/server/server

# Rebuild Go API (optional: only needed if you want auto-rebuild)
echo "🔨 Building Go API..."
go build -o "$GO_BINARY" ./cmd/server

# Start Go server in background
echo "🚀 Starting Go API..."
"$GO_BINARY" &
GO_PID=$!

# Start Caddy
echo "🌐 Starting Caddy..."
caddy run --config ./Caddyfile --adapter caddyfile

# Cleanup Go server when Caddy stops
echo "🛑 Caddy stopped. Killing Go API (PID $GO_PID)..."
kill $GO_PID
