#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
MIGRATIONS_DIR="$ROOT_DIR/internal/storage/pgstorage/migrations"

# 🔥 Правильный URL: postgresql://user:password@host:port/dbname
DEFAULT_DATABASE_URL="postgresql://teammate_search:teammate_search@teammate-db:5432/teammates_data?sslmode=disable"
DATABASE_URL="${DATABASE_URL:-$DEFAULT_DATABASE_URL}"

echo "Ожидание доступности PostgreSQL..."

# Ждём, пока БД будет готова принимать подключения
until pg_isready -h teammate-db -p 5432 --username=teammate_search --dbname=teammates_data --timeout=1; do
  echo "БД пока недоступна, ждём..."
  sleep 2
done

echo "БД доступна, применяем миграции..."

if ! command -v psql >/dev/null 2>&1; then
  echo "Ошибка: psql не установлен. Установите postgresql-client." >&2
  exit 1
fi

for migration in $(ls "$MIGRATIONS_DIR"/*.sql | sort); do
  echo "Применяем миграцию: $migration"
  psql "$DATABASE_URL" -v ON_ERROR_STOP=1 -q -f "$migration"
done

echo "✅ Все миграции успешно применены"
