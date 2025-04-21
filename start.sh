#!/bin/sh 
set -e

if [ "$RUN_DB_MIGRATION" = "true" ]; then
  echo "run db migration"
  /app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up
fi

echo "start up the app"
echo "$@"

echo "starting server"
exec /app/main "$@"  # Use exec to replace the shell with the main app