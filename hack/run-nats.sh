#!/usr/bin/env bash
set -euo pipefail

PORT="${NATS_PORT:-4222}"
MONITOR_PORT="${NATS_MONITOR_PORT:-8222}"
DATA_DIR="${NATS_DATA_DIR:-./data/nats}"

echo "[INFO] Starting NATS (port=$PORT, monitor=$MONITOR_PORT, data=$DATA_DIR)"

mkdir -p "$DATA_DIR"

if command -v nats-server >/dev/null 2>&1; then
  echo "[INFO] Using local nats-server binary."
  exec nats-server \
    --port "$PORT" \
    --http_port "$MONITOR_PORT" \
    --jetstream \
    --store_dir "$DATA_DIR"
else
  echo "[INFO] nats-server not found. Falling back to Docker."
  docker run --rm -it \
    -p "${PORT}:4222" \
    -p "${MONITOR_PORT}:8222" \
    -v "$(pwd)/${DATA_DIR}:/data" \
    nats:latest \
    -js -p 4222 -m 8222 -sd /data
fi