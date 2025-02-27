#!/bin/sh
set -e

# 等待数据库就绪
echo "Waiting for postgres..."

# 运行数据库迁移
echo "Running migrations..."
. /app/app.env
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

# 启动应用
echo "Starting the app..."
exec /app/main "$@"
