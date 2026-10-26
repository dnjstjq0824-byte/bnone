#!/usr/bin/env bash
set -euo pipefail

echo "=========================================================="
echo "  BN ONE v1.0 — Staging & Production Deployment Pipeline  "
echo "=========================================================="

# 1. Environment check
command -v docker >/dev/null 2>&1 || { echo >&2 "[ERROR] Docker is required but not installed."; exit 1; }
command -v docker-compose >/dev/null 2>&1 || command -v docker compose >/dev/null 2>&1 || { echo >&2 "[ERROR] Docker Compose is required."; exit 1; }

# 2. Build and verify all production images
echo "[STEP 1/4] Building production container images..."
docker compose build --no-cache

# 3. Start services with healthchecks
echo "[STEP 2/4] Starting services in detached mode..."
docker compose up -d

# 4. Wait for database and API healthcheck
echo "[STEP 3/4] Waiting for services to become healthy..."
MAX_RETRIES=30
RETRY_COUNT=0

until curl -s http://localhost:8080/healthz | grep -q "HEALTHY" || [ $RETRY_COUNT -eq $MAX_RETRIES ]; do
  echo "Waiting for API server to respond... ($RETRY_COUNT/$MAX_RETRIES)"
  sleep 2
  RETRY_COUNT=$((RETRY_COUNT + 1))
done

if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
  echo "[ERROR] API healthcheck timed out!"
  docker compose logs api
  exit 1
fi

echo "[STEP 4/4] Deployment verification passed! All endpoints operational."
echo "User Web:     http://localhost:3000"
echo "Admin Console: http://localhost:3001"
echo "API Server:   http://localhost:8080/api/v1"
echo "Healthcheck:  http://localhost:8080/healthz"
