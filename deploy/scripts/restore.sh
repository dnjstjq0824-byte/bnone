#!/usr/bin/env bash
set -euo pipefail

if [ -z "${1:-}" ]; then
  echo "Usage: $0 <path_to_backup.sql.gz>"
  exit 1
fi

BACKUP_FILE="$1"

if [ ! -f "${BACKUP_FILE}" ]; then
  echo "Error: Backup file ${BACKUP_FILE} does not exist."
  exit 1
fi

echo "[RESTORE] Restoring PostgreSQL database from ${BACKUP_FILE}..."
PGPASSWORD="${POSTGRES_PASSWORD:-bnone_secret}" gunzip -c "${BACKUP_FILE}" | psql -h "${POSTGRES_HOST:-localhost}" -U "${POSTGRES_USER:-bnone}" -d "${POSTGRES_DB:-bnone_db}"

echo "[RESTORE] Database restored successfully."

