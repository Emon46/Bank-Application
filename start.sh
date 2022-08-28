#!/bin/sh

set -e
echo "run postgres db migration"
source /app/app.env
/app/migrate --path /app/migration --database "$DB_SOURCE" -verbose up

echo "start our bank application"
ls -la /app
exec "$@"
