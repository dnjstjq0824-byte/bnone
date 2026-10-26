#!/usr/bin/env bash
set -euo pipefail

BACKUP_DIR="${BACKUP_DIR:-/backups}"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_FILE="${BACKUP_DIR}/bnone_db_${TIMESTAMP}.sql.gz"

mkdir -p "${BACKUP_DIR}"

echo "[BACKUP] Starting PostgreSQL database backup to ${BACKUP_FILE}..."
PGPASSWORD="${POSTGRES_PASSWORD:-bnone_secret}" pg_dump -h "${POSTGRES_HOST:-localhost}" -U "${POSTGRES_USER:-bnone}" -d "${POSTGRES_DB:-bnone_db}" | gzip > "${BACKUP_FILE}"

echo "[BACKUP] Backup completed successfully. Size: $(du -h "${BACKUP_FILE}" | cut -f1)"

# Keep last 14 days of backups
find "${BACKUP_DIR}" -name "bnone_db_*.sql.gz" -mtime +14 -delete
echo "[BACKUP] Retention policy applied (retained 14 days)."

