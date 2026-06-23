#!/usr/bin/env bash
#
# Smoke test: build the server, boot it, hit the live API with a known
# request, and assert the response looks right. Run after each push.
#
# Usage: PORT=8080 ./scripts/smoke.sh
set -euo pipefail

PORT="${PORT:-8080}"
BASE="http://127.0.0.1:${PORT}"
BINARY="bin/server"

cleanup() {
  if [[ -n "${SERVER_PID:-}" ]]; then
    kill "${SERVER_PID}" 2>/dev/null || true
    wait "${SERVER_PID}" 2>/dev/null || true
  fi
}
trap cleanup EXIT

echo "==> building"
go build -o "${BINARY}" ./cmd/server

echo "==> starting server on :${PORT}"
PORT="${PORT}" "./${BINARY}" &
SERVER_PID=$!

echo "==> waiting for /healthz"
for _ in $(seq 1 50); do
  if curl -fsS "${BASE}/healthz" >/dev/null 2>&1; then
    break
  fi
  sleep 0.1
done

echo "==> POST /suggest"
RESP="$(curl -fsS -X POST "${BASE}/suggest" \
  -H 'Content-Type: application/json' \
  -d '{"ingredient":"butter"}')"
echo "    response: ${RESP}"

if ! echo "${RESP}" | grep -q '"alternatives"'; then
  echo "SMOKE FAILED: response missing \"alternatives\"" >&2
  exit 1
fi

echo "==> checking 404 for unknown ingredient"
CODE="$(curl -s -o /dev/null -w '%{http_code}' -X POST "${BASE}/suggest" \
  -H 'Content-Type: application/json' \
  -d '{"ingredient":"unobtanium"}')"
if [[ "${CODE}" != "404" ]]; then
  echo "SMOKE FAILED: expected 404 for unknown ingredient, got ${CODE}" >&2
  exit 1
fi

echo "SMOKE OK"
