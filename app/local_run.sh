# local_run.sh
#!/bin/bash
set -e

# Path to Go API binary (relative to /app)
GO_BINARY=./backend/go/cmd/server/server

# Rebuild Go API
echo "🔨 Building Go API..."
(cd backend/go && go build -o ./cmd/server/server ./cmd/server)

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