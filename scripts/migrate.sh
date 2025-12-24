#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
MIGRATIONS_DIR="$ROOT_DIR/internal/storage/pgstorage/migrations"

DEFAULT_DATABASE_URL="postgresql://teammate_data:teammate_data@localhost:5432/teammate_data?sslmode=disable"
DATABASE_URL="${DATABASE_URL:-$DEFAULT_DATABASE_URL}"

if ! command -v psql >/dev/null 2>&1; then
  echo "psql client is required to run migrations" >&2
  exit 1
fi

for migration in $(ls "$MIGRATIONS_DIR"/*.sql | sort); do
  echo "Applying migration: $migration"
  psql "$DATABASE_URL" -f "$migration"
done

