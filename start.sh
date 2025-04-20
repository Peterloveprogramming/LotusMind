#!/bin/sh 
set -e

echo "run db migration"
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start up the app"
echo "$@"

echo "starting server"
exec /app/main "$@"  # Use exec to replace the shell with the main app