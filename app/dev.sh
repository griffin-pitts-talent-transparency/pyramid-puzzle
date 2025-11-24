#!/bin/bash
set -e

# Frontend only — no Go server yet
echo "Starting Caddy server for static files..."
caddy run
