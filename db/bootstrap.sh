#!/bin/sh -eu
echo "db-init starting…"

# wait loop (30 × 2 s ≈ 1 min)
i=0
until pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_ADMIN_USER"; do
  i=$((i+1))
  [ $i -ge 30 ] && { echo "postgres never became ready"; exit 1; }
  echo "[$i/30] db not ready…"
  sleep 2
done

export PGPASSWORD="$DB_ADMIN_PASSWORD"
psql -v ON_ERROR_STOP=1 \
     -h "$DB_HOST" -p "$DB_PORT" -U "$DB_ADMIN_USER" \
     --set=worker="$DB_USER" \
     --set=wpw="$DB_PASSWORD" \
     --set=appdb="$DB_NAME" \
     -f /scripts/init.sql

echo "role $DB_USER and db $DB_NAME ready."
