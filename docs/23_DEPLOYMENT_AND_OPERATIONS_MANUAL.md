# 23. DEPLOYMENT AND OPERATIONS MANUAL

**Document ID:** BN-ONE-OPS-023  
**Version:** 1.0.0  
**Effective Date:** 2026-10-24  

---

## 1. Container Architecture & Orchestration

BN ONE delivers complete containerized definitions for local development, staging environments, and production Kubernetes/Docker clusters.

```text
bn-one/
  docker-compose.yml       <-- Full-stack orchestration (PostgreSQL, Redis, API, Worker, Web, Admin)
  services/
    api/Dockerfile         <-- Multi-stage Go 1.27 scratch / alpine container
    worker/Dockerfile      <-- Multi-stage Go 1.27 background worker container
  apps/
    web/Dockerfile         <-- Multi-stage Node.js build + Nginx alpine production image
    admin/Dockerfile       <-- Multi-stage Node.js build + Nginx alpine production image
```

---

## 2. Docker Compose Quickstart

To run the complete production-grade system locally with a single command:

```bash
# 1. Clone repository
git clone https://github.com/bookingnetwork/bn-one.git
cd bn-one

# 2. Boot all services
docker compose up -d --build

# 3. Apply database migrations & seed fixtures
docker compose exec api /app/migrate -path=/app/migrations -database="postgres://bnone:bnone_secret@postgres:5432/bnone_db?sslmode=disable" up
docker compose exec postgres psql -U bnone -d bnone_db -f /fixtures/01_seed_initial_data.sql
```

### Port Allocations:
- **User Web:** `http://localhost:3000`
- **Admin Console:** `http://localhost:3001`
- **API Server:** `http://localhost:8080/api/v1`
- **PostgreSQL:** `localhost:5432`
- **Redis:** `localhost:6379`

---

## 3. Environment Variables Reference

| Variable Name | Default Value | Description |
|---|---|---|
| `PORT` | `8080` | API listener TCP port |
| `DATABASE_URL` | `postgres://bnone:bnone_secret@localhost:5432/bnone_db?sslmode=disable` | Database connection string |
| `REDIS_URL` | `redis://localhost:6379/0` | Cache & background queue broker |
| `JWT_SECRET` | *Must be 32+ char secret* | HMAC-SHA256 signature key |
| `BASE_CURRENCY` | `KRW` | Platform default fallback currency |
| `LOG_LEVEL` | `info` | Logging verbosity (`debug`, `info`, `warn`, `error`) |

---

## 4. Disaster Recovery & Backup Procedures

1. **Database Snapshots:** PostgreSQL automated daily snapshots via WAL-G / pg_dump stored in encrypted S3/GCS buckets with 30-day retention.
2. **Point-In-Time Recovery (PITR):** Continuous Write-Ahead Log (WAL) archiving enables recovery to any point within the last 14 days.
3. **Health Check Probes:** `/healthz` and `/readyz` endpoints monitored by Kubernetes Liveness and Readiness probes.

