#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
DIST="$ROOT/dist"

mkdir -p "$DIST"

echo "Building Go binary for linux/amd64..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o "$DIST/server" ./cmd/server

echo "Creating deployment zip..."
(cd "$DIST" && zip -j function.zip server)

echo "Done: $DIST/function.zip ($(du -sh "$DIST/function.zip" | cut -f1))"
