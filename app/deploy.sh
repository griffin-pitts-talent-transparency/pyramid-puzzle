#!/bin/bash
# Deploy static files + SQLite data + API to lebron-games.com

set -euo pipefail

REMOTE=ssh-digitalocean-lebron-games

STATIC_LOCAL_DIR="static"
STATIC_REMOTE_DIR="/srv/default/www"

DATA_LOCAL_DIR="data"
DATA_REMOTE_DIR="/srv/default/data"

SECURE_LOCAL_DIR="secure"
SECURE_REMOTE_DIR="/srv/default/secure"

API_LOCAL_DIR="./cmd/server"
API_REMOTE_BIN="/srv/default/bin"
API_BINARY_NAME="pyramid-puzzle-server"

echo "--------------------------------------------------------"
echo "🚀 Deploying static site → $REMOTE:$STATIC_REMOTE_DIR"

# Ensure remote directory exists
ssh "$REMOTE" "mkdir -p $STATIC_REMOTE_DIR"

# Deploy root index.html
scp static/index.html "$REMOTE:$STATIC_REMOTE_DIR/"

# Deploy partials/
rsync -avz --delete \
    static/partials/ "$REMOTE:$STATIC_REMOTE_DIR/partials/"

# Deploy css/
rsync -avz --delete \
    static/css/ "$REMOTE:$STATIC_REMOTE_DIR/css/"

# Deploy js/
rsync -avz --delete \
    static/js/ "$REMOTE:$STATIC_REMOTE_DIR/js/"

echo "--------------------------------------------------------"
echo "📦 Deploying SQLite DB + SQL files → $REMOTE:$DATA_REMOTE_DIR"
ssh "$REMOTE" "mkdir -p $DATA_REMOTE_DIR"
rsync -avz \
    --exclude='.DS_Store' \
    "$DATA_LOCAL_DIR/" "$REMOTE:$DATA_REMOTE_DIR/"

echo "--------------------------------------------------------"
echo "🔐 Deploying secure SQL files (excluding userdata.db) → $REMOTE:$SECURE_REMOTE_DIR"
ssh "$REMOTE" "mkdir -p $SECURE_REMOTE_DIR/queries $SECURE_REMOTE_DIR/schema $SECURE_REMOTE_DIR/data"

# Do NOT deploy userdata.db automatically — must be manually deployed once if needed
rsync -avz --delete \
    "$SECURE_LOCAL_DIR/queries/" "$REMOTE:$SECURE_REMOTE_DIR/queries/"

rsync -avz --delete \
    "$SECURE_LOCAL_DIR/schema/" "$REMOTE:$SECURE_REMOTE_DIR/schema/"

echo "--------------------------------------------------------"
echo "🔧 Building API binary (CGO enabled, glibc target)"
CGO_ENABLED=1 \
CC=/opt/homebrew/bin/x86_64-unknown-linux-gnu-gcc \
GOOS=linux GOARCH=amd64 \
go build -o "$API_BINARY_NAME" "$API_LOCAL_DIR"

echo "--------------------------------------------------------"
echo "🚀 Deploying API binary → $REMOTE:$API_REMOTE_BIN"
ssh "$REMOTE" "mkdir -p $API_REMOTE_BIN"
scp "$API_BINARY_NAME" "$REMOTE:$API_REMOTE_BIN/pyramid-puzzle-server.new"
ssh "$REMOTE" "mv $API_REMOTE_BIN/pyramid-puzzle-server.new $API_REMOTE_BIN/pyramid-puzzle-server"
rm "$API_BINARY_NAME"

echo "--------------------------------------------------------"
echo "🔁 Restarting API service (pyramid-puzzle-server)"
ssh "$REMOTE" "sudo systemctl restart pyramid-puzzle-server"

echo "--------------------------------------------------------"
echo "🔁 Reloading Caddy (static site)"
ssh "$REMOTE" "sudo systemctl reload caddy"

echo "--------------------------------------------------------"
echo "✅ Deployment complete! Visit https://lebron-games.com"
